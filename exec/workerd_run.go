package exec

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
	"vvorker/conf"
	"vvorker/defs"
	"vvorker/entities"
	"vvorker/services/control"
	"vvorker/utils"
	"vvorker/utils/database"

	"github.com/go-co-op/gocron/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type execManager struct {
	//用于外层循坏的退出
	signMap *defs.SyncMap[string, bool]
	//用于执行cancel函数
	chanMap *defs.SyncMap[string, chan struct{}]
	// 用于存储进程 ID
	pidMap *defs.SyncMap[string, int]
	// 用于记录 worker 运行状态
	runningMap *defs.SyncMap[string, bool]
	// 用于存储调度器实例
	scheduler gocron.Scheduler
}

var ExecManager *execManager

type WorkerLogData struct {
	UID    string    `gorm:"index" json:"uid"`
	Output string    `json:"output"`
	Time   time.Time `gorm:"index" json:"time"`
	Type   string    `gorm:"index" json:"type"`
	LogUID string    `gorm:"index" json:"log_uid"`
}

// 定义合并后的日志模型
type WorkerLog struct {
	gorm.Model
	*WorkerLogData
}

var (
	// 合并后的日志 channel，用于异步处理
	workerLogChan = make(chan WorkerLog, 1000)
	// 批量插入的大小
	batchSize = 100
)

// 初始化时启动合并后的日志处理 goroutine
func init() {
	var scheduler gocron.Scheduler
	var err error
	if scheduler, err = gocron.NewScheduler(); err != nil {
		logrus.Fatalf("Failed to initialize scheduler: %v", err)
	}

	ExecManager = &execManager{
		signMap:    new(defs.SyncMap[string, bool]),
		chanMap:    new(defs.SyncMap[string, chan struct{}]),
		pidMap:     new(defs.SyncMap[string, int]),
		runningMap: new(defs.SyncMap[string, bool]), // 初始化运行状态映射
		scheduler:  scheduler,                       // 初始化调度器
	}

	ExecManager.scheduler.Start()
	go processWorkerLogs()
}

// 处理合并后日志的函数
func processWorkerLogs() {
	var logs []WorkerLog
	for {
		select {
		case log := <-workerLogChan:
			logs = append(logs, log)
			if len(logs) >= batchSize {
				// 批量插入数据库
				if err := dbCreateWorkerLogs(logs); err != nil {
					logrus.Errorf("Failed to batch insert worker logs: %v", err)
				}
				logs = nil
			}
		case <-time.After(2 * time.Second):
			if len(logs) > 0 {
				// 定时批量插入
				if err := dbCreateWorkerLogs(logs); err != nil {
					logrus.Errorf("Failed to batch insert worker logs: %v", err)
				}
				logs = nil
			}
		}
	}
}

// 批量插入合并后日志到数据库的函数
func dbCreateWorkerLogs(logs []WorkerLog) error {
	db := database.GetDB()
	return db.CreateInBatches(logs, len(logs)).Error
}

func (m *execManager) RunCmd(uid string, argv []string) {
	if _, ok := m.chanMap.Get(uid); ok {
		logrus.Warnf("workerd %s is already running!", uid)
		return
	}

	db := database.GetDB()
	var worker entities.Worker
	if err := db.Where("uid = ?", uid).First(&worker).Error; err != nil {
		logrus.Warnf("workerconfig error: %v", err)
		return
	}
	workerconfig, werr := conf.ParseWorkerConfig(worker.Template)
	if werr != nil {
		logrus.Warnf("workerconfig error: %v", werr)
		workerconfig = conf.DefaultWorkerConfig()
	}
	schedluers := workerconfig.Schedulers
	allJobs := make([]gocron.Job, 0)
	for _, scheduler := range schedluers {
		if scheduler.Cron == "" {
			continue
		}
		s := ExecManager.scheduler
		j, err := s.NewJob(
			gocron.CronJob(scheduler.Cron, true),
			gocron.NewTask(
				control.SendSchedulerEvent,
				uid,
				scheduler.Cron,
			),
		)
		if err != nil {
			logrus.Warnf("Failed to create scheduler job: %v", err)
			continue
		}
		logrus.Infof("workerd %s scheduler job created! id: %s", uid, j.ID())
		allJobs = append(allJobs, j)
	}

	c := make(chan struct{})
	m.chanMap.Set(uid, c)
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context, uid string, argv []string, m *execManager) {
		defer func(uid string, m *execManager) {
			m.signMap.Delete(uid)
		}(uid, m)

		logrus.Infof("workerd %s running!", uid)
		workerdDir := filepath.Join(
			conf.AppConfigInstance.WorkerdDir,
			defs.WorkerInfoPath,
			uid,
		)

		for {
			// 检查上下文是否被取消，如果取消则退出循环
			select {
			case <-ctx.Done():
				logrus.Infof("workerd %s context cancelled, exiting loop", uid)
				for _, job := range allJobs {
					ExecManager.scheduler.RemoveJob(job.ID())
				}
				return
			default:
			}

			args := []string{"serve",
				filepath.Join(workerdDir, defs.CapFileName),
			}
			args = append(args, "--verbose")
			args = append(args, argv...)

			cmd := exec.CommandContext(ctx, conf.AppConfigInstance.WorkerdBinPath, args...)
			cmd.Dir = workerdDir
			cmd.SysProcAttr = &syscall.SysProcAttr{}

			// 创建一个管道来捕获标准输出
			stdoutPipe, err := cmd.StdoutPipe()
			if err != nil {
				logrus.Errorf("Failed to create stdout pipe for workerd %s: %v", uid, err)
			}

			// 创建一个管道来捕获错误输出
			stderrPipe, err := cmd.StderrPipe()
			if err != nil {
				logrus.Errorf("Failed to create stderr pipe for workerd %s: %v", uid, err)
			}

			if err := cmd.Start(); err != nil {
				logrus.Errorf("Failed to start workerd %s: %v", uid, err)
				m.runningMap.Set(uid, false)

				continue
			}

			// 保存进程 ID
			m.pidMap.Set(uid, cmd.Process.Pid)

			// 读取标准输出并发送到 channel
			go func(uid string) {
				buf := make([]byte, 1024)
				for {
					select {
					case <-ctx.Done(): // 监听上下文取消信号
						return
					default:
						n, err := stdoutPipe.Read(buf)
						if n > 0 {
							workerLogChan <- WorkerLog{
								WorkerLogData: &WorkerLogData{
									UID:    uid,
									Output: string(buf[:n]),
									Time:   time.Now(),
									Type:   "stdout",
									LogUID: utils.GenerateUID(),
								},
							}
							logrus.Infof("workerd %s stdout: %s", uid, string(buf[:n]))
						}
						if err != nil {
							return
						}
					}
				}
			}(uid)

			// 读取错误输出并发送到 channel
			go func(uid string) {
				buf := make([]byte, 1024)
				for {
					select {
					case <-ctx.Done(): // 监听上下文取消信号
						return
					default:
						n, err := stderrPipe.Read(buf)
						if n > 0 {
							workerLogChan <- WorkerLog{
								WorkerLogData: &WorkerLogData{
									UID:    uid,
									Output: string(buf[:n]),
									Time:   time.Now(),
									Type:   "error",
									LogUID: utils.GenerateUID(),
								},
							}
							logrus.Errorf("workerd %s error: %s", uid, string(buf[:n]))
						}
						if err != nil {
							return
						}
					}
				}
			}(uid)
			m.runningMap.Set(uid, true)

			if err := cmd.Wait(); err != nil {
				logrus.Errorf("Workerd %s exited with error: %v", uid, err)
				m.runningMap.Set(uid, false)
			}

			if exit, ok := m.signMap.Get(uid); ok && exit {
				return
			}
			time.Sleep(3 * time.Second)
		}
	}(ctx, uid, argv, m)

	go func(cancel context.CancelFunc, uid string, m *execManager) {
		defer func(uid string, m *execManager) {
			m.chanMap.Delete(uid)
			// m.pidMap.Delete(uid) // 不要试图删除pid
		}(uid, m)

		if channel, ok := m.chanMap.Get(uid); ok {
			<-channel
			cancel() // 调用 cancel 函数取消上下文
		}
	}(cancel, uid, m)
}

// ExitCmd 根据 uid 停止某个正在运行的 worker
func (m *execManager) ExitCmd(uid string) {
	defer func(uid string, m *execManager) {
		m.signMap.Delete(uid)
		m.runningMap.Set(uid, false) // 标记 worker 为停止状态
	}(uid, m)

	if channel, ok := m.chanMap.Get(uid); ok {
		channel <- struct{}{}
		logrus.Infof("workerd %s is being stopped!", uid)
	} else {
		logrus.Warnf("workerd %s is not running, cannot stop it!", uid)
	}

	// 尝试获取进程 ID
	pid, ok := m.pidMap.Get(uid)
	if !ok {
		logrus.Warnf("No process ID found for workerd %s", uid)
		return
	} else {
		logrus.Infof("workerd %s pid is %d", uid, pid)
	}

	// 获取进程句柄
	process, err := os.FindProcess(pid)
	if err != nil {
		logrus.Errorf("Failed to find process for workerd %s: %v", uid, err)
		return
	}

	// 等待进程退出
	_, err = process.Wait()
	if err != nil {
		logrus.Errorf("Error waiting for workerd %s to exit: %v", uid, err)
	} else {
		logrus.Infof("workerd %s has stopped", uid)
	}
}

func (m *execManager) ExitAllCmd() {
	for uid := range m.chanMap.ToMap() {
		m.ExitCmd(uid)
	}
}

func (m *execManager) GetWorkerStatusByUID(uid string) int {
	mu, ok := m.runningMap.Get(uid)
	if !ok {
		return 0
	}
	if mu {
		return 1
	}
	return 0
}

package exec

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
	"vorker/conf"
	"vorker/defs"

	"github.com/sirupsen/logrus"
)

type execManager struct {
	//用于外层循坏的退出
	signMap *defs.SyncMap[string, bool]
	//用于执行cancel函数
	chanMap *defs.SyncMap[string, chan struct{}]
	// 新增：用于存储进程 ID
	pidMap *defs.SyncMap[string, int]
}

var ExecManager *execManager

func init() {
	ExecManager = &execManager{
		signMap: new(defs.SyncMap[string, bool]),
		chanMap: new(defs.SyncMap[string, chan struct{}]),
		pidMap:  new(defs.SyncMap[string, int]),
	}
}

func (m *execManager) RunCmd(uid string, argv []string) {
	if _, ok := m.chanMap.Get(uid); ok {
		logrus.Warnf("workerd %s is already running!", uid)
		return
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
			args := []string{"serve",
				filepath.Join(workerdDir, defs.CapFileName),
			}
			// 判断操作系统是否为 Windows
			if runtime.GOOS != "windows" {
				args = append(args, "--watch")
			}
			args = append(args, "--verbose")
			args = append(args, argv...)
			cmd := exec.CommandContext(ctx, conf.AppConfigInstance.WorkerdBinPath, args...)
			cmd.Dir = workerdDir
			cmd.SysProcAttr = &syscall.SysProcAttr{}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				logrus.Errorf("Failed to start workerd %s: %v", uid, err)
				time.Sleep(3 * time.Second)
				continue
			}
			// 保存进程 ID
			m.pidMap.Set(uid, cmd.Process.Pid)

			if err := cmd.Wait(); err != nil {
				logrus.Errorf("Workerd %s exited with error: %v", uid, err)
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
			// m.pidMap.Delete(uid)
		}(uid, m)

		if channel, ok := m.chanMap.Get(uid); ok {
			<-channel
			m.signMap.Set(uid, true)
			cancel()
			return
		} else {
			logrus.Errorf("workerd %s is not running!", uid)
			return
		}
	}(cancel, uid, m)
}

// ExitCmd 根据 uid 停止某个正在运行的 worker
func (m *execManager) ExitCmd(uid string) {
	if channel, ok := m.chanMap.Get(uid); ok {
		channel <- struct{}{}
		logrus.Infof("workerd %s is being stopped!", uid)
	} else {
		logrus.Warnf("workerd %s is not running, cannot stop it!", uid)
	}
	// 如果是windows，需要等待workerd退出
	if runtime.GOOS == "windows" {
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
}

func (m *execManager) ExitAllCmd() {
	for uid := range m.chanMap.ToMap() {
		m.ExitCmd(uid)
	}
}

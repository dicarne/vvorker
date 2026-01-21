package models

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"time"
	"vvorker/conf"
	"vvorker/defs"
	"vvorker/entities"
	"vvorker/exec"
	"vvorker/ext/kv/src/sys_cache"
	"vvorker/funcs"
	workercopy "vvorker/models/worker_copy"
	"vvorker/rpc"
	"vvorker/tunnel"
	"vvorker/utils"
	"vvorker/utils/database"
	"vvorker/utils/generate"

	"github.com/codeclysm/extract/v3"

	"github.com/google/uuid"
	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

type Worker struct {
	gorm.Model
	*entities.Worker
	EnableAccessControl bool   `json:"EnableAccessControl"`
	Description         string `json:"Description"`
}

func init() {
	go func() {
		if conf.AppConfigInstance.LitefsEnabled {
			if !conf.IsMaster() {
				return
			}
			utils.WaitForPort("localhost", conf.AppConfigInstance.LitefsPrimaryPort)
		}
		db := database.GetDB()

		for err := db.AutoMigrate(&Worker{}); err != nil; err = db.AutoMigrate(&Worker{}) {
			logrus.WithError(err).Errorf("auto migrate worker error, sleep 5s and retry")
			time.Sleep(5 * time.Second)
		}
	}()
	go func() {
		if conf.AppConfigInstance.LitefsEnabled {
			utils.WaitForPort("localhost", conf.AppConfigInstance.LitefsPrimaryPort)
		}
		for {
			if len(conf.AppConfigInstance.NodeID) == 0 {
				logrus.Error("[workerd init()] node is not initialized, retrying after 1 seconds")
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}
		tunnel.GetClient()
		if conf.IsMaster() {
			utils.WaitForPort("localhost", conf.AppConfigInstance.TunnelAPIPort)
			logrus.Info("Waitting for client...")
		} else {
			utils.WaitForPort(conf.AppConfigInstance.TunnelHost, conf.AppConfigInstance.TunnelAPIPort)
		}
		time.Sleep(time.Second * 2)
		NodeWorkersInit()
	}()
}

func (w *Worker) TableName() string {
	return "workers"
}

func GetWorkerByUID(userID uint, uid string) (*Worker, error) {
	var worker Worker
	db := database.GetDB()

	if err := db.Where(&Worker{
		Worker: &entities.Worker{
			UserID: uint64(userID),
		},
	}).Where(
		&Worker{
			Worker: &entities.Worker{
				UID: uid,
			},
		},
	).First(&worker).Error; err != nil {
		// 如果不是拥有者，检查是否是协作者
		if IsWorkerMember(uid, uint64(userID)) {
			if err := db.Where(&Worker{Worker: &entities.Worker{UID: uid}}).First(&worker).Error; err != nil {
				return nil, err
			}
			return &worker, nil
		}
		return nil, err
	}
	return &worker, nil
}

func HasWorker(userID uint, uid string) bool {
	var worker Worker
	db := database.GetDB()

	if err := db.Where(&Worker{
		Worker: &entities.Worker{
			UserID: uint64(userID),
		},
	}).Where(
		&Worker{
			Worker: &entities.Worker{
				UID: uid,
			},
		},
	).Select("UID").First(&worker).Error; err != nil {
		return false
	}
	return true
}

func AdminGetWorkerByName(name string) (*Worker, error) {
	var worker Worker
	db := database.GetDB()

	if err := db.Where(
		&Worker{
			Worker: &entities.Worker{
				Name: name,
			},
		},
	).First(&worker).Error; err != nil {
		return nil, err
	}
	return &worker, nil
}

type WorkerSimple struct {
	UID                 string `json:"UID"`
	EnableAccessControl bool   `json:"EnableAccessControl"`
	NodeName            string `json:"NodeName"`
}

func AdminGetWorkerByNameSimple(name string) (*WorkerSimple, error) {
	var worker WorkerSimple
	db := database.GetDB()
	bytes, err := sys_cache.GlobalCache("worker_uid_name:"+name, func() ([]byte, error) {
		if err := db.Model(&Worker{}).Select(
			"UID",
			"EnableAccessControl",
			"NodeName",
		).Where(
			&Worker{
				Worker: &entities.Worker{
					Name: name,
				},
			},
		).First(&worker).Error; err != nil {
			return nil, err
		}
		v, err := json.Marshal(worker)
		if err != nil {
			return nil, err
		}
		return v, nil
	}, 10)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(bytes, &worker)
	return &worker, nil
}

func GetWorkersByNames(userID uint, names []string) ([]*Worker, error) {
	var workers []*Worker
	db := database.GetDB()

	if err := db.Where(&Worker{
		Worker: &entities.Worker{
			UserID: uint64(userID),
		},
	}).Where("name in (?)", names).Find(&workers).Error; err != nil {
		return nil, err
	}
	return workers, nil
}

func GetWorkersByUserID(userID uint64) ([]*Worker, error) {
	var workers []*Worker
	db := database.GetDB()

	if err := db.Where(&Worker{
		Worker: &entities.Worker{
			UserID: userID,
		},
	}).Order("updated_at desc").Find(&workers).Error; err != nil {
		return nil, err
	}
	return workers, nil
}

func GetWorkersByUIDs(userID uint, uids []string) ([]*Worker, error) {
	var workers []*Worker
	db := database.GetDB()
	if err := db.Where(&Worker{
		Worker: &entities.Worker{
			UserID: uint64(userID),
		},
	}).Where("uid in (?)", uids).Find(&workers).Error; err != nil {
		return nil, err
	}
	return workers, nil
}

func AdminGetWorkersByNames(names []string) ([]*Worker, error) {
	var workers []*Worker
	db := database.GetDB()

	if err := db.Where("name in (?)", names).Find(&workers).Error; err != nil {
		return nil, err
	}
	return workers, nil
}

func GetAllWorkers(userID uint) ([]*Worker, error) {
	var workers []*Worker
	db := database.GetDB()

	// 获取用户拥有的 Workers
	if err := db.Where(&Worker{
		Worker: &entities.Worker{
			UserID: uint64(userID),
		},
	}).Order("updated_at desc").Find(&workers).Error; err != nil {
		return nil, err
	}

	// 创建一个 map 来快速查找用户拥有的 worker UIDs
	ownedWorkerUIDs := make(map[string]bool)
	for _, worker := range workers {
		ownedWorkerUIDs[worker.UID] = true
	}

	// 获取用户参与协作的 Workers
	collabWorkerUIDs, err := GetUserCollaboratedWorkers(uint64(userID))
	if err != nil {
		return workers, nil // 协作者列表获取失败不影响返回拥有的 workers
	}

	// 将协作 workers 添加到列表，但不修改 Description
	for _, uid := range collabWorkerUIDs {
		// 跳过已经是拥有的 workers
		if ownedWorkerUIDs[uid] {
			continue
		}
		var worker Worker
		if err := db.Where(&Worker{Worker: &entities.Worker{UID: uid}}).First(&worker).Error; err != nil {
			continue
		}
		workers = append(workers, &worker)
	}

	return workers, nil
}

func AdminGetAllWorkers() ([]*Worker, error) {
	var workers []*Worker
	db := database.GetDB()

	if err := db.Find(&workers).Error; err != nil {
		return nil, err
	}
	return workers, nil
}

func AdminGetAllWorkersTunnelMap() (map[string]string, error) {
	workers, err := AdminGetAllWorkers()
	if err != nil {
		return nil, err
	}
	tunnelMap := make(map[string]string)
	for _, worker := range workers {
		tunnelMap[worker.Name] = worker.TunnelID
	}
	return tunnelMap, nil
}

func AdminGetWorkersByNodeName(nodeName string) ([]*Worker, error) {
	var workers []*Worker
	db := database.GetDB()

	if err := db.Where(&Worker{
		Worker: &entities.Worker{
			NodeName: nodeName,
		},
	}).Find(&workers).Error; err != nil {
		return nil, err
	}
	return workers, nil
}

func GetWorkers(userID uint, offset, limit int) ([]*Worker, error) {
	var workers []*Worker
	db := database.GetDB()

	if err := db.Where(&Worker{
		Worker: &entities.Worker{
			UserID: uint64(userID),
		},
	}).Order("updated_at desc").Offset(offset).Limit(limit).Find(&workers).Error; err != nil {
		return nil, err
	}
	return workers, nil
}

func Trans2Entities(workers []*Worker) []*entities.Worker {
	var entities []*entities.Worker = []*entities.Worker{}
	for _, worker := range workers {
		entities = append(entities, worker.ToEntity())
	}
	return entities
}

func (w *Worker) GetWorkerClientID() string {
	return w.UID + "-worker-" + strconv.Itoa(int(w.LocalID))
}

func (w *Worker) Create() error {
	c := context.Background()
	if w.MaxCount == 0 {
		w.MaxCount = 1
	}
	db := database.GetDB()
	if w.NodeName == conf.AppConfigInstance.NodeName {
		db.Model(&workercopy.WorkerCopy{}).Unscoped().Where(&workercopy.WorkerCopy{WorkerUID: w.UID}).Delete(&workercopy.WorkerCopy{})

		for i := 0; i < int(w.MaxCount); i++ {
			w.LocalID = int32(i)

			logrus.Infof("create worker copy %v", i)

			port := tunnel.GetPortManager().ClaimWorkerPort(c, w.GetWorkerClientID())
			w.Port = port
			tunnel.GetClient().AddWorker(w.GetWorkerClientID(), utils.WorkerHostPrefix(w.GetName()), int(w.GetPort()))

			controlPort := tunnel.GetPortManager().ClaimWorkerPort(c, w.GetWorkerClientID()+"-control")
			w.ControlPort = controlPort
			tunnel.GetClient().AddWorker(w.GetWorkerClientID()+"-control", w.GetUID()+"-control", int(controlPort))

			wCopy := &workercopy.WorkerCopy{
				WorkerUID:   w.UID,
				LocalID:     uint(i),
				Port:        uint(port),
				ControlPort: uint(controlPort),
			}
			if err := db.Create(wCopy).Error; err != nil {
				logrus.WithError(err).Errorf("create worker copy error: %v", wCopy)
				return err
			}
		}

		if err := w.UpdateFile(); err != nil {
			return err
		}
		if !conf.IsMaster() && conf.AppConfigInstance.LitefsEnabled {
			return nil
		}
	} else {
		n, err := GetNodeByNodeName(w.NodeName)
		if err != nil {
			return err
		}
		wp, err := proto.Marshal(w)
		if err != nil {
			return err
		}
		go rpc.EventNotify(n.Node, defs.EventAddWorker, map[string][]byte{
			defs.KeyWorkerProto: wp,
		})
	}

	return db.Create(w).Error
}

func (w *Worker) Update() error {
	c := context.Background()
	// if w.ID == 0 {
	// 	return errors.New("worker has no id")
	// }
	db := database.GetDB()
	if w.MaxCount == 0 {
		w.MaxCount = 1
	}
	if w.NodeName == conf.AppConfigInstance.NodeName {
		workercopies := &[]workercopy.WorkerCopy{}
		db.Model(&workercopy.WorkerCopy{}).Where(&workercopy.WorkerCopy{WorkerUID: w.UID}).Find(workercopies)
		for i := range *workercopies {
			w.LocalID = int32(i)
			tunnel.GetClient().Delete(w.GetWorkerClientID())
			tunnel.GetClient().Delete(w.GetWorkerClientID() + "-control")
		}
		db.Model(&workercopy.WorkerCopy{}).Unscoped().Where(&workercopy.WorkerCopy{WorkerUID: w.UID}).Delete(&workercopy.WorkerCopy{})

		for i := 0; i < int(w.MaxCount); i++ {
			w.LocalID = int32(i)
			logrus.Infof("update worker copy %v", i)
			port := tunnel.GetPortManager().ClaimWorkerPort(c, w.GetWorkerClientID())
			w.Port = port
			tunnel.GetClient().AddWorker(w.GetWorkerClientID(), utils.WorkerHostPrefix(w.GetName()), int(port))

			controlPort := tunnel.GetPortManager().ClaimWorkerPort(c, w.GetWorkerClientID()+"-control")
			w.ControlPort = controlPort

			tunnel.GetClient().AddWorker(w.GetWorkerClientID()+"-control", w.GetUID()+"-control", int(controlPort))

			wCopy := &workercopy.WorkerCopy{
				WorkerUID:   w.UID,
				LocalID:     uint(i),
				Port:        uint(port),
				ControlPort: uint(controlPort),
			}
			if err := db.Create(wCopy).Error; err != nil {
				logrus.WithError(err).Errorf("create worker copy error: %v", wCopy)
				return err
			}
		}

		if err := w.UpdateFile(); err != nil {
			return err
		}

	}
	if !conf.IsMaster() && conf.AppConfigInstance.LitefsEnabled {
		return nil
	}

	return db.Save(w).Error
}

func (w *Worker) Delete() error {
	if w.NodeName == conf.AppConfigInstance.NodeName {
		db := database.GetDB()
		workercopies := &[]workercopy.WorkerCopy{}
		db.Model(&workercopy.WorkerCopy{}).Where(&workercopy.WorkerCopy{WorkerUID: w.UID}).Find(workercopies)
		for i := range *workercopies {
			w.LocalID = int32(i)
			tunnel.GetClient().Delete(w.GetWorkerClientID())
			tunnel.GetClient().Delete(w.GetWorkerClientID() + "-control")
		}
		db.Model(&workercopy.WorkerCopy{}).Unscoped().Where(&workercopy.WorkerCopy{WorkerUID: w.UID}).Delete(&workercopy.WorkerCopy{})
		db.Model(&WorkerMember{}).Unscoped().Where(&WorkerMember{WorkerUID: w.UID}).Delete(&WorkerMember{})
		db.Model(&WorkerInformation{}).Unscoped().Where(&WorkerInformation{WorkerInformationBase: &WorkerInformationBase{UID: w.UID}}).Delete(&WorkerInformation{})
	} else {
		n, err := GetNodeByNodeName(w.NodeName)
		if err != nil {
			logrus.WithError(err).Warnf("delete worker %s error, node %s not found, will remove it from db", w.UID, w.NodeName)
		} else {
			wp, err := proto.Marshal(w)
			if err != nil {
				return err
			}
			go rpc.EventNotify(n.Node, defs.EventDeleteWorker, map[string][]byte{
				defs.KeyWorkerProto: wp,
			})
		}
	}
	if err := w.DeleteFile(); err != nil {
		return err
	}

	if !conf.IsMaster() && conf.AppConfigInstance.LitefsEnabled {
		return nil
	}
	db := database.GetDB()

	return db.Unscoped().Where(
		&Worker{Worker: &entities.Worker{
			UID: w.UID,
		}}).Delete(&Worker{}).Error
}

func (w *Worker) Flush() error {
	if w.NodeName != conf.AppConfigInstance.NodeName {
		n, err := GetNodeByNodeName(w.NodeName)
		if err != nil {
			return err
		}
		wp, err := proto.Marshal(w)
		if err != nil {
			return err
		}
		return rpc.EventNotify(n.Node, defs.EventFlushWorker, map[string][]byte{
			defs.KeyWorkerProto: wp,
		})
	}
	if len(w.TunnelID) == 0 {
		w.TunnelID = uuid.New().String()
	}

	exec.ExecManager.ExitCmd(w.UID)

	if err := w.DeleteFile(); err != nil {
		return err
	}
	logrus.Infof("flush worker %s", w.Name)
	if err := w.Update(); err != nil {
		return err
	}

	if err := generate.GenWorkerConfig(w.ToEntity(), w); err != nil {
		return err
	}

	exec.ExecManager.RunCmd(w.UID, []string{})

	return nil
}

func (w *Worker) ToEntity() *entities.Worker {
	ans := w.Worker
	ans.Port = int32(w.GetPort())
	ans.ControlPort = int32(w.GetControlPort())
	return ans
}

func (w *Worker) DeleteFile() error {
	return os.RemoveAll(
		filepath.Join(
			conf.AppConfigInstance.WorkerdDir,
			defs.WorkerInfoPath,
			w.UID,
		),
	)
}

func (w *Worker) UpdateFile() error {
	if conf.AppConfigInstance.FileStorageUseOSS && len(w.ActiveVersionID) == 0 {
		code, err := funcs.DownloadFileFromSysBucket(fmt.Sprintf("code/%s", w.GetUID()))
		if err != nil {
			return err
		}
		codeBytes, err := io.ReadAll(code)
		if err != nil {
			return err
		}
		w.Code = codeBytes
	}
	if len(w.ActiveVersionID) == 0 {
		return utils.WriteFile(
			filepath.Join(
				conf.AppConfigInstance.WorkerdDir,
				defs.WorkerInfoPath,
				w.UID,
				defs.WorkerCodePath,
				w.Entry),
			string(w.Code))
	}

	c := context.Background()

	file, err := GetFileByVersionUID(c, w.ActiveVersionID)
	if err != nil {
		return err
	}

	return extract.Tar(c, bytes.NewReader(file.Data), filepath.Join(
		conf.AppConfigInstance.WorkerdDir, defs.WorkerInfoPath,
		w.UID, defs.WorkerCodePath), nil)
}

func (w *Worker) Run() ([]byte, error) {
	addr := fmt.Sprintf("http://%s:%d", conf.AppConfigInstance.TunnelHost, conf.AppConfigInstance.TunnelEntryPort)
	resp, err := req.C().R().SetHeader(
		defs.HeaderHost, fmt.Sprintf("%s%s", w.Name, conf.AppConfigInstance.WorkerURLSuffix),
	).Get(addr)
	if err != nil {
		return nil, err
	}
	return resp.Bytes(), nil
}

func (w *Worker) GetPort() int {
	c := context.Background()
	workerPort, ok := tunnel.GetPortManager().GetWorkerPort(c, w.GetWorkerClientID())
	if !ok {
		return 0
	}
	return int(workerPort)
}

func DiffWorkers(newWorkerList []entities.WorkerUIDVersion) ([]entities.WorkerUIDVersion, error) {
	db := database.GetDB()
	var workers []Worker
	if err := db.Model(&Worker{}).Select("uid", "version").Find(&workers).Error; err != nil {
		return nil, err
	}
	newWorkerMap := lo.SliceToMap(newWorkerList, func(w entities.WorkerUIDVersion) (string, string) { return w.UID, w.Version })
	oldWorkerMap := lo.SliceToMap(workers, func(w Worker) (string, string) { return w.UID, w.Version })
	var differentWorkers []entities.WorkerUIDVersion
	for _, newWorker := range newWorkerList {
		if _, ok := oldWorkerMap[newWorker.UID]; ok {
			if newWorker.Version != oldWorkerMap[newWorker.UID] {
				differentWorkers = append(differentWorkers, entities.WorkerUIDVersion{
					UID:     newWorker.UID,
					Version: newWorker.Version,
				})
				logrus.Infof("sync workers update worker, worker is: %+v, oldversion %s; newversion %s", newWorker, oldWorkerMap[newWorker.UID], newWorker.Version)
			}
		} else {
			differentWorkers = append(differentWorkers, entities.WorkerUIDVersion{
				UID:     newWorker.UID,
				Version: newWorker.Version,
			})
			logrus.Infof("sync workers add worker, worker is: %+v", newWorker)
		}
	}

	// delete workers that not in workerList
	for _, worker := range workers {
		if _, ok := newWorkerMap[worker.UID]; !ok {
			ow := Worker{Worker: &entities.Worker{UID: worker.UID}}
			ww := db.Where(&ow).First(&Worker{})
			if ww.Error != nil {
				continue
			}
			if err := ow.Delete(); err != nil {
				logrus.WithError(err).Errorf("sync workers delete worker error, worker is: %+v", worker)
				continue
			}
			logrus.Infof("sync workers delete worker, worker is: %+v", worker)
		}
	}

	return differentWorkers, nil
}

func SyncWorkers(workerList []entities.WorkerUIDVersion) error {
	partialFail := false
	UIDs := []string{}

	for _, workerUIDVersion := range workerList {
		worker, err := rpc.GetWorkerByUID(conf.AppConfigInstance.MasterEndpoint, workerUIDVersion.UID)
		if err != nil {
			logrus.WithError(err).Errorf("sync workers get worker error, uid is: %s, err: %v", workerUIDVersion.UID, err)
			continue
		}
		modelWorker := &Worker{Worker: worker}
		UIDs = append(UIDs, worker.UID)

		logrus.Infof("sync workers db create, new worker is: %+v", entities.Worker{
			UID: worker.GetUID(), Name: worker.GetName(), NodeName: worker.GetNodeName(),
		})

		exec.ExecManager.ExitCmd(worker.UID)

		if err := modelWorker.Delete(); err != nil && err != gorm.ErrRecordNotFound {
			logrus.WithError(err).Errorf("sync workers db delete error, worker is: %+v", worker)
			partialFail = true
			continue
		}

		if err := modelWorker.Create(); err != nil {
			logrus.WithError(err).Errorf("sync workers db create error, worker is: %+v", worker)
			partialFail = true
			continue
		}

		if err := modelWorker.DeleteFile(); err != nil {
			logrus.WithError(err).Errorf("sync workers delete file error, worker is: %+v", worker)
			partialFail = true
			continue
		}

		if err := modelWorker.UpdateFile(); err != nil {
			logrus.WithError(err).Errorf("sync workers update file error, worker is: %+v", worker)
			partialFail = true
			continue
		}
		if err := generate.GenWorkerConfig(modelWorker.ToEntity(), modelWorker); err != nil {
			logrus.WithError(err).Errorf("sync workers gen config error, worker is: %+v", worker)
			partialFail = true
			continue
		}

		exec.ExecManager.RunCmd(worker.UID, []string{})
	}

	if partialFail {
		return errors.New("partial fail")
	}

	return nil
}

func (w *Worker) WorkerNameToPort(name string) (int, error) {
	anow := Worker{}
	db := database.GetDB()
	if err := db.Where("name = ?", name).First(&anow).Error; err != nil {
		return 0, err
	}
	return anow.GetPort(), nil
}

func (w *Worker) WorkerNameToUID(name string) (string, error) {
	anow := Worker{}
	db := database.GetDB()
	if err := db.Where("name = ?", name).First(&anow).Error; err != nil {
		return "", err
	}
	return anow.GetUID(), nil
}

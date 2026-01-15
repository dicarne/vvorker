package workerd

import (
	"bytes"
	"fmt"
	"runtime/debug"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/entities"
	"vvorker/exec"
	"vvorker/ext/kv/src/sys_cache"
	"vvorker/funcs"
	"vvorker/models"
	"vvorker/utils"
	"vvorker/utils/generate"
	permissions "vvorker/utils/permissions"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UpdateWorkerReq struct {
	*entities.Worker
	Description string `json:"Description"`
}

func UpdateWorker(userID uint, UID string, worker *entities.Worker, desc string) error {
	FillWorkerValue(worker, true, UID, userID)

	workerRecord, err := models.GetWorkerByUID(userID, UID)
	if err != nil {
		return err
	}

	curNodeName := conf.AppConfigInstance.NodeName

	if workerRecord.NodeName == curNodeName {
		exec.ExecManager.ExitCmd(workerRecord.GetUID())
	}

	// 删除旧的worker
	err = workerRecord.Delete()
	if err != nil {
		return err
	}

	// 创建新的worker
	newWorker := &models.Worker{Worker: worker,
		EnableAccessControl: workerRecord.EnableAccessControl,
		Description:         desc,
	}
	newWorker.Version = utils.GenerateUID()
	code := newWorker.Code
	if conf.AppConfigInstance.FileStorageUseOSS {
		newWorker.Code = nil
	}

	if conf.AppConfigInstance.FileStorageUseOSS {
		err = funcs.UploadFileToSysBucket(fmt.Sprintf("code/%s", newWorker.GetUID()), bytes.NewReader(code))
		if err != nil {
			return err
		}
	}

	err = newWorker.Create()
	if err != nil {
		return err
	}

	if worker.NodeName == curNodeName {
		err := generate.GenWorkerConfig(newWorker.ToEntity(), newWorker)
		if err != nil {
			return err
		}
		exec.ExecManager.RunCmd(worker.GetUID(), []string{})
	}
	return nil
}

// 更新worker
func UpdateEndpointJSON(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()

	var worker *UpdateWorkerReq
	if err := permissions.BindJSON(c, &worker); err != nil {
		return
	}

	UID := worker.UID
	userID := c.GetUint(common.UIDKey)
	oldworker, err := permissions.CanWriteWorker(c, uint64(userID), UID)
	if err != nil {
		return
	}

	if worker.Worker.Code == nil {
		worker.Worker.Code = oldworker.Code
	}
	if worker.Worker.Template == "" {
		worker.Worker.Template = oldworker.Template
	}
	if worker.Worker.MaxCount == 0 {
		worker.Worker.MaxCount = oldworker.MaxCount
		if worker.Worker.MaxCount == 0 {
			worker.Worker.MaxCount = 1
		}
	}

	if err := UpdateWorker(userID, UID, worker.Worker, worker.Description); err != nil {
		logrus.WithError(err).Errorf("update worker error, worker is: [%+v]", worker.Worker)
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	if oldworker.Name != worker.Worker.Name {
		cacheKey := fmt.Sprintf("db:workerd:uid_name_%s_cache:", oldworker.Name)
		lockKey := fmt.Sprintf("db:workerd:uid_name_%s_lock:", oldworker.Name)
		sys_cache.Del(cacheKey)
		sys_cache.Del(lockKey)
	}

	common.RespOK(c, "update worker success", gin.H{
		"version": worker.Worker.Version,
	})
}

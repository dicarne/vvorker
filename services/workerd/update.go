package workerd

import (
	"runtime/debug"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/entities"
	"vvorker/exec"
	"vvorker/models"
	"vvorker/utils/generate"
	permissions "vvorker/utils/permissions"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UpdateWorkerReq struct {
	*entities.Worker
}

func UpdateEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()

	UID := c.Param("uid")
	var newWorker *UpdateWorkerReq
	if err := permissions.BindJSON(c, &newWorker); err != nil {
		return
	}

	userID := c.GetUint(common.UIDKey)
	oldworker, err := permissions.CanWriteWorker(c, uint64(userID), UID)
	if err != nil {
		return
	}

	if newWorker.Worker.Code == nil {
		// 如果没有更新代码，则从旧的代码中查找原本的代码进行覆盖
		newWorker.Worker.Code = oldworker.Code
	}
	if err := UpdateWorker(userID, UID, newWorker.Worker); err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	common.RespOK(c, "update worker success", nil)
}

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

	if err := UpdateWorker(userID, UID, worker.Worker); err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	common.RespOK(c, "update worker success", nil)
}

func UpdateWorker(userID uint, UID string, worker *entities.Worker) error {
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
	newWorker := &models.Worker{Worker: worker, EnableAccessControl: workerRecord.EnableAccessControl}
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

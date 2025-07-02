package workerd

import (
	"fmt"
	"runtime/debug"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/entities"
	"vvorker/exec"
	"vvorker/models"
	"vvorker/utils/generate"

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
	if len(UID) == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is empty", nil)
		return
	}

	var worker *UpdateWorkerReq
	if err := c.ShouldBindJSON(&worker); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}

	userID := c.GetUint(common.UIDKey)

	if worker.Worker.Code == nil {
		oldworker, err := models.GetWorkerByUID(userID, UID)
		if err != nil {
			common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
			return
		}
		worker.Worker.Code = oldworker.Code
	}
	if err := UpdateWorker(userID, UID, worker.Worker); err != nil {
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
	if err := c.ShouldBindJSON(&worker); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}
	UID := worker.UID
	if len(UID) == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is empty", nil)
		return
	}
	userID := c.GetUint(common.UIDKey)

	if worker.Worker.Code == nil {
		oldworker, err := models.GetWorkerByUID(userID, UID)
		if err != nil {
			common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
			return
		}
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
	if worker == nil {
		return fmt.Errorf("worker is nil")
	}

	curNodeName := conf.AppConfigInstance.NodeName

	if workerRecord.NodeName == curNodeName {
		exec.ExecManager.ExitCmd(workerRecord.GetUID())
	}

	err = workerRecord.Delete()
	if err != nil {
		return err
	}

	newWorker := &models.Worker{Worker: worker}
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

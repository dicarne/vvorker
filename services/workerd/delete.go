package workerd

import (
	"fmt"
	"runtime/debug"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/exec"
	"vvorker/models"
	permissions "vvorker/utils/permissions"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func DeleteEndpoint(c *gin.Context) {
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

	userID := c.GetUint(common.UIDKey)
	// 只有拥有者可以删除 worker
	_, err := permissions.CanManageWorkerMembers(c, uint64(userID), UID)
	if err != nil {
		return
	}

	if err := Delete(userID, UID); err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	models.DeleteWorkerInformationByUID(UID)

	common.RespOK(c, "delete worker success", nil)
}

func Delete(userID uint, UID string) error {
	worker, err := models.GetWorkerByUID(userID, UID)
	if err != nil {
		return err
	}
	if worker == nil {
		return fmt.Errorf("worker not found")
	}

	if worker.NodeName == conf.AppConfigInstance.NodeName {
		exec.ExecManager.ExitCmd(worker.GetUID())
	}
	if err = worker.Delete(); err != nil {
		return err
	}

	return nil
}

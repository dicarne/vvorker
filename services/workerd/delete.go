package workerd

import (
	"vvorker/common"
	"vvorker/conf"
	"vvorker/entities"
	"vvorker/exec"
	"vvorker/models"
	"vvorker/utils/database"
	permissions "vvorker/utils/permissions"

	"github.com/gin-gonic/gin"
)

func DeleteEndpoint(c *gin.Context) {

	UID := c.Param("uid")
	if len(UID) == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is empty", nil)
		return
	}

	userID, ok := common.RequireUID32(c)
	if !ok {
		return
	}
	// 只有拥有者可以删除 worker
	_, err := permissions.CanManageWorkerMembers(c, uint64(userID), UID)
	if err != nil {
		// CanManageWorkerMembers 内部已经调用了 RespErr
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
	// 权限已经在 DeleteEndpoint 中检查过了
	// 这里直接查询 worker，不进行权限检查
	db := database.GetDB()
	var worker models.Worker
	if err := db.Where(&models.Worker{Worker: &entities.Worker{UID: UID}}).First(&worker).Error; err != nil {
		return err
	}

	if worker.NodeName == conf.AppConfigInstance.NodeName {
		exec.ExecManager.ExitCmd(worker.GetUID())
	}
	if err := worker.Delete(); err != nil {
		return err
	}

	return nil
}

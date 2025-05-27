package workerd

import (
	"runtime/debug"
	"vvorker/common"
	"vvorker/exec"
	"vvorker/models"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GetWorkerLogsReq struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func GetWorkerLogsEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()

	var req *GetWorkerLogsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}

	UID := c.Param("uid")
	if len(UID) == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is empty", nil)
		return
	}

	userID := c.GetUint(common.UIDKey)
	if !models.HasWorker(userID, UID) {
		common.RespErr(c, common.RespCodeInvalidRequest, "worker not found", nil)
		return
	}

	db := database.GetDB()
	var logs []*exec.WorkerLog
	if err := db.Where("uid = ?", UID).Offset(req.Offset).Limit(req.Limit).Order("time desc").Find(&logs).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	var logs2 []*exec.WorkerLogData
	for _, log := range logs {
		logs2 = append(logs2, log.WorkerLogData)
	}
	common.RespOK(c, "get worker logs success", logs2)
}

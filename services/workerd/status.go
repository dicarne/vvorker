package workerd

import (
	"runtime/debug"
	"vvorker/common"
	"vvorker/exec"
	permissions "vvorker/utils/permissions"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GetWorkersStatusResp struct {
	UIDS []string `json:"uids"`
}

func GetWorkersStatusByUIDEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()

	req := GetWorkersStatusResp{}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, "参数解析失败", nil)
		return
	}

	userID := c.GetUint(common.UIDKey)
	// 检查用户是否有权限访问所有请求的 Workers
	for _, uid := range req.UIDS {
		_, err := permissions.CanReadWorker(c, uint64(userID), uid)
		if err != nil {
			// CanReadWorker 内部已经调用了 RespErr
			return
		}
	}

	status := make(map[string]int)
	for _, uid := range req.UIDS {
		status[uid] = exec.ExecManager.GetWorkerStatusByUID(uid)
	}
	common.RespOK(c, "ok", gin.H{
		"status": status,
	})
}

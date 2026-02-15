package workerd

import (
	"vvorker/common"
	"vvorker/exec"
	permissions "vvorker/utils/permissions"

	"github.com/gin-gonic/gin"
)

type GetWorkersStatusResp struct {
	UIDS []string `json:"uids" binding:"required,min=1"`
}

func GetWorkersStatusByUIDEndpoint(c *gin.Context) {

	req := GetWorkersStatusResp{}
	if err := c.BindJSON(&req); err != nil {
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

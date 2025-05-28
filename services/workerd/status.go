package workerd

import (
	"runtime/debug"
	"vvorker/common"
	"vvorker/exec"

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
	status := make(map[string]int)
	for _, uid := range req.UIDS {
		status[uid] = exec.ExecManager.GetWorkerStatusByUID(uid)
	}
	common.RespOK(c, "ok", gin.H{
		"status": status,
	})
}

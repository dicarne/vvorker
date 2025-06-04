package appconf

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	num, err := models.AdminGetUserNumber()
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"WorkerURLSuffix": conf.AppConfigInstance.WorkerURLSuffix,
			"Scheme":          conf.AppConfigInstance.Scheme,
			"EnableRegister":  conf.AppConfigInstance.EnableRegister || num == 0,
			"UrlType":         conf.AppConfigInstance.WorkerHostMode,
			"ApiUrl":          fmt.Sprintf("%s:%d", conf.AppConfigInstance.DefaultWorkerHost, conf.AppConfigInstance.WorkerPort),
		},
	})
}

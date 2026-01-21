package appconf

import (
	"net/http"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/models"

	"github.com/gin-gonic/gin"
)

func GetEndpoint(c *gin.Context) {

	num, err := models.AdminGetUserNumber()
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		return
	}
	urlPrefix := ""
	if conf.AppConfigInstance.WorkerHostMode == "path" {
		urlPrefix = conf.AppConfigInstance.WorkerHostPath
	}
	if conf.AppConfigInstance.SSOBaseURL != "" {
		if conf.AppConfigInstance.WorkerHostPath == "" {
			urlPrefix = conf.AppConfigInstance.SSOBaseURL
		} else {
			urlPrefix = conf.AppConfigInstance.SSOBaseURL + "/" + conf.AppConfigInstance.WorkerHostPath
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"WorkerURLSuffix": conf.AppConfigInstance.WorkerURLSuffix,
			"Scheme":          conf.AppConfigInstance.Scheme,
			"EnableRegister":  conf.AppConfigInstance.EnableRegister || num == 0,
			"UrlType":         conf.AppConfigInstance.WorkerHostMode,
			"ApiUrl":          conf.AppConfigInstance.APIWebBaseURL,
			"UrlPrefix":       urlPrefix,
			"Version":         conf.Version,
		},
	})
}

package agent

import (
	"vvorker/common"
	"vvorker/conf"
	"vvorker/defs"
	"vvorker/entities"
	"vvorker/exec"
	"vvorker/models"
	"vvorker/utils/generate"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func AddWorkerEventHandler(c *gin.Context, req *entities.NotifyEventRequest) {
	worker, err := entities.ToWorkerEntity(req.Extra[defs.KeyWorkerProto])
	if err != nil {
		logrus.WithError(err).Error("add worker event handler error")
		common.RespErr(c, common.RespCodeInvalidRequest, common.RespMsgInvalidRequest, nil)
		return
	}

	w := &models.Worker{Worker: worker}

	if err := w.Create(); err != nil {
		logrus.WithError(err).Error("add worker event handler error")
		common.RespErr(c, common.RespCodeInvalidRequest, common.RespMsgInvalidRequest, nil)
		return
	}

	if worker.NodeName == conf.AppConfigInstance.NodeName {
		if err := generate.GenWorkerConfig(w.ToEntity(), w); err != nil {
			logrus.WithError(err).Error("add worker event handler error")
			common.RespErr(c, common.RespCodeInvalidRequest, common.RespMsgInvalidRequest, nil)
			return
		}
		exec.ExecManager.RunCmd(worker.GetUID(), []string{})
	}

	logrus.Info("add worker event handler success")
	common.RespOK(c, common.RespMsgOK, nil)
}

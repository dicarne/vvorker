package agent

import (
	"vvorker/common"
	"vvorker/defs"
	"vvorker/entities"
	"vvorker/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func FlushWorkerEventHandler(c *gin.Context, req *entities.NotifyEventRequest) {
	worker, err := entities.ToWorkerEntity(req.Extra[defs.KeyWorkerProto])
	if err != nil {
		logrus.WithError(err).Error("flush worker event handler error")
		common.RespErr(c, common.RespCodeInvalidRequest, common.RespMsgInvalidRequest, nil)
		return
	}

	logrus.Infoln("flush worker event handler", worker.GetUID())

	if err := (&models.Worker{Worker: worker}).Flush(); err != nil {
		logrus.WithError(err).Error("flush worker event handler error")
		common.RespErr(c, common.RespCodeInvalidRequest, common.RespMsgInvalidRequest, nil)
		return
	}

	common.RespOK(c, common.RespMsgOK, nil)
}

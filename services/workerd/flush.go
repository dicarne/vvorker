package workerd

import (
	"runtime/debug"
	"vvorker/common"
	"vvorker/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func FlushEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	UID := c.Param("uid")
	if len(UID) == 0 {
		logrus.Errorf("uid is empty, ctx: %v", c)
		return
	}

	userID := c.GetUint(common.UIDKey)

	if err := Flush(userID, UID); err != nil {
		c.JSON(500, gin.H{"code": 3, "error": err.Error()})
		logrus.Errorf("failed to flush worker, err: %v, ctx: %v", err, c)
		return
	}

	common.RespOK(c, "flush worker success", nil)
}

func FlushAllEndpoint(c *gin.Context) {
	userID := c.GetUint(common.UIDKey)
	workers, err := models.GetAllWorkers(userID)
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	err = nil
	for _, worker := range workers {
		if err = worker.Flush(); err != nil {
			logrus.Errorf("failed to flush worker, err: %v, ctx: %v", err, c)
			continue
		}
	}
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		logrus.Warnf("partial failure, ctx: %v", c)
		return
	}

	if err := GenCapnpConfig(); err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		logrus.WithError(err).Error("gen capnp config error")
		return
	}

	common.RespOK(c, "flush worker success", nil)
}

func Flush(userID uint, UID string) error {
	worker, err := models.GetWorkerByUID(userID, UID)
	if err != nil {
		return err
	}
	err = worker.Flush()
	if err != nil {
		return err
	}
	return nil
}

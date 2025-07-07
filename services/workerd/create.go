package workerd

import (
	"runtime/debug"
	"vvorker/common"
	"vvorker/defs"
	"vvorker/entities"
	"vvorker/models"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CreateEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	worker := &models.Worker{}

	if err := c.BindJSON(worker); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}
	worker.Code = []byte(defs.DefaultCode)

	if !isCreateParamValidate() {
		common.RespErr(c, common.RespCodeInvalidRequest, common.RespMsgInvalidRequest, nil)
		return
	}
	userID := uint64(c.GetUint(common.UIDKey))

	newUID, err := Create(uint(userID), worker.Worker)
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	common.RespOK(c, "create worker success", gin.H{
		"UID":  newUID,
		"Name": worker.GetName(),
	})
}

// Create creates a new worker in the database and update the workerd capnp config file
func Create(userID uint, worker *entities.Worker) (string, error) {
	FillWorkerValue(worker, false, "", userID)

	if err := (&models.Worker{Worker: worker}).Create(); err != nil {
		logrus.Errorf("failed to create worker, err: %v", err)
		return "", err
	}

	err := Flush(userID, worker.GetUID())
	if err != nil {
		logrus.Errorf("failed to flush worker config, err: %v", err)
		return "", err
	}
	return worker.GetUID(), nil
}

func Recover(userID uint, worker *entities.Worker) error {
	db := database.GetDB()
	worker.UserID = uint64(userID)
	w2 := models.Worker{}
	if err := db.Where("uid = ?", worker.GetUID()).Assign(worker).FirstOrCreate(&w2).Error; err != nil {
		logrus.Errorf("failed to create or update worker, err: %v", err)
		return err
	}
	worker = w2.Worker
	err := Flush(userID, worker.GetUID())
	if err != nil {
		logrus.Errorf("failed to flush worker config, err: %v", err)
		return err
	}
	return nil
}

func isCreateParamValidate() bool {
	// TODO: validate the create params
	return true
}

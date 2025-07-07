package permissions

import (
	"fmt"
	"vvorker/common"
	"vvorker/entities"
	"vvorker/models"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
)

func CanReadWorker(c *gin.Context, uid uint64, worker_uid string) (*models.Worker, error) {
	if uid == 0 || worker_uid == "" {
		common.RespErr(c, common.RespCodeInvalidRequest, "invalid request param", nil)
		return nil, fmt.Errorf("invalid request param")
	}

	db := database.GetDB()
	worker := &models.Worker{}
	tx := db.Model(&models.Worker{}).Where(&models.Worker{Worker: &entities.Worker{
		UID:    worker_uid,
		UserID: uid,
	}}).First(worker)

	if tx.Error != nil {
		common.RespErr(c, common.RespCodeNotAuthed, "forbidden", nil)
		return nil, tx.Error
	}

	return worker, nil
}

func CanWriteWorker(c *gin.Context, uid uint64, worker_uid string) (*models.Worker, error) {
	if uid == 0 || worker_uid == "" {
		common.RespErr(c, common.RespCodeInvalidRequest, "invalid request param", nil)
		return nil, fmt.Errorf("invalid request param")
	}

	db := database.GetDB()
	worker := &models.Worker{}
	tx := db.Model(&models.Worker{}).Where(&models.Worker{Worker: &entities.Worker{
		UID:    worker_uid,
		UserID: uid,
	}}).First(worker)

	if tx.Error != nil {
		common.RespErr(c, common.RespCodeNotAuthed, "forbidden", nil)
		return nil, tx.Error
	}

	return worker, nil
}

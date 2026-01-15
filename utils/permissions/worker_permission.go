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
	// 检查是否是拥有者或协作者
	tx := db.Model(&models.Worker{}).Where(&models.Worker{Worker: &entities.Worker{
		UID:    worker_uid,
		UserID: uid,
	}}).First(worker)

	if tx.Error != nil {
		// 如果不是拥有者，检查是否是协作者
		if models.IsWorkerMember(worker_uid, uid) {
			// 获取 worker 详情
			tx = db.Model(&models.Worker{}).Where(&models.Worker{Worker: &entities.Worker{
				UID: worker_uid,
			}}).First(worker)
			if tx.Error != nil {
				common.RespErr(c, common.RespCodeNotAuthed, "forbidden", nil)
				return nil, tx.Error
			}
			return worker, nil
		}
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
	// 检查是否是拥有者或协作者
	tx := db.Model(&models.Worker{}).Where(&models.Worker{Worker: &entities.Worker{
		UID:    worker_uid,
		UserID: uid,
	}}).First(worker)

	if tx.Error != nil {
		// 如果不是拥有者，检查是否是协作者
		if models.IsWorkerMember(worker_uid, uid) {
			// 获取 worker 详情
			tx = db.Model(&models.Worker{}).Where(&models.Worker{Worker: &entities.Worker{
				UID: worker_uid,
			}}).First(worker)
			if tx.Error != nil {
				common.RespErr(c, common.RespCodeNotAuthed, "forbidden", nil)
				return nil, tx.Error
			}
			return worker, nil
		}
		common.RespErr(c, common.RespCodeNotAuthed, "forbidden", nil)
		return nil, tx.Error
	}

	return worker, nil
}

func CanManageWorkerMembers(c *gin.Context, uid uint64, worker_uid string) (*models.Worker, error) {
	if uid == 0 || worker_uid == "" {
		common.RespErr(c, common.RespCodeInvalidRequest, "invalid request param", nil)
		return nil, fmt.Errorf("invalid request param")
	}

	// 只有拥有者可以管理成员
	canManage, err := models.CanManageMembers(worker_uid, uid)
	if err != nil || !canManage {
		common.RespErr(c, common.RespCodeNotAuthed, "forbidden", nil)
		return nil, err
	}

	db := database.GetDB()
	worker := &models.Worker{}
	tx := db.Model(&models.Worker{}).Where(&models.Worker{Worker: &entities.Worker{
		UID: worker_uid,
	}}).First(worker)

	if tx.Error != nil {
		common.RespErr(c, common.RespCodeNotAuthed, "forbidden", nil)
		return nil, tx.Error
	}

	return worker, nil
}

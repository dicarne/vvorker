package access

import (
	"runtime/debug"
	"vvorker/common"
	"vvorker/entities"
	"vvorker/models"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EnableAccessControlRequest struct {
	Enable    bool   `json:"enable"`
	WorkerUID string `json:"worker_uid"`
}

func UpdateEnableAccessControlEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	uid := uint64(c.GetUint(common.UIDKey))
	if uid == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is required", nil)
		return
	}
	request := EnableAccessControlRequest{}
	if err := c.BindJSON(&request); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}

	db := database.GetDB()
	var user models.Worker
	if err := db.Where(&models.Worker{Worker: &entities.Worker{UID: request.WorkerUID, UserID: uid}}).First(&user).Error; err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, "worker not found", nil)
		return
	}
	user.EnableAccessControl = request.Enable
	if err := db.Save(&user).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	common.RespOK(c, common.RespMsgOK, nil)
}

type AccessControlRequest struct {
	WorkerUID string `json:"worker_uid"`
}

func GetAccessControlEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	uid := uint64(c.GetUint(common.UIDKey))
	if uid == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is required", nil)
		return
	}
	request := AccessControlRequest{}
	if err := c.BindJSON(&request); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}

	db := database.GetDB()
	var user models.Worker
	if err := db.Where(&models.Worker{Worker: &entities.Worker{UID: request.WorkerUID, UserID: uid}}).First(&user).Error; err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, "worker not found", nil)
		return
	}
	common.RespOK(c, common.RespMsgOK, gin.H{
		"EnableAccessControl": user.EnableAccessControl,
	})
}

func AddAccessRuleEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	uid := uint64(c.GetUint(common.UIDKey))
	if uid == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is required", nil)
		return
	}
	request := models.AccessRule{}
	if err := c.BindJSON(&request); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}

	db := database.GetDB()
	var user models.Worker
	if err := db.Where(&models.Worker{Worker: &entities.Worker{UID: request.WorkerUID, UserID: uid}}).First(&user).Error; err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, "worker not found", nil)
		return
	}
	if err := db.Create(&request).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	common.RespOK(c, common.RespMsgOK, nil)
}

type DeleteAccessRuleRequest struct {
	WorkerUID string `json:"worker_uid"`
	RuleID    uint   `json:"rule_id"`
}

func DeleteAccessRuleEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	uid := uint64(c.GetUint(common.UIDKey))
	if uid == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is required", nil)
		return
	}
	request := DeleteAccessRuleRequest{}
	if err := c.BindJSON(&request); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}

	db := database.GetDB()
	var user models.Worker
	if err := db.Where(&models.Worker{Worker: &entities.Worker{UID: request.WorkerUID, UserID: uid}}).First(&user).Error; err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, "worker not found", nil)
		return
	}
	if err := db.Delete(&models.AccessRule{}, &models.AccessRule{
		Model:     gorm.Model{ID: request.RuleID},
		WorkerUID: request.WorkerUID,
	}).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	common.RespOK(c, common.RespMsgOK, nil)
}

type ListAccessRuleRequest struct {
	WorkerUID string `json:"worker_uid"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
}

func ListAccessRuleEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	uid := uint64(c.GetUint(common.UIDKey))
	if uid == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is required", nil)
		return
	}
	request := ListAccessRuleRequest{}
	if err := c.BindJSON(&request); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}

	db := database.GetDB()
	var user models.Worker
	if err := db.Where(&models.Worker{Worker: &entities.Worker{UID: request.WorkerUID, UserID: uid}}).First(&user).Error; err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, "worker not found", nil)
		return
	}
	var total int64
	if err := db.Model(&models.AccessRule{}).Where(&models.AccessRule{WorkerUID: request.WorkerUID}).Count(&total).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	var accessRules []models.AccessRule
	if err := db.Where(&models.AccessRule{WorkerUID: request.WorkerUID}).Offset((request.Page - 1) * request.PageSize).Limit(request.PageSize).Find(&accessRules).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	common.RespOK(c, common.RespMsgOK, gin.H{
		"total":        total,
		"access_rules": accessRules,
	})
}

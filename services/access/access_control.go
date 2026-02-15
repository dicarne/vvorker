package access

import (
	"vvorker/common"
	"vvorker/entities"
	"vvorker/ext/kv/src/sys_cache"
	"vvorker/models"
	"vvorker/utils"
	"vvorker/utils/database"
	permissions "vvorker/utils/permissions"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type EnableAccessControlRequest struct {
	Enable    *bool  `json:"enable" binding:"required"`
	WorkerUID string `json:"worker_uid" binding:"required"`
}

func UpdateEnableAccessControlEndpoint(c *gin.Context) {

	uid, ok := common.RequireUID(c)
	if !ok {
		return
	}
	request := EnableAccessControlRequest{}
	if err := c.BindJSON(&request); err != nil {
		return
	}

	// 检查用户是否有写权限（拥有者或协作者）
	_, err := permissions.CanWriteWorker(c, uid, request.WorkerUID)
	if err != nil {
		// CanWriteWorker 内部已经调用了 RespErr
		return
	}

	db := database.GetDB()
	var user models.Worker
	if err := db.Where(&models.Worker{Worker: &entities.Worker{UID: request.WorkerUID}}).First(&user).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, "worker not found", nil)
		return
	}
	user.EnableAccessControl = *request.Enable
	if err := db.Save(&user).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	sys_cache.DeleteGlobalCache("worker_uid_name:" + user.Name)
	common.RespOK(c, common.RespMsgOK, nil)
}

type AccessControlRequest struct {
	WorkerUID string `json:"worker_uid" binding:"required"`
}

func GetAccessControlEndpoint(c *gin.Context) {

	uid, ok := common.RequireUID(c)
	if !ok {
		return
	}
	request := AccessControlRequest{}
	if err := c.BindJSON(&request); err != nil {
		return
	}

	// 检查用户是否有读权限（拥有者或协作者）
	_, err := permissions.CanReadWorker(c, uid, request.WorkerUID)
	if err != nil {
		// CanReadWorker 内部已经调用了 RespErr
		return
	}

	db := database.GetDB()
	var user models.Worker
	if err := db.Where(&models.Worker{Worker: &entities.Worker{UID: request.WorkerUID}}).First(&user).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, "worker not found", nil)
		return
	}
	common.RespOK(c, common.RespMsgOK, gin.H{
		"EnableAccessControl": user.EnableAccessControl,
	})
}

func AddAccessRuleEndpoint(c *gin.Context) {

	uid, ok := common.RequireUID(c)
	if !ok {
		return
	}
	request := models.AccessRule{}
	if err := c.BindJSON(&request); err != nil {
		return
	}
	request.Length = len(request.Path)
	request.RuleUID = utils.GenerateUID()

	// 检查用户是否有写权限（拥有者或协作者）
	_, err := permissions.CanWriteWorker(c, uid, request.WorkerUID)
	if err != nil {
		// CanWriteWorker 内部已经调用了 RespErr
		return
	}

	if err := database.GetDB().Create(&request).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	common.RespOK(c, common.RespMsgOK, nil)
}

func UpdateAccessRuleEndpoint(c *gin.Context) {

	uid, ok := common.RequireUID(c)
	if !ok {
		return
	}
	request := models.AccessRule{}
	if err := c.BindJSON(&request); err != nil {
		return
	}
	request.Length = len(request.Path)

	// 检查用户是否有写权限（拥有者或协作者）
	_, err := permissions.CanWriteWorker(c, uid, request.WorkerUID)
	if err != nil {
		// CanWriteWorker 内部已经调用了 RespErr
		return
	}

	db := database.GetDB()
	if err := db.Unscoped().Where(&models.AccessRule{RuleUID: request.RuleUID, WorkerUID: request.WorkerUID}).Delete(&models.AccessRule{}).Error; err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, "rule not found", nil)
		return
	}
	request.RuleUID = utils.GenerateUID()
	request.ID = 0
	if err := db.Create(&request).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	common.RespOK(c, common.RespMsgOK, nil)
}

type DeleteAccessRuleRequest struct {
	WorkerUID string `json:"worker_uid" binding:"required"`
	RuleUID   string `json:"rule_uid" binding:"required"`
}

func DeleteAccessRuleEndpoint(c *gin.Context) {

	uid, ok := common.RequireUID(c)
	if !ok {
		return
	}
	request := DeleteAccessRuleRequest{}
	if err := c.BindJSON(&request); err != nil {
		return
	}

	// 检查用户是否有写权限（拥有者或协作者）
	_, err := permissions.CanWriteWorker(c, uid, request.WorkerUID)
	if err != nil {
		// CanWriteWorker 内部已经调用了 RespErr
		return
	}

	db := database.GetDB()
	if err := db.Delete(&models.AccessRule{}, &models.AccessRule{
		RuleUID:   request.RuleUID,
		WorkerUID: request.WorkerUID,
	}).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	common.RespOK(c, common.RespMsgOK, nil)
}

type ListAccessRuleRequest struct {
	WorkerUID string `json:"worker_uid" binding:"required"`
	Page      int    `json:"page" binding:"gte=1"`
	PageSize  int    `json:"page_size" binding:"gte=1"`
}

func ListAccessRuleEndpoint(c *gin.Context) {

	uid, ok := common.RequireUID(c)
	if !ok {
		return
	}
	request := ListAccessRuleRequest{}
	if err := c.BindJSON(&request); err != nil {
		return
	}

	// 检查用户是否有读权限（拥有者或协作者）
	_, err := permissions.CanReadWorker(c, uid, request.WorkerUID)
	if err != nil {
		// CanReadWorker 内部已经调用了 RespErr
		return
	}

	db := database.GetDB()
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

type SwitchAccessRuleRequest struct {
	WorkerUID string `json:"worker_uid" binding:"required"`
	RuleUID   string `json:"rule_uid" binding:"required"`
	Disable   *bool  `json:"disable" binding:"required"`
}

func SwitchAccessRuleEndpoint(c *gin.Context) {

	uid, ok := common.RequireUID(c)
	if !ok {
		return
	}
	request := SwitchAccessRuleRequest{}
	if err := c.BindJSON(&request); err != nil {
		return
	}

	// 检查用户是否有写权限（拥有者或协作者）
	_, err := permissions.CanWriteWorker(c, uid, request.WorkerUID)
	if err != nil {
		// CanWriteWorker 内部已经调用了 RespErr
		return
	}

	db := database.GetDB()
	stat := 1
	if *request.Disable {
		stat = 2
	}

	if err := db.Model(&models.AccessRule{}).Where(&models.AccessRule{RuleUID: request.RuleUID, WorkerUID: request.WorkerUID}).Update("status", stat).Error; err != nil {
		logrus.Error(err)
		common.RespErr(c, common.RespCodeInvalidRequest, "rule not found", nil)
		return
	}

	common.RespOK(c, common.RespMsgOK, nil)
}

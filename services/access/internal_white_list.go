package access

import (
	"vvorker/common"
	"vvorker/entities"
	"vvorker/models"
	"vvorker/utils/database"
	permissions "vvorker/utils/permissions"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InternalWhiteListCreateRequest 创建内部白名单请求结构体
type InternalWhiteListCreateRequest struct {
	WorkerUID       string `json:"worker_uid" binding:"required"`
	AllowWorkerName string `json:"name" binding:"required"`
	Description     string `json:"description"`
}

// CreateInternalWhiteListEndpoint 创建内部白名单端点
func CreateInternalWhiteListEndpoint(c *gin.Context) {

	uid := uint64(c.GetUint(common.UIDKey))
	request := InternalWhiteListCreateRequest{}
	if err := c.BindJSON(&request); err != nil {
		return
	}
	if uid == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is required", nil)
		return
	}

	// 检查用户是否有写权限（拥有者或协作者）
	_, err := permissions.CanWriteWorker(c, uid, request.WorkerUID)
	if err != nil {
		// CanWriteWorker 内部已经调用了 RespErr
		return
	}

	db := database.GetDB()
	// 查找 allow name 对应的 worker 的 uid
	var allowWorker models.Worker
	if err := db.Where(&models.Worker{
		Worker: &entities.Worker{
			Name: request.AllowWorkerName,
		},
	}).First(&allowWorker).Error; err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, "allowed worker not found", nil)
		return
	}

	whiteList := models.InternalServerWhiteList{
		WorkerUID:      request.WorkerUID,
		AllowWorkerUID: allowWorker.UID,
		Description:    request.Description,
	}
	if err := db.Create(&whiteList).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	common.RespOK(c, common.RespMsgOK, gin.H{
		"internal_white_list": whiteList,
	})
}

// InternalWhiteListListRequest 列出内部白名单请求结构体
type InternalWhiteListListRequest struct {
	WorkerUID string `json:"worker_uid" binding:"required"`
	Page      int    `json:"page" binding:"gte=1"`
	PageSize  int    `json:"page_size" binding:"gte=1"`
}

// ListInternalWhiteListEndpoint 列出内部白名单端点
func ListInternalWhiteListEndpoint(c *gin.Context) {

	uid := uint64(c.GetUint(common.UIDKey))
	request := InternalWhiteListListRequest{}
	if err := c.BindJSON(&request); err != nil {
		return
	}
	if uid == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is required", nil)
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
	if err := db.Model(&models.InternalServerWhiteList{}).Where(&models.InternalServerWhiteList{WorkerUID: request.WorkerUID}).Count(&total).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	// 定义一个包含内部白名单和 worker 名称的结构体
	type WhiteListWithWorkerName struct {
		models.InternalServerWhiteList
		WorkerName string `gorm:"column:name"`
	}
	var whiteLists []WhiteListWithWorkerName

	// 使用 JOIN 操作查询内部白名单和对应的 worker 名称
	if err := db.Table("internal_server_white_lists").
		Select("internal_server_white_lists.*, workers.name").
		Joins("JOIN workers ON internal_server_white_lists.allow_worker_uid = workers.uid").
		Where("internal_server_white_lists.worker_uid = ?", request.WorkerUID).
		Offset((request.Page - 1) * request.PageSize).
		Limit(request.PageSize).
		Find(&whiteLists).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	common.RespOK(c, common.RespMsgOK, gin.H{
		"internal_white_lists": whiteLists,
	})
}

// InternalWhiteListUpdateRequest 更新内部白名单请求结构体
type InternalWhiteListUpdateRequest struct {
	WorkerUID   string `json:"worker_uid" binding:"required"`
	Description string `json:"description"`
	ID          uint   `json:"id" binding:"required,gt=0"`
}

// UpdateInternalWhiteListEndpoint 更新内部白名单端点
func UpdateInternalWhiteListEndpoint(c *gin.Context) {

	uid := uint64(c.GetUint(common.UIDKey))
	request := InternalWhiteListUpdateRequest{}
	if err := c.BindJSON(&request); err != nil {
		return
	}
	if uid == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is required", nil)
		return
	}

	// 检查用户是否有写权限（拥有者或协作者）
	_, err := permissions.CanWriteWorker(c, uid, request.WorkerUID)
	if err != nil {
		// CanWriteWorker 内部已经调用了 RespErr
		return
	}

	db := database.GetDB()
	updateData := map[string]interface{}{
		"description": request.Description,
	}
	if err := db.Model(&models.InternalServerWhiteList{}).Where("worker_uid = ?", request.WorkerUID).Updates(updateData).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	common.RespOK(c, common.RespMsgOK, nil)
}

// InternalWhiteListDeleteRequest 删除内部白名单请求结构体
type InternalWhiteListDeleteRequest struct {
	WorkerUID string `json:"worker_uid" binding:"required"`
	ID        uint   `json:"id" binding:"required,gt=0"`
}

// DeleteInternalWhiteListEndpoint 删除内部白名单端点
func DeleteInternalWhiteListEndpoint(c *gin.Context) {

	uid := uint64(c.GetUint(common.UIDKey))
	request := InternalWhiteListDeleteRequest{}
	if err := c.BindJSON(&request); err != nil {
		return
	}
	if uid == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is required", nil)
		return
	}

	// 检查用户是否有写权限（拥有者或协作者）
	_, err := permissions.CanWriteWorker(c, uid, request.WorkerUID)
	if err != nil {
		// CanWriteWorker 内部已经调用了 RespErr
		return
	}

	db := database.GetDB()
	if err := db.Where(&models.InternalServerWhiteList{WorkerUID: request.WorkerUID, Model: gorm.Model{
		ID: request.ID,
	}}).Delete(&models.InternalServerWhiteList{}).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	common.RespOK(c, common.RespMsgOK, nil)
}

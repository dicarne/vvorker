package access

import (
	"vvorker/common"
	"vvorker/models/secrets"
	"vvorker/utils/database"
	permissions "vvorker/utils/permissions"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SecretCreateRequest 创建密钥请求结构体
type SecretCreateRequest struct {
	WorkerUID string `json:"worker_uid"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

// CreateSecretEndpoint 创建密钥端点
func CreateSecretEndpoint(c *gin.Context) {

	uid := uint64(c.GetUint(common.UIDKey))
	request := SecretCreateRequest{}
	if err := c.BindJSON(&request); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}
	if uid == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is required", nil)
		return
	}
	if request.WorkerUID == "" {
		common.RespErr(c, common.RespCodeInvalidRequest, "worker_uid is required", nil)
		return
	}
	if request.Key == "" {
		common.RespErr(c, common.RespCodeInvalidRequest, "key is required", nil)
		return
	}

	// 检查用户是否有写权限（拥有者或协作者）
	_, err := permissions.CanWriteWorker(c, uid, request.WorkerUID)
	if err != nil {
		// CanWriteWorker 内部已经调用了 RespErr
		return
	}

	db := database.GetDB()
	secret := secrets.Secret{
		WorkerUID: request.WorkerUID,
		Key:       request.Key,
		Value:     request.Value,
	}
	if err := db.Create(&secret).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	common.RespOK(c, common.RespMsgOK, gin.H{
		"secret": secret,
	})
}

// SecretListRequest 列出密钥请求结构体
type SecretListRequest struct {
	WorkerUID string `json:"worker_uid"`
}

// ListSecretEndpoint 列出密钥端点
func ListSecretEndpoint(c *gin.Context) {

	uid := uint64(c.GetUint(common.UIDKey))
	request := SecretListRequest{}
	if err := c.BindJSON(&request); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}
	if uid == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is required", nil)
		return
	}
	if request.WorkerUID == "" {
		common.RespErr(c, common.RespCodeInvalidRequest, "worker_uid is required", nil)
		return
	}

	// 检查用户是否有读权限（拥有者或协作者）
	_, err := permissions.CanReadWorker(c, uid, request.WorkerUID)
	if err != nil {
		// CanReadWorker 内部已经调用了 RespErr
		return
	}

	db := database.GetDB()
	var secretList []secrets.Secret
	if err := db.Where(&secrets.Secret{WorkerUID: request.WorkerUID}).Find(&secretList).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	// 隐藏密钥值，统一显示为六个星号
	// for i := range secretList {
	// 	secretList[i].Value = "******"
	// }

	common.RespOK(c, common.RespMsgOK, gin.H{
		"secrets": secretList,
	})
}

// SecretUpdateRequest 更新密钥请求结构体
type SecretUpdateRequest struct {
	WorkerUID string `json:"worker_uid"`
	ID        uint   `json:"id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

// UpdateSecretEndpoint 更新密钥端点
func UpdateSecretEndpoint(c *gin.Context) {

	uid := uint64(c.GetUint(common.UIDKey))
	request := SecretUpdateRequest{}
	if err := c.BindJSON(&request); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}
	if uid == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is required", nil)
		return
	}
	if request.WorkerUID == "" {
		common.RespErr(c, common.RespCodeInvalidRequest, "worker_uid is required", nil)
		return
	}
	if request.ID == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "id is required", nil)
		return
	}

	// 检查用户是否有写权限（拥有者或协作者）
	_, err := permissions.CanWriteWorker(c, uid, request.WorkerUID)
	if err != nil {
		// CanWriteWorker 内部已经调用了 RespErr
		return
	}

	db := database.GetDB()
	updateData := map[string]interface{}{}
	if request.Key != "" {
		updateData["key"] = request.Key
	}
	if request.Value != "" {
		updateData["value"] = request.Value
	}

	if len(updateData) > 0 {
		if err := db.Model(&secrets.Secret{}).
			Where("worker_uid = ? AND id = ?", request.WorkerUID, request.ID).
			Updates(updateData).Error; err != nil {
			common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
			return
		}
	}

	common.RespOK(c, common.RespMsgOK, nil)
}

// SecretDeleteRequest 删除密钥请求结构体
type SecretDeleteRequest struct {
	WorkerUID string `json:"worker_uid"`
	ID        uint   `json:"id"`
}

// DeleteSecretEndpoint 删除密钥端点
func DeleteSecretEndpoint(c *gin.Context) {

	uid := uint64(c.GetUint(common.UIDKey))
	request := SecretDeleteRequest{}
	if err := c.BindJSON(&request); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}
	if uid == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is required", nil)
		return
	}
	if request.WorkerUID == "" {
		common.RespErr(c, common.RespCodeInvalidRequest, "worker_uid is required", nil)
		return
	}
	if request.ID == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "id is required", nil)
		return
	}

	// 检查用户是否有写权限（拥有者或协作者）
	_, err := permissions.CanWriteWorker(c, uid, request.WorkerUID)
	if err != nil {
		// CanWriteWorker 内部已经调用了 RespErr
		return
	}

	db := database.GetDB()
	if err := db.Where(&secrets.Secret{WorkerUID: request.WorkerUID, Model: gorm.Model{
		ID: request.ID,
	}}).Delete(&secrets.Secret{}).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	common.RespOK(c, common.RespMsgOK, nil)
}

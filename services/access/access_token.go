package access

import (
	"vvorker/common"
	"vvorker/models"
	"vvorker/utils"
	"vvorker/utils/database"
	"vvorker/utils/permissions"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccessTokenCreateRequest struct {
	WorkerUID      string `json:"worker_uid" binding:"required"`
	Description    string `json:"description"`
	Forever        bool   `json:"forever"`
	ExpirationTime string `json:"expiration_time"`
}

func CreateAccessTokenEndpoint(c *gin.Context) {

	uid, ok := common.RequireUID(c)
	if !ok {
		return
	}
	request := AccessTokenCreateRequest{}
	if err := c.BindJSON(&request); err != nil {
		return
	}

	_, err := permissions.CanWriteWorker(c, uid, request.WorkerUID)
	if err != nil {
		return
	}

	db := database.GetDB()
	accessToken := models.ExternalServerToken{
		WorkerUID:      request.WorkerUID,
		Description:    request.Description,
		Forever:        request.Forever,
		ExpirationTime: request.ExpirationTime,
		Token:          utils.GenerateUID(),
	}
	if err := db.Create(&accessToken).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	common.RespOK(c, common.RespMsgOK, gin.H{
		"access_token": accessToken.Token,
	})
}

type AccessTokenListRequest struct {
	WorkerUID string `json:"worker_uid" binding:"required"`
	Page      int    `json:"page" binding:"gte=1"`
	PageSize  int    `json:"page_size" binding:"gte=1"`
}

func ListAccessTokenEndpoint(c *gin.Context) {

	uid, ok := common.RequireUID(c)
	if !ok {
		return
	}
	request := AccessTokenListRequest{}
	if err := c.BindJSON(&request); err != nil {
		return
	}

	_, err := permissions.CanReadWorker(c, uid, request.WorkerUID)
	if err != nil {
		return
	}

	db := database.GetDB()
	var total int64
	if err := db.Model(&models.ExternalServerToken{}).Where(&models.ExternalServerToken{WorkerUID: request.WorkerUID}).Count(&total).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	var accessTokens []models.ExternalServerToken
	if err := db.Where(&models.ExternalServerToken{WorkerUID: request.WorkerUID}).Offset((request.Page - 1) * request.PageSize).Limit(request.PageSize).Find(&accessTokens).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	// 遍历每个key，只保留前三位，其他用*替代
	for i := range accessTokens {
		accessTokens[i].Token = accessTokens[i].Token[:3] + "************" + accessTokens[i].Token[len(accessTokens[i].Token)-3:]
	}
	common.RespOK(c, common.RespMsgOK, gin.H{
		"access_tokens": accessTokens,
	})
}

type AccessTokenDeleteRequest struct {
	WorkerUID string `json:"worker_uid" binding:"required"`
	ID        uint   `json:"id" binding:"required,gt=0"`
}

func DeleteAccessTokenEndpoint(c *gin.Context) {

	uid, ok := common.RequireUID(c)
	if !ok {
		return
	}
	request := AccessTokenDeleteRequest{}
	if err := c.BindJSON(&request); err != nil {
		return
	}

	_, err := permissions.CanWriteWorker(c, uid, request.WorkerUID)
	if err != nil {
		return
	}

	db := database.GetDB()
	if err := db.Where(&models.ExternalServerToken{WorkerUID: request.WorkerUID, Model: gorm.Model{
		ID: request.ID,
	}}).Delete(&models.ExternalServerToken{}).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	common.RespOK(c, common.RespMsgOK, nil)
}

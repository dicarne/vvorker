package access

import (
	"runtime/debug"
	"vvorker/common"
	"vvorker/entities"
	"vvorker/models"
	"vvorker/utils"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AccessTokenCreateRequest struct {
	WorkerUID      string `json:"worker_uid"`
	Description    string `json:"description"`
	Forever        bool   `json:"forever"`
	ExpirationTime string `json:"expiration_time"`
}

func CreateAccessTokenEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	uid := uint64(c.GetUint(common.UIDKey))
	request := AccessTokenCreateRequest{}
	if err := c.BindJSON(&request); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}
	if uid == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is required", nil)
		return
	}
	db := database.GetDB()
	worker := models.Worker{}

	if err := db.Where(&models.Worker{
		Worker: &entities.Worker{
			UID:    request.WorkerUID,
			UserID: uid,
		},
	}).First(&worker).Error; err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, "worker not found", nil)
		return
	}

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
		"access_token": accessToken,
	})
}

type AccessTokenListRequest struct {
	WorkerUID string `json:"worker_uid"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
}

func ListAccessTokenEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	uid := uint64(c.GetUint(common.UIDKey))
	request := AccessTokenListRequest{}
	if err := c.BindJSON(&request); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}
	if uid == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is required", nil)
		return
	}
	db := database.GetDB()
	worker := models.Worker{}

	if err := db.Where(&models.Worker{
		Worker: &entities.Worker{
			UID:    request.WorkerUID,
			UserID: uid,
		},
	}).First(&worker).Error; err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, "worker not found", nil)
		return
	}
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
	WorkerUID string `json:"worker_uid"`
}

func DeleteAccessTokenEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()

	uid := uint64(c.GetUint(common.UIDKey))
	request := AccessTokenDeleteRequest{}
	if err := c.BindJSON(&request); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}
	if uid == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is required", nil)
		return
	}
	db := database.GetDB()
	worker := models.Worker{}

	if err := db.Where(&models.Worker{
		Worker: &entities.Worker{
			UID:    request.WorkerUID,
			UserID: uid,
		},
	}).First(&worker).Error; err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, "worker not found", nil)
		return
	}
	if err := db.Where(&models.ExternalServerToken{WorkerUID: request.WorkerUID}).Delete(&models.ExternalServerToken{}).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	common.RespOK(c, common.RespMsgOK, nil)
}

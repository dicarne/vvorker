package access

import (
	"net/http"
	"runtime/debug"
	"vvorker/common"
	"vvorker/models"
	"vvorker/utils"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AccessKeyCreateRequest struct {
	Name string `json:"name"`
}

type AccessKeyCreateResponse struct {
	AccessKey string `json:"key"`
	Name      string `json:"name"`
}

type AccessKeyDeleteRequest struct {
	AccessKey string `json:"key"`
}

func CreateAccessKeyEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	uid := uint64(c.GetUint(common.UIDKey))

	request := AccessKeyCreateRequest{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	if uid == 0 {
		c.JSON(400, gin.H{"code": 1, "msg": "uid is required"})
		return
	}
	db := database.GetDB()
	accessKey := models.AccessKey{
		UserId: uid,
		Name:   request.Name,
		Key:    "ac::" + utils.GenerateUID(),
	}
	if err := db.Create(&accessKey).Error; err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Create Access Key Failed.", gin.H{"error": err.Error()})
		return
	}
	common.RespOK(c, "success", AccessKeyCreateResponse{AccessKey: accessKey.Key, Name: accessKey.Name})
}

func AccessKeyToUserID(accessKey string) (uint64, error) {
	db := database.GetDB()
	var accessKeyModel models.AccessKey
	if err := db.Where("key = ?", accessKey).First(&accessKeyModel).Error; err != nil {
		return 0, err
	}
	return accessKeyModel.UserId, nil
}

func GetAccessKeysEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	uid := uint64(c.GetUint(common.UIDKey))
	db := database.GetDB()
	var accessKeys []models.AccessKey
	if err := db.Where("user_id = ?", uid).Find(&accessKeys).Error; err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Get Access Keys Failed.", gin.H{"error": err.Error()})
		return
	}
	common.RespOK(c, "success", accessKeys)
}

func DeleteAccessKeyEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	uid := uint64(c.GetUint(common.UIDKey))
	request := AccessKeyDeleteRequest{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	db := database.GetDB()
	var accessKeyModel models.AccessKey
	if err := db.Where("user_id =?", uid).Where("key =?", request.AccessKey).First(&accessKeyModel).Error; err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Delete Access Key Failed.", gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&accessKeyModel).Error; err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Delete Access Key Failed.", gin.H{"error": err.Error()})
		return
	}
	common.RespOK(c, "success", nil)
}

package access

import (
	"net/http"
	"vvorker/common"
	"vvorker/models"
	"vvorker/utils"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
)

type AccessKeyCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type AccessKeyCreateResponse struct {
	AccessKey string `json:"key"`
	Name      string `json:"name"`
}

type AccessKeyDeleteRequest struct {
	AccessKey string `json:"key" binding:"required"`
}

func CreateAccessKeyEndpoint(c *gin.Context) {

	uid, ok := common.RequireUID(c)
	if !ok {
		return
	}

	request := AccessKeyCreateRequest{}
	if err := c.BindJSON(&request); err != nil {
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
	if err := db.Where(&models.AccessKey{
		Key: accessKey,
	}).First(&accessKeyModel).Error; err != nil {
		return 0, err
	}
	return accessKeyModel.UserId, nil
}

func GetAccessKeysEndpoint(c *gin.Context) {

	uid, _ := common.RequireUID(c)
	db := database.GetDB()
	var accessKeys []models.AccessKey
	if err := db.Where(&models.AccessKey{
		UserId: uid,
	}).Find(&accessKeys).Error; err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Get Access Keys Failed.", gin.H{"error": err.Error()})
		return
	}
	common.RespOK(c, "success", accessKeys)
}

func DeleteAccessKeyEndpoint(c *gin.Context) {

	uid, _ := common.RequireUID(c)
	request := AccessKeyDeleteRequest{}
	if err := c.BindJSON(&request); err != nil {
		return
	}
	db := database.GetDB()
	var accessKeyModel models.AccessKey
	if err := db.Where(&models.AccessKey{
		UserId: uid,
		Key:    request.AccessKey,
	}).First(&accessKeyModel).Error; err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Delete Access Key Failed.", gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&accessKeyModel).Error; err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Delete Access Key Failed.", gin.H{"error": err.Error()})
		return
	}
	common.RespOK(c, "success", nil)
}

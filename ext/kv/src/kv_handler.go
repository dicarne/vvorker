package kv

import (
	"net/http"
	"strconv"
	"vorker/common"
	"vorker/entities"
	"vorker/models"
	"vorker/utils"
	"vorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CreateKVResourcesEndpoint(c *gin.Context) {
	var req = entities.CreateNewResourcesRequest{}
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusBadRequest, "invalid request", gin.H{"error": err.Error()})
		return
	}
	if !req.Validate() {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusBadRequest, "invalid request", gin.H{"error": "invalid request"})
		return
	}
	db := database.GetDB()
	userID, err := strconv.ParseUint(req.UserID, 10, 64)
	if err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusBadRequest, "Failed to convert UserID to uint64", gin.H{"error": err.Error()})
		return
	}
	kvResource := &models.KV{
		UserID: userID,
		Name:   req.Name,
		UID:    utils.GenerateUID(),
	}
	if err := db.Create(kvResource).Error; err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusInternalServerError, "Failed to create KV resource", gin.H{"error": err.Error()})
		return
	}
	// 使用 common.RespOK 返回成功响应
	common.RespOK(c, "KV resource created successfully", gin.H{"uid": kvResource.UID, "status": 0})
}

// 删除指定KV资源
func DeleteKVResourcesEndpoint(c *gin.Context) {
	uid := c.GetUint64(common.UIDKey)

	var req = entities.DeleteResourcesReq{}
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusBadRequest, "invalid request", gin.H{"error": err.Error()})
		return
	}
	if !req.Validate() {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusBadRequest, "invalid request", gin.H{"error": "invalid request"})
		return
	}

	condition := models.PostgreSQL{UID: req.UID, UserID: uid}

	db := database.GetDB()

	if rr := db.Delete(&condition, condition); rr.Error != nil || rr.RowsAffected == 0 {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusInternalServerError, "Failed to delete KV resource", gin.H{"error": rr.Error.Error()})
		return
	}
	// 使用 common.RespOK 返回成功响应
	common.RespOK(c, "KV resource deleted successfully", gin.H{"status": 0})
}

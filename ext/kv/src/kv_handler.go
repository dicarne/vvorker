package kv

import (
	"net/http"
	"vvorker/common"
	"vvorker/entities"
	"vvorker/models"
	"vvorker/utils"
	"vvorker/utils/database"

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
	userID := uint64(c.GetUint(common.UIDKey))
	if userID == 0 {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusBadRequest, "Failed to convert UserID to uint64", gin.H{"error": "uid is required"})
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
	common.RespOK(c, "success", entities.CreateNewResourcesResponse{
		UID:  kvResource.UID,
		Name: kvResource.Name,
		Type: "kv",
	})
}

// 删除指定KV资源
func DeleteKVResourcesEndpoint(c *gin.Context) {
	uid := uint64(c.GetUint(common.UIDKey))

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

	condition := models.KV{UID: req.UID, UserID: uid}

	db := database.GetDB()

	if rr := db.Delete(&condition, condition); rr.Error != nil || rr.RowsAffected == 0 {
		// 使用 common.RespErr 返回错误响应
		msg := ""
		if rr.Error != nil {
			msg = rr.Error.Error()
		} else if rr.RowsAffected == 0 {
			msg = "resource not found"
		}
		common.RespErr(c, http.StatusInternalServerError, "Failed to delete KV resource", gin.H{"error": msg})
		return
	}
	// 使用 common.RespOK 返回成功响应
	common.RespOK(c, "success", entities.DeleteResourcesResp{
		Status: 0,
	})
}

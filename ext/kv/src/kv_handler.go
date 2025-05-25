package kv

import (
	"net/http"
	"strconv"
	"vorker/entities"
	"vorker/models"
	"vorker/utils"
	"vorker/utils/database"

	"github.com/gin-gonic/gin"
)

func CreateKVResourcesEndpoint(c *gin.Context) {
	var req = entities.CreateNewResourcesRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !req.Validate() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	db := database.GetDB()
	userID, err := strconv.ParseUint(req.UserID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to convert UserID to uint64: " + err.Error()})
		return
	}
	kvResource := &models.KV{
		UserID: userID,
		Name:   req.Name,
		UID:    utils.GenerateUID(),
	}
	if err := db.Create(kvResource).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create KV resource: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"uid": kvResource.UID, "status": 0})
}

// 删除指定KV资源
func DeleteKVResourcesEndpoint(c *gin.Context) {
	var req = entities.DeleteResourcesReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !req.Validate() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	db := database.GetDB()
	db.Model(&models.KV{}).Delete(&models.KV{
		UID: req.UID,
	})
	c.JSON(http.StatusOK, gin.H{"status": 0})
}

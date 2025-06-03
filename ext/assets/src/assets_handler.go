package assets

import (
	"vvorker/models"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
)

type UploadAssetsReq struct {
	UID       string `json:"uid"`
	WorkerUID string `json:"worker_uid"`
	Path      string `json:"path"`
}

func UploadAssetsEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
		}
	}()
	var req UploadAssetsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	if req.UID == "" || req.Path == "" {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	userID := c.GetUint("uid")
	if userID == 0 {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var w models.Worker
	if err := database.GetDB().Where("uid =?", req.WorkerUID).First(&w).Error; err != nil {
		c.JSON(404, gin.H{"error": "Worker not found"})
		return
	}

	if w.UserID != uint64(userID) {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	db := database.GetDB()

	var file models.File
	if err := db.Where("uid = ?", req.UID).First(&file).Error; err != nil {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}
	var asset = models.Assets{
		UID:       req.UID,
		WorkerUID: req.WorkerUID,
		Path:      req.Path,
		UserID:    uint64(userID),
		MIME:      file.Mimetype,
		Hash:      file.Hash,
	}
	nkv := models.Assets{}
	// 如果有，则更新，如果无，则新增
	if err := db.Where("uid = ?", req.UID).Assign(&asset).FirstOrCreate(&nkv).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to save asset"})
		return
	}
	c.JSON(200, gin.H{"message": "Asset saved successfully"})
}

type GetAssetsReq struct {
	WorkerUID string `json:"worker_uid"`
	Path      string `json:"path"`
}

func GetAssetsEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
		}
	}()

	var req GetAssetsReq
	req.Path = c.GetHeader("vvorker-asset-path")
	req.WorkerUID = c.GetHeader("vvorker-asset-worker-uid")

	if req.WorkerUID == "" || req.Path == "" {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	db := database.GetDB()
	var asset models.Assets
	if err := db.Where("worker_uid = ? AND path = ?", req.WorkerUID, req.Path).First(&asset).Error; err != nil {
		c.JSON(404, gin.H{"error": "Asset not found"})
		return
	}

	var file models.File
	if err := db.Where("uid =?", asset.UID).First(&file).Error; err != nil {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}

	c.Data(200, file.Mimetype, file.Data)
}

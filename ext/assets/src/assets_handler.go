package assets

import (
	"mime"
	"vvorker/entities"
	"vvorker/models"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	if err := database.GetDB().Where(&models.Worker{
		Worker: &entities.Worker{
			UID: req.WorkerUID,
		},
	}).First(&w).Error; err != nil {
		logrus.Errorf("Worker not found: %v", err)
		c.JSON(404, gin.H{"error": "Worker not found"})
		return
	}

	if w.UserID != uint64(userID) {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	db := database.GetDB()

	var file models.File
	if err := db.Where(&models.File{
		UID: req.UID,
	}).First(&file).Error; err != nil {
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
	if err := db.Where(&models.Assets{
		Path:      req.Path,
		WorkerUID: req.WorkerUID,
	}).Assign(&asset).FirstOrCreate(&nkv).Error; err != nil {
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
	if err := db.Where(&models.Assets{
		WorkerUID: req.WorkerUID,
		Path:      req.Path,
	}).First(&asset).Error; err != nil {
		c.JSON(404, gin.H{"error": "Asset not found"})
		return
	}

	var file models.File
	if err := db.Where(&models.File{
		UID: asset.UID,
	}).First(&file).Error; err != nil {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}

	mimeType := mime.TypeByExtension(file.Mimetype)
	if mimeType == "" {
		// 如果没有匹配的 MIME 类型，默认使用 application/octet-stream
		mimeType = "application/octet-stream"
	}

	c.Data(200, mimeType, file.Data)
}

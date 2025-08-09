package assets

import (
	"mime"
	"vvorker/entities"
	"vvorker/ext/kv/src/sys_cache"
	"vvorker/models"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var AssetBucket = "assets"

type ClearAssetsReq struct {
	WorkerUID string `json:"worker_uid"`
}

func ClearAssetsEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
		}
	}()
	var req = ClearAssetsReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	if req.WorkerUID == "" {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	userID := c.GetUint("uid")
	if userID == 0 {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	db := database.GetDB()

	var w models.Worker
	if err := db.Where(&models.Worker{
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

	assets := []models.Assets{}
	if err := db.Where(&models.Assets{
		WorkerUID: req.WorkerUID,
	}).Find(&assets).Error; err != nil {
		logrus.Errorf("Failed to find assets: %v", err)
		c.JSON(500, gin.H{"error": "Failed to find assets"})
		return
	}

	if err := db.Where(&models.Assets{
		WorkerUID: req.WorkerUID,
	}).Delete(&assets).Error; err != nil {
		logrus.Errorf("Failed to delete assets: %v", err)
		c.JSON(500, gin.H{"error": "Failed to delete assets"})
		return
	}

	deleteCount := 0
	for _, a := range assets {
		// 如果Assets 中没有任何资源引用了a.UID
		count := int64(0)
		if err := db.Model(&models.Assets{}).Where(&models.Assets{
			UID: a.UID,
		}).Count(&count).Error; err != nil {
			logrus.Errorf("Failed to count assets: %v", err)
			c.JSON(500, gin.H{"error": "Failed to count assets"})
			return
		}
		if count == 0 {
			// 则删除File中UID为a.UID的文件
			if err := db.Unscoped().Where(&models.File{
				UID: a.UID,
			}).Delete(&models.File{}).Error; err != nil {
				logrus.Errorf("Failed to delete file: %v", err)
				c.JSON(500, gin.H{"error": "Failed to delete file"})
			}
			deleteCount++
		}
	}

	c.JSON(200, gin.H{"message": "Assets cleared successfully", "delete_count": deleteCount})
}

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
		c.JSON(404, gin.H{"error": "Worker not found", "worker_uid": req.WorkerUID})
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

	cache, err := sys_cache.Get(AssetBucket + ":" + asset.UID + ":data")
	if len(cache) == 0 || err != nil {
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
		sys_cache.Put(AssetBucket+":"+asset.UID+":data", file.Data, 3600)
		sys_cache.Put(AssetBucket+":"+asset.UID+":mime", []byte(mimeType), 3600)
		c.Data(200, mimeType, file.Data)
	} else {
		bmimeType, _ := sys_cache.Get(AssetBucket + ":" + asset.UID + ":mime")
		mimeType := string(bmimeType)
		if mimeType == "" {
			// 如果没有匹配的 MIME 类型，默认使用 application/octet-stream
			mimeType = "application/octet-stream"
		}
		c.Data(200, mimeType, cache)
	}

}

package kv

import (
	"errors"
	"net/http"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/defs"
	"vvorker/entities"
	"vvorker/models"
	"vvorker/utils"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nutsdb/nutsdb"
	"github.com/sirupsen/logrus"
)

var db *nutsdb.DB

var buckets *defs.SyncMap[string, bool]

var SystemBucket = "system"

func init() {
	db2, err := nutsdb.Open(
		nutsdb.DefaultOptions,
		nutsdb.WithDir(conf.AppConfigInstance.LocalKVDir), // 数据库会自动创建这个目录文件
	)
	db = db2
	if err != nil {
		logrus.Panic(err)
	}
	buckets = defs.NewSyncMap(map[string]bool{})
}

func Close() {
	if db != nil {
		db.Close()
	}
}

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
	userID := uint64(c.GetUint(common.UIDKey))
	if userID == 0 {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusBadRequest, "Failed to convert UserID to uint64", gin.H{"error": "uid is required"})
		return
	}
	kvResource, err := CreateKV(userID, req.Name)
	if err != nil {
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

func CreateKV(userID uint64, name string) (*models.KV, error) {
	kvModel := &models.KV{}
	kvModel.UserID = userID
	kvModel.Name = name
	kvModel.UID = utils.GenerateUID()
	if err := database.GetDB().Create(kvModel).Error; err != nil {
		return nil, err
	}
	return kvModel, nil
}

func RecoverKV(userID uint64, kv *models.KV) error {
	kv.UserID = userID
	db := database.GetDB()
	nkv := models.KV{}
	// 如果有，则更新，如果无，则新增
	if err := db.Where("uid = ?", kv.UID).Assign(kv).FirstOrCreate(&nkv).Error; err != nil {
		return err
	}
	return nil
}

// 删除指定KV资源
func DeleteKVResourcesEndpoint(c *gin.Context) {
	uid := uint64(c.GetUint(common.UIDKey))

	var req = entities.DeleteResourcesReq{}
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		common.RespErr(c, http.StatusBadRequest, "invalid request", gin.H{"error": err.Error()})
		return
	}
	if !req.Validate() {
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

type InvokeKVRequest struct {
	RID    string `json:"rid"`
	Key    string `json:"key"`
	Value  string `json:"value"`
	Method string `json:"method"`
	TTL    int    `json:"ttl"`
	Offset int    `json:"offset"`
	Size   int    `json:"size"`
}

func InvokeKVEndpoint(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			common.RespErr(c, http.StatusInternalServerError, "Failed to invoke KV resource", gin.H{"error": err})
		}
	}()
	var req = InvokeKVRequest{}
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		common.RespErr(c, http.StatusBadRequest, "invalid request", gin.H{"error": err.Error()})
		return
	}

	switch req.Method {
	case "get":
		{
			value, err := Get(req.RID, req.Key)
			if err != nil {
				common.RespErr(c, http.StatusInternalServerError, "Failed to get KV resource", gin.H{"error": err.Error()})
				return
			}
			common.RespOK(c, "success", string(value))
		}
	case "set":
		{
			if err := Put(req.RID, req.Key, []byte(req.Value), req.TTL); err != nil {
				common.RespErr(c, http.StatusInternalServerError, "Failed to set KV resource", gin.H{"error": err.Error()})
				return
			}
			common.RespOK(c, "success", nil)
		}
	case "del":
		{
			if err := Del(req.RID, req.Key); err != nil {
				common.RespErr(c, http.StatusInternalServerError, "Failed to delete KV resource", gin.H{"error": err.Error()})
				return
			}
			common.RespOK(c, "success", nil)
		}
	case "keys":
		{
			keys, err := Keys(req.RID, req.Key, req.Offset, req.Size)
			if err != nil {
				common.RespErr(c, http.StatusInternalServerError, "Failed to get KV resource", gin.H{"error": err.Error()})
				return
			}
			common.RespOK(c, "success", keys)
		}
	default:
		common.RespErr(c, http.StatusBadRequest, "invalid request", gin.H{"error": "invalid request"})
		return
	}
}

func ExistBucket(bucket string) error {
	if _, exist := buckets.Get(bucket); exist {
		return nil
	}
	return db.Update(func(tx *nutsdb.Tx) error {
		tx.NewKVBucket(bucket)
		buckets.Set(bucket, true)
		return nil
	})
}

func Put(bucket string, key string, value []byte, ttl int) error {
	ExistBucket(bucket)
	return db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put(bucket, []byte(key), value, uint32(ttl))
	})
}

func Get(bucket string, key string) ([]byte, error) {
	ExistBucket(bucket)
	var value []byte
	err := db.View(func(tx *nutsdb.Tx) error {
		v, err := tx.Get(bucket, []byte(key))
		if err != nil {
			return err
		}
		value = v
		return nil
	})
	if err != nil {
		if errors.Is(err, nutsdb.ErrKeyNotFound) {
			return []byte(""), nil
		}
		return nil, err
	}
	return value, nil
}

func Del(bucket string, key string) error {
	ExistBucket(bucket)
	return db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete(bucket, []byte(key))
	})
}

func Keys(bucket string, prefix string, offset int, size int) ([]string, error) {
	ExistBucket(bucket)
	var result []string
	err := db.View(
		func(tx *nutsdb.Tx) error {
			prefixBytes := []byte(prefix)
			// Based on compiler feedback, tx.PrefixScan appears to return a slice of keys ([]byte),
			// not a slice of Entry structs.
			entries, err := tx.PrefixScan(bucket, prefixBytes, offset, size)
			if err != nil {
				return err
			}
			result = make([]string, 0, len(entries))
			for _, entry := range entries {
				result = append(result, string(entry))
			}
			return nil
		})

	if err != nil {
		return nil, err
	}
	return result, nil
}

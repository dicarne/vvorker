package kv

import (
	"net/http"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/entities"
	_ "vvorker/ext/kv/src/kv_nutsdb"
	kvnutsdb "vvorker/ext/kv/src/kv_nutsdb"
	kvredis "vvorker/ext/kv/src/kv_redis"
	kvtypes "vvorker/ext/kv/src/kv_types"
	"vvorker/models"
	"vvorker/utils"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var kvStorage kvtypes.IKVStorage

func init() {
	if conf.AppConfigInstance.KVProvider == "redis" {
		kvStorage = &kvredis.KVRedis{}
	} else {
		kvStorage = &kvnutsdb.KVNutsDB{}
	}
}

func Close() {
	kvStorage.Close()
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
	userID, ok := common.RequireUID(c)
	if !ok {
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
	uid, _ := common.RequireUID(c)

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

/////////////////////

func InvokeKVEndpoint(c *gin.Context) {
	var req = kvtypes.InvokeKVRequest{}
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		common.RespErr(c, http.StatusBadRequest, "invalid request", gin.H{"error": err.Error()})
		return
	}

	switch req.Method {
	case "get":
		{
			value, err := kvStorage.Get(req.RID, req.Key)
			if err != nil {
				common.RespErr(c, http.StatusInternalServerError, "Failed to get KV resource", gin.H{"error": err.Error()})
				return
			}
			if value != nil {
				common.RespOK(c, "success", string(value))
			} else {
				common.RespOK(c, "success", nil)
			}
		}
	case "set":
		{
			if req.Options.NX {
				code, err := kvStorage.PutNX(req.RID, req.Key, []byte(req.Value), req.Options.EX)
				if code != 0 {
					common.RespOK(c, "success", code)
					return
				}
				if err != nil {
					common.RespErr(c, http.StatusInternalServerError, "Failed to setNX KV resource", gin.H{"error": err.Error()})
					return
				}

			} else if req.Options.XX {
				code, err := kvStorage.PutXX(req.RID, req.Key, []byte(req.Value), req.Options.EX)
				if code != 0 {
					common.RespOK(c, "success", code)
					return
				}
				if err != nil {
					common.RespErr(c, http.StatusInternalServerError, "Failed to setXX KV resource", gin.H{"error": err.Error()})
					return
				}

			} else {
				code, err := kvStorage.Put(req.RID, req.Key, []byte(req.Value), req.Options.EX)
				if code != 0 {
					common.RespOK(c, "success", code)
					return
				}
				if err != nil {
					common.RespErr(c, http.StatusInternalServerError, "Failed to setEX KV resource", gin.H{"error": err.Error()})
					return
				}

			}
			common.RespOK(c, "success", nil)
		}
	case "del":
		{
			if err := kvStorage.Del(req.RID, req.Key); err != nil {
				common.RespErr(c, http.StatusInternalServerError, "Failed to delete KV resource", gin.H{"error": err.Error()})
				return
			}
			common.RespOK(c, "success", nil)
		}
	case "keys":
		{
			keys, err := kvStorage.Keys(req.RID, req.Key, req.Offset, req.Size)
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

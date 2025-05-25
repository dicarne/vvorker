package workerd

import (
	"encoding/json"
	"runtime/debug"
	"strconv"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/defs"
	"vvorker/entities"
	"vvorker/models"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetWorkersEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	offsetStr := c.Param("offset")
	if len(offsetStr) == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "offset is empty", nil)
		return
	}
	limitStr := c.Param("limit")
	if len(limitStr) == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "limit is empty", nil)
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "offset is invalid", nil)
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "limit is invalid", nil)
		return
	}
	userID := c.GetUint(common.UIDKey)

	workers, err := models.GetWorkers(userID, offset, limit)
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	common.RespOK(c, "get worker success", models.Trans2Entities(workers))
}

func GetAllWorkersEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	userID := c.GetUint(common.UIDKey)
	workers, err := models.GetAllWorkers(userID)
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	common.RespOK(c, "get all workers success", models.Trans2Entities(workers))
}

func GetWorkerEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	userID := c.GetUint(common.UIDKey)
	uid := c.Param("uid")
	if len(uid) == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is empty", nil)
		return
	}
	worker, err := models.GetWorkerByUID(userID, uid)
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	common.RespOK(c, "get workers success", models.Trans2Entities([]*models.Worker{worker}))
}

func AgentSyncWorkers(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	req := &entities.AgentSyncWorkersReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.RespErr(c, defs.CodeInvalidRequest, err.Error(), nil)
		return
	}

	nodeName := c.GetString(defs.KeyNodeName)
	// get node's workerlist
	workers, err := models.AdminGetWorkersByNodeName(nodeName)
	if err != nil {
		common.RespErr(c, defs.CodeInternalError, err.Error(), nil)
		return
	}

	// build response
	// TODO: chunk loading
	resp := &entities.AgentSyncWorkersResp{
		WorkerList: &entities.WorkerList{
			NodeName: nodeName,
			Workers:  models.Trans2Entities(workers),
		},
	}
	common.RespOK(c, "sync workers success", resp)
}

func FillWorkerConfig(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	req := &entities.AgentFillWorkerReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.RespErr(c, defs.CodeInvalidRequest, err.Error(), nil)
		return
	}

	db := database.GetDB()
	worker := &models.Worker{}
	con := db.Model(&models.Worker{}).Where(&models.Worker{Worker: &entities.Worker{
		UID: req.UID,
	}}).First(worker)
	if con.Error != nil {
		common.RespErr(c, defs.CodeInternalError, con.Error.Error(), nil)
		return
	}
	newTemplate := FinishWorkerConfig(worker)

	common.RespOK(c, "fill worker config success", &entities.AgentFillWorkerResp{
		NewTemplate: newTemplate,
	})
}

func FinishWorkerConfig(worker *models.Worker) string {
	UserID := worker.UserID
	workerconfig, err := conf.ParseWorkerConfig(worker.Template)
	if err == nil {
		db := database.GetDB()
		for i, ext := range workerconfig.PgSql {
			if len(ext.ResourceID) != 0 {
				var pgresources = models.PostgreSQL{}
				db.Model(&models.PostgreSQL{}).Where(&models.PostgreSQL{UID: ext.ResourceID, UserID: uint64(UserID)}).First(&pgresources)
				if pgresources.ID != 0 {
					ext.Database = pgresources.Database
					ext.Password = pgresources.Password
					ext.User = pgresources.Username
				} else {
					ext.ResourceID = ""
				}
				workerconfig.PgSql[i] = ext
			}
		}

		for i, ext := range workerconfig.KV {
			if len(ext.ResourceID) != 0 {
				var kvresources = models.KV{}
				db.Model(&models.KV{}).Where(&models.KV{UID: ext.ResourceID, UserID: uint64(UserID)}).First(&kvresources)
				// 配置redis
				logrus.Printf("kvresources.ID: %v", kvresources.ID)
				if kvresources.ID != 0 {
				} else {
					ext.ResourceID = ""
				}
				workerconfig.KV[i] = ext
			}
		}

		for i, ext := range workerconfig.OSS {
			if len(ext.ResourceID) != 0 {
				var ossresources = models.OSS{}
				db.Model(&models.OSS{}).Where(&models.OSS{UID: ext.ResourceID, UserID: uint64(UserID)}).First(&ossresources)
				// 配置oss
				if ossresources.ID != 0 {
					ext.Bucket = ossresources.Bucket
					ext.Region = ossresources.Region
					ext.AccessKeyId = ossresources.AccessKey
					ext.AccessKeySecret = ossresources.SecretKey
				} else {
					ext.ResourceID = ""
				}
				workerconfig.OSS[i] = ext
			}
		}

		workerBytes, werr := json.Marshal(workerconfig)
		if werr != nil {
			logrus.Errorf("Failed to marshal worker config: %v", werr)
			return worker.Template
		} else {
			return string(workerBytes)
		}

	}
	return worker.Template
}

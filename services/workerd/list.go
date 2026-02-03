package workerd

import (
	"encoding/json"
	"strconv"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/defs"
	"vvorker/entities"
	"vvorker/funcs"
	"vvorker/models"
	"vvorker/models/secrets"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetWorkersEndpoint(c *gin.Context) {

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

type SimpleWorker struct {
	UID           string `json:"UID"`
	Name          string `json:"Name"`
	NodeName      string `json:"NodeName"`
	AccessControl bool   `json:"AccessControl"`
	Description   string `json:"Description"`
	IsCollab      bool   `json:"IsCollab"`
}

func GetAllWorkersEndpoint(c *gin.Context) {

	userID := c.GetUint(common.UIDKey)

	// 获取用户拥有的 Workers
	ownedWorkers, err := models.GetWorkersByUserID(uint64(userID))
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	// 获取用户参与协作的 Workers
	collabWorkerUIDs, err := models.GetUserCollaboratedWorkers(uint64(userID))
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	// 创建一个 map 来快速查找拥有的 worker UIDs
	ownedWorkerUIDs := make(map[string]bool)
	for _, worker := range ownedWorkers {
		ownedWorkerUIDs[worker.UID] = true
	}

	// 获取所有协作 workers 的详细信息
	var collabWorkers []*models.Worker
	for _, uid := range collabWorkerUIDs {
		// 跳过已经是拥有的 workers
		if ownedWorkerUIDs[uid] {
			continue
		}
		var worker models.Worker
		if err := database.GetDB().Where(&models.Worker{Worker: &entities.Worker{UID: uid}}).First(&worker).Error; err != nil {
			continue
		}
		collabWorkers = append(collabWorkers, &worker)
	}

	// 合并拥有的和协作的 workers
	var simpleWorkers []*SimpleWorker
	for _, worker := range ownedWorkers {
		simpleWorkers = append(simpleWorkers, &SimpleWorker{
			UID:           worker.UID,
			Name:          worker.Name,
			NodeName:      worker.NodeName,
			AccessControl: worker.EnableAccessControl,
			Description:   worker.Description,
			IsCollab:      false,
		})
	}
	for _, worker := range collabWorkers {
		simpleWorkers = append(simpleWorkers, &SimpleWorker{
			UID:           worker.UID,
			Name:          worker.Name,
			NodeName:      worker.NodeName,
			AccessControl: worker.EnableAccessControl,
			Description:   worker.Description,
			IsCollab:      true,
		})
	}

	common.RespOK(c, "get all workers success", simpleWorkers)
}

type GetWorkerRespose struct {
	UID         string `json:"UID"`
	NodeName    string `json:"NodeName"`
	Name        string `json:"Name"`
	Version     string `json:"Version"`
	MaxCount    int32  `json:"MaxCount"`
	Description string `json:"Description"`
}

func GetWorkerEndpoint(c *gin.Context) {

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
	worker.Worker.Code = nil
	common.RespOK(c, "get workers success", []GetWorkerRespose{
		{
			UID:         worker.UID,
			NodeName:    worker.NodeName,
			Name:        worker.Name,
			Version:     worker.Version,
			MaxCount:    worker.MaxCount,
			Description: worker.Description,
		},
	})
}

type GetWorkerEndpointJSONReq struct {
	UID string `json:"uid"`
}

func GetWorkerEndpointJSON(c *gin.Context) {

	userID := c.GetUint(common.UIDKey)
	req := &GetWorkerEndpointJSONReq{}
	if err := json.NewDecoder(c.Request.Body).Decode(req); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}
	uid := req.UID

	if len(uid) == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is empty", nil)
		return
	}
	worker, err := models.GetWorkerByUID(userID, uid)
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	worker.Worker.Code = nil
	worker.Worker.Template = ""
	common.RespOK(c, "get workers success", []*models.Worker{worker})
}

func GetWorkerEndpointAgent(c *gin.Context) {

	req := &entities.AgentGetWorkerByUIDReq{}
	if err := json.NewDecoder(c.Request.Body).Decode(req); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}
	uid := req.UID

	db := database.GetDB()
	worker := &models.Worker{}
	if err := db.Model(&models.Worker{}).Where(&models.Worker{
		Worker: &entities.Worker{
			UID: uid,
		},
	}).First(&worker).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	common.RespOK(c, "get workers success", models.Trans2Entities([]*models.Worker{worker}))
}

func AgentSyncWorkers(c *gin.Context) {

	req := &entities.AgentSyncWorkersReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.RespErr(c, defs.CodeInvalidRequest, err.Error(), nil)
		return
	}

	nodeName := c.GetString(defs.KeyNodeName)
	// get node's workerlist
	// workers, err := models.AdminGetWorkersByNodeName(nodeName)
	// if err != nil {
	// 	common.RespErr(c, defs.CodeInternalError, err.Error(), nil)
	// 	return
	// }
	logrus.Infof("sync workers, node name: %s", nodeName)

	var workers []*models.Worker
	db := database.GetDB()

	if err := db.Model(&models.Worker{}).Where(&models.Worker{
		Worker: &entities.Worker{
			NodeName: nodeName,
		},
	}).Select("uid", "version").Find(&workers).Error; err != nil {
		common.RespErr(c, defs.CodeInternalError, err.Error(), nil)
		return
	}
	var workerUIDVersions []entities.WorkerUIDVersion
	for _, worker := range workers {
		workerUIDVersions = append(workerUIDVersions, entities.WorkerUIDVersion{
			UID:     worker.UID,
			Version: worker.Version,
		})
	}

	resp := &entities.AgentDiffSyncWorkersResp{
		WorkerUIDVersions: workerUIDVersions,
	}
	common.RespOK(c, "sync workers success", resp)
}

func FillWorkerConfig(c *gin.Context) {

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

				funcs.MigratePostgreSQLDatabase(worker.UserID, ext.ResourceID)
			} else {
				if len(ext.Migrate) != 0 {
					funcs.MigratePostgreSQLDatabase(worker.UserID, "worker_resource:pgsql:"+worker.UID+":"+ext.Migrate)
				}
			}
		}

		for i, ext := range workerconfig.Mysql {
			if len(ext.ResourceID) != 0 {
				var mysqlresources = models.MySQL{}
				db.Model(&models.MySQL{}).Where(&models.MySQL{UID: ext.ResourceID, UserID: uint64(UserID)}).First(&mysqlresources)
				if mysqlresources.ID != 0 {
					if conf.AppConfigInstance.ServerMySQLOneDBName != "" {
						ext.Database = conf.AppConfigInstance.ServerMySQLOneDBName
						ext.User = conf.AppConfigInstance.ServerMySQLUser
						ext.Password = conf.AppConfigInstance.ServerMySQLPassword
					} else {
						ext.Database = mysqlresources.Database
						ext.Password = mysqlresources.Password
						ext.User = mysqlresources.Username
					}
				} else {
					ext.ResourceID = ""
				}
				workerconfig.Mysql[i] = ext

				funcs.MigrateMySQLDatabase(worker.UserID, ext.ResourceID)
			} else {
				if len(ext.Migrate) != 0 {
					funcs.MigrateMySQLDatabase(worker.UserID, "worker_resource:mysql:"+worker.UID+":"+ext.Migrate)
				}
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
				ext.Provider = conf.AppConfigInstance.KVProvider
				workerconfig.KV[i] = ext
			}
		}

		for i, ext := range workerconfig.OSS {
			if len(ext.ResourceID) != 0 {
				var ossresources = models.OSS{}
				db.Model(&models.OSS{}).Where(&models.OSS{UID: ext.ResourceID, UserID: uint64(UserID)}).First(&ossresources)
				// 配置oss
				if ossresources.ID != 0 {
					if !conf.AppConfigInstance.MinioSingleBucketMode {
						ext.Bucket = ossresources.Bucket
						ext.Region = ossresources.Region
						ext.AccessKeyId = ossresources.AccessKey
						ext.AccessKeySecret = ossresources.SecretKey
					} else {
						ext.Bucket = ossresources.Bucket
						ext.Region = conf.AppConfigInstance.ServerMinioRegion
						ext.AccessKeyId = conf.AppConfigInstance.ServerMinioAccess
						ext.AccessKeySecret = conf.AppConfigInstance.ServerMinioSecret
					}
				} else {
					ext.ResourceID = ""
				}
				workerconfig.OSS[i] = ext
			}
		}
		var workerSecrets []secrets.Secret
		db.Model(secrets.Secret{}).Where(&secrets.Secret{
			WorkerUID: worker.UID,
		}).Find(&workerSecrets)

		vars := gin.H{}
		err := json.Unmarshal(workerconfig.Vars, &vars)
		if err == nil {
			for i, s := range workerSecrets {
				vars[s.Key] = workerSecrets[i].Value
			}
		}

		varsBytes, _ := json.Marshal(vars)
		workerconfig.Vars = varsBytes

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

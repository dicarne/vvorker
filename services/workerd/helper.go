package workerd

import (
	"encoding/json"
	"vorker/conf"
	"vorker/defs"
	"vorker/entities"
	"vorker/models"
	"vorker/rpc"
	"vorker/utils"
	"vorker/utils/database"

	"github.com/lucasepe/codename"
	"github.com/sirupsen/logrus"
)

func FillWorkerValue(worker *entities.Worker, keepUID bool, UID string, UserID uint) {
	if !keepUID {
		worker.UID = utils.GenerateUID()
	}
	worker.UserID = uint64(UserID)
	worker.HostName = defs.DefaultHostName

	if len(worker.NodeName) == 0 {
		assignNode, err := models.GetAssignNode()
		if err == nil {
			worker.NodeName = assignNode.GetName()
		} else {
			worker.NodeName = defs.DefaultNodeName
		}
	}
	if node, err := models.GetNodeByNodeName(worker.NodeName); err == nil {
		worker.TunnelID = node.UID
	} else {
		worker.TunnelID = conf.AppConfigInstance.NodeID
	}

	worker.ExternalPath = defs.DefaultExternalPath

	if len(worker.Code) == 0 {
		worker.Code = []byte(defs.DefaultCode)
	}
	if len(worker.Entry) == 0 {
		worker.Entry = defs.DefaultEntry
	}
	if len(worker.Template) == 0 {
		workerconfig, err := conf.ParseWorkerConfig(worker.Template)
		if err != nil {
			db := database.GetDB()
			for i, ext := range workerconfig.PgSql {
				if len(ext.ResourceID) != 0 {
					var pgresources = models.PostgreSQL{}
					db.Model(&models.PostgreSQL{}).Where(&models.PostgreSQL{UID: ext.ResourceID}).First(&pgresources)
					ext.Database = pgresources.Database
					ext.Password = conf.AppConfigInstance.ServerPostgresPassword
					ext.User = conf.AppConfigInstance.ServerPostgresUser

					workerconfig.PgSql[i] = ext
				}
			}

			for i, ext := range workerconfig.KV {
				if len(ext.ResourceID) != 0 {
					var kvresources = models.KV{}
					db.Model(&models.KV{}).Where(&models.KV{UID: ext.ResourceID}).First(&kvresources)
					// 配置redis
					workerconfig.KV[i] = ext
				}
			}

			for i, ext := range workerconfig.OSS {
				if len(ext.ResourceID) != 0 {
					var ossresources = models.OSS{}
					db.Model(&models.OSS{}).Where(&models.OSS{UID: ext.ResourceID}).First(&ossresources)
					// 配置oss
					ext.Bucket = ossresources.Bucket
					ext.Region = ossresources.Region
					ext.AccessKeyId = ossresources.AccessKey
					ext.AccessKeySecret = ossresources.SecretKey
					workerconfig.OSS[i] = ext
				}
			}

			workerBytes, werr := json.Marshal(workerconfig)
			if werr != nil {
				logrus.Errorf("Failed to marshal worker config: %v", werr)
			} else {
				worker.Template = string(workerBytes)
			}
		}
	}

	// if the worker name is not unique, use the uid as the name
	if wl, err :=
		models.AdminGetWorkersByNames([]string{worker.Name}); len(wl) > 0 ||
		err != nil ||
		len(worker.Name) == 0 {
		if len(wl) == 1 {
			if UID == wl[0].UID {
				return
			}
		}
		rng, _ := codename.DefaultRNG()
		worker.Name = codename.Generate(rng, 0)
	}
}

func SyncAgent(w *entities.Worker) {
	go func(worker *entities.Worker) {
		if worker.NodeName == conf.AppConfigInstance.NodeName {
			return
		}

		targetNode, err := models.GetNodeByNodeName(worker.NodeName)
		if err != nil {
			logrus.Errorf("worker node is invalid, db error: %v", err)
			return
		}
		if err := rpc.EventNotify(targetNode.Node, defs.EventSyncWorkers, nil); err != nil {
			logrus.Errorf("emit event: %v error, err: %v", defs.EventSyncWorkers, err)
			return
		}
	}(w)
}

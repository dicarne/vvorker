package workerd

import (
	"vvorker/conf"
	"vvorker/defs"
	"vvorker/entities"
	"vvorker/models"
	"vvorker/rpc"
	"vvorker/utils"

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

	// if the worker name is not unique, use the uid as the name
	if wl, err := models.AdminGetWorkersByNames([]string{worker.Name}); len(wl) > 0 || err != nil || len(worker.Name) == 0 {
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

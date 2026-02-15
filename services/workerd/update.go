package workerd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/entities"
	"vvorker/exec"
	"vvorker/ext/kv/src/sys_cache"
	"vvorker/funcs"
	"vvorker/models"
	"vvorker/utils"
	"vvorker/utils/generate"
	permissions "vvorker/utils/permissions"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UpdateWorkerReq struct {
	*entities.Worker
	Description string `json:"Description"`
}

func UpdateWorker(userID uint, UID string, worker *entities.Worker, desc string) (string, error) {
	FillWorkerValue(worker, true, UID, userID)

	workerRecord, err := models.GetWorkerByUID(userID, UID)
	if err != nil {
		return "", err
	}

	// 创建部署任务
	traceID := utils.GenerateUID()
	if err := models.CreateTask(traceID, UID, "running", "deployment"); err != nil {
		logrus.WithError(err).Warn("failed to create deployment task")
	}

	curNodeName := conf.AppConfigInstance.NodeName

	if workerRecord.NodeName == curNodeName {
		exec.ExecManager.ExitCmd(workerRecord.GetUID())
	}

	// 删除旧的worker
	err = workerRecord.Delete()
	if err != nil {
		if traceID != "" {
			models.CompleteTask(traceID, "failed")
		}
		return traceID, err
	}

	// 创建新的worker
	newWorker := &models.Worker{Worker: worker,
		EnableAccessControl: workerRecord.EnableAccessControl,
		Description:         desc,
	}
	newWorker.Version = utils.GenerateUID()
	code := newWorker.Code
	if conf.AppConfigInstance.FileStorageUseOSS {
		newWorker.Code = nil
	}

	if conf.AppConfigInstance.FileStorageUseOSS {
		err = funcs.UploadFileToSysBucket(fmt.Sprintf("code/%s", newWorker.GetUID()), bytes.NewReader(code))
		if err != nil {
			if traceID != "" {
				models.CompleteTask(traceID, "failed")
			}
			return traceID, err
		}
	}

	err = newWorker.Create()
	if err != nil {
		if traceID != "" {
			models.CompleteTask(traceID, "failed")
		}
		return traceID, err
	}

	if worker.NodeName == curNodeName {
		err = generate.GenWorkerConfig(newWorker.ToEntity(), newWorker)
		if err != nil {
			if traceID != "" {
				models.CompleteTask(traceID, "failed")
			}
			return traceID, err
		}
		exec.ExecManager.RunCmd(worker.GetUID())
	}
	return traceID, nil
}

// 更新worker
func UpdateEndpointJSON(c *gin.Context) {

	var worker *UpdateWorkerReq
	if err := permissions.BindJSON(c, &worker); err != nil {
		return
	}

	UID := worker.UID
	userID, ok := common.RequireUID(c)
	if !ok {
		return
	}
	oldworker, err := permissions.CanWriteWorker(c, uint64(userID), UID)
	if err != nil {
		return
	}

	if worker.Worker.Code == nil {
		worker.Worker.Code = oldworker.Code
	}
	if worker.Worker.Template == "" {
		worker.Worker.Template = oldworker.Template
	}
	if worker.Worker.MaxCount == 0 {
		worker.Worker.MaxCount = oldworker.MaxCount
		if worker.Worker.MaxCount == 0 {
			worker.Worker.MaxCount = 1
		}
	}

	config := &conf.WorkerConfig{}
	err = json.Unmarshal([]byte(worker.Template), config)
	if err != nil {
		logrus.WithError(err).Errorf("update worker error, worker is: [%+v], error config json.", worker.Worker)
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	worker.Worker.SemVersion = config.Version

	traceID, err := UpdateWorker(uint(oldworker.UserID), UID, worker.Worker, worker.Description)
	if err != nil {
		logrus.WithError(err).Errorf("update worker error, worker is: [%+v]", worker.Worker)
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	c.Set("trace_id", traceID)

	if oldworker.Name != worker.Worker.Name {
		cacheKey := fmt.Sprintf("db:workerd:uid_name_%s_cache:", oldworker.Name)
		lockKey := fmt.Sprintf("db:workerd:uid_name_%s_lock:", oldworker.Name)
		sys_cache.Del(cacheKey)
		sys_cache.Del(lockKey)
	}

	traceIDValue, _ := c.Get("trace_id")
	traceID, _ = traceIDValue.(string)
	result := gin.H{
		"version": worker.Worker.Version,
	}
	if traceID != "" {
		result["task_id"] = traceID
	}
	common.RespOK(c, "update worker success", result)
}

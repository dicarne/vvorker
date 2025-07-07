package workerd

import (
	"encoding/json"
	"io"
	"runtime/debug"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/entities"
	"vvorker/exec"
	"vvorker/models"
	"vvorker/utils/generate"
	permissions "vvorker/utils/permissions"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UpdateWorkerReq struct {
	*entities.Worker
}

// 更新worker，请采用 UpdateWorkerWithFile
func UpdateEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()

	UID := c.Param("uid")
	var newWorker *UpdateWorkerReq
	if err := permissions.BindJSON(c, &newWorker); err != nil {
		return
	}

	userID := c.GetUint(common.UIDKey)
	oldworker, err := permissions.CanWriteWorker(c, uint64(userID), UID)
	if err != nil {
		return
	}

	if newWorker.Worker.Code == nil {
		// 如果没有更新代码，则从旧的代码中查找原本的代码进行覆盖
		newWorker.Worker.Code = oldworker.Code
	}
	if err := UpdateWorker(userID, UID, newWorker.Worker); err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	common.RespOK(c, "update worker success", nil)
}

// 【弃用】更新worker，请采用 UpdateWorkerWithFile
func UpdateEndpointJSON(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()

	var worker *UpdateWorkerReq
	if err := permissions.BindJSON(c, &worker); err != nil {
		return
	}

	UID := worker.UID
	userID := c.GetUint(common.UIDKey)
	oldworker, err := permissions.CanWriteWorker(c, uint64(userID), UID)
	if err != nil {
		return
	}

	if worker.Worker.Code == nil {
		worker.Worker.Code = oldworker.Code
	}

	if err := UpdateWorker(userID, UID, worker.Worker); err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	common.RespOK(c, "update worker success", nil)
}

func UpdateWorker(userID uint, UID string, worker *entities.Worker) error {
	FillWorkerValue(worker, true, UID, userID)

	workerRecord, err := models.GetWorkerByUID(userID, UID)
	if err != nil {
		return err
	}

	curNodeName := conf.AppConfigInstance.NodeName

	if workerRecord.NodeName == curNodeName {
		exec.ExecManager.ExitCmd(workerRecord.GetUID())
	}

	// 删除旧的worker
	err = workerRecord.Delete()
	if err != nil {
		return err
	}

	// 创建新的worker
	newWorker := &models.Worker{Worker: worker, EnableAccessControl: workerRecord.EnableAccessControl}
	err = newWorker.Create()
	if err != nil {
		return err
	}

	if worker.NodeName == curNodeName {
		err := generate.GenWorkerConfig(newWorker.ToEntity(), newWorker)
		if err != nil {
			return err
		}
		exec.ExecManager.RunCmd(worker.GetUID(), []string{})
	}
	return nil
}

func UpdateWorkerWithFile(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in UpdateWorkerWithFile: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()

	// 1. Parse the multipart form
	if err := c.Request.ParseMultipartForm(1 << 30); err != nil { // 32MB max memory
		common.RespErr(c, common.RespCodeInvalidParams, "failed to parse form", nil)
		return
	}

	// 2. Get the JSON data
	jsonData := c.PostForm("data")
	if jsonData == "" {
		common.RespErr(c, common.RespCodeAuthErr, "missing data field", nil)
		return
	}

	logrus.Info("jsonData: ", jsonData)

	var req UpdateWorkerReq
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		common.RespErr(c, common.RespCodeInvalidParams, "invalid json data", nil)
		return
	}
	userID := c.GetUint(common.UIDKey)

	// 3. Get the file
	oldWorker, err := permissions.CanWriteWorker(c, uint64(userID), req.UID)
	if err != nil {
		return
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		req.Worker.Code = oldWorker.Code
	} else {
		defer file.Close()
		// 4. Read file content
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			common.RespErr(c, common.RespCodeInternalError, "failed to read file", nil)
			return
		}

		// 5. Use the file content as the worker code
		req.Worker.Code = fileBytes
	}

	if err := UpdateWorker(userID, req.UID, req.Worker); err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	common.RespOK(c, "update worker with file success", nil)
}

package workerd

import (
	"vvorker/common" // 假设 common 包路径是这个
	"vvorker/models"

	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetWorkerInformationByIDEndpoint 根据 ID 获取 worker 信息
func GetWorkerInformationByIDEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()

	id := c.Param("id")
	workerInfo, err := models.GetWorkerInformationByUID(id)
	if err != nil {
		workerInfo = &models.WorkerInformation{
			WorkerInformationBase: &models.WorkerInformationBase{
				UID:         id,
				Description: "",
				Example:     "",
			},
		}
	}

	// 使用 common.RespOK 统一成功响应
	common.RespOK(c, "ok", workerInfo.WorkerInformationBase)
}

// UpdateWorkerInformationEndpoint 更新 worker 信息
func UpdateWorkerInformationEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()

	id := c.Param("id")
	var updatedInfo = &models.WorkerInformation{}
	if err := c.ShouldBindJSON(updatedInfo); err != nil {
		// 使用 common.RespErr 统一错误响应
		common.RespErr(c, 400, "参数解析失败", nil)
		return
	}

	userID := c.GetUint(common.UIDKey)
	// 如果worker中包含这个id
	if models.HasWorker(userID, id) {
		updatedInfo.UID = id
		if err := models.UpdateWorkerInformationByUID(id, updatedInfo); err != nil {
			logrus.WithError(err).Error("更新 worker 信息失败")
			// 使用 common.RespErr 统一错误响应
			common.RespErr(c, 500, "更新失败", nil)
			return
		}
		// 使用 common.RespOK 统一成功响应
		common.RespOK(c, "更新成功", nil)
	} else {
		// 使用 common.RespErr 统一错误响应
		common.RespErr(c, common.RespCodeInternalError, "worker information id err", nil)
	}
}

// DeleteWorkerInformationEndpoint 删除 worker 信息
func DeleteWorkerInformationEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()

	id := c.Param("id")
	if err := models.DeleteWorkerInformationByUID(id); err != nil {
		logrus.WithError(err).Error("删除 worker 信息失败")
		// 使用 common.RespErr 统一错误响应
		common.RespErr(c, 500, "删除失败", nil)
		return
	}

	// 使用 common.RespOK 统一成功响应
	common.RespOK(c, "删除成功", nil)
}

package workerd

import (
	"vvorker/common" // 假设 common 包路径是这个
	"vvorker/models"
	permissions "vvorker/utils/permissions"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetWorkerInformationByIDEndpoint 根据 ID 获取 worker 信息
func GetWorkerInformationByIDEndpoint(c *gin.Context) {

	id := c.Param("id")
	userID, ok := common.RequireUID(c)
	if !ok {
		return
	}
	// 检查用户是否有权限访问 Worker（拥有者或协作者）
	_, err := permissions.CanReadWorker(c, uint64(userID), id)
	if err != nil {
		// CanReadWorker 内部已经调用了 RespErr
		return
	}

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

	id := c.Param("id")
	var updatedInfo = &models.WorkerInformation{}
	if err := c.BindJSON(updatedInfo); err != nil {
		// 使用 common.RespErr 统一错误响应
		common.RespErr(c, 400, "参数解析失败", nil)
		return
	}

	userID, ok := common.RequireUID(c)
	if !ok {
		return
	}
	// 检查用户是否有写权限（拥有者或协作者）
	_, err := permissions.CanWriteWorker(c, uint64(userID), id)
	if err != nil {
		// CanWriteWorker 内部已经调用了 RespErr
		return
	}

	updatedInfo.UID = id
	if err := models.UpdateWorkerInformationByUID(id, updatedInfo); err != nil {
		logrus.WithError(err).Error("更新 worker 信息失败")
		// 使用 common.RespErr 统一错误响应
		common.RespErr(c, 500, "更新失败", nil)
		return
	}
	// 使用 common.RespOK 统一成功响应
	common.RespOK(c, "更新成功", nil)
}

// DeleteWorkerInformationEndpoint 删除 worker 信息
func DeleteWorkerInformationEndpoint(c *gin.Context) {

	id := c.Param("id")
	userID, ok := common.RequireUID(c)
	if !ok {
		return
	}
	// 只有拥有者可以删除 worker 信息
	_, err := permissions.CanManageWorkerMembers(c, uint64(userID), id)
	if err != nil {
		// CanManageWorkerMembers 内部已经调用了 RespErr
		return
	}

	if err := models.DeleteWorkerInformationByUID(id); err != nil {
		logrus.WithError(err).Error("删除 worker 信息失败")
		// 使用 common.RespErr 统一错误响应
		common.RespErr(c, 500, "删除失败", nil)
		return
	}

	// 使用 common.RespOK 统一成功响应
	common.RespOK(c, "删除成功", nil)
}

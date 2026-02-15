package workerd

import (
	"vvorker/common"
	"vvorker/exec"
	"vvorker/utils/database"
	permissions "vvorker/utils/permissions"

	"github.com/gin-gonic/gin"
)

// GetWorkerLogsReq 定义获取工作者日志请求结构体
type GetWorkerLogsReq struct {
	Page     int `json:"page"`      // 页码，从 1 开始
	PageSize int `json:"page_size"` // 每页记录数
}

// 定义返回结构体
type GetWorkerLogsResp struct {
	Total int                   `json:"total"` // 日志总数
	Logs  []*exec.WorkerLogData `json:"logs"`  // 日志列表
}

func GetWorkerLogsEndpoint(c *gin.Context) {
	var req *GetWorkerLogsReq
	if err := c.BindJSON(&req); err != nil {
		return
	}

	// 处理默认值，防止 page 和 page_size 为 0
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 计算 offset
	offset := (req.Page - 1) * req.PageSize

	UID := c.Param("uid")
	if len(UID) == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "uid is empty", nil)
		return
	}

	userID := c.GetUint(common.UIDKey)
	// 检查用户是否有权限访问 Worker（拥有者或协作者）
	_, err := permissions.CanReadWorker(c, uint64(userID), UID)
	if err != nil {
		// CanReadWorker 内部已经调用了 RespErr
		return
	}

	db := database.GetDB()
	var logs []*exec.WorkerLog
	var total int64
	// 先查询日志总数
	if err := db.Model(&exec.WorkerLog{}).Where("uid = ?", UID).Limit(10000).Count(&total).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	// 使用计算出的 offset 和 page_size 进行查询
	if err := db.Where("uid = ?", UID).Offset(offset).Limit(req.PageSize).Order("time desc").Find(&logs).Error; err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	var logs2 []*exec.WorkerLogData
	for _, log := range logs {
		logs2 = append(logs2, log.WorkerLogData)
	}

	// 封装返回数据
	resp := GetWorkerLogsResp{
		Total: int(total),
		Logs:  logs2,
	}
	common.RespOK(c, "get worker logs success", resp)
}

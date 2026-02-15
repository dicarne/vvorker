package task

import (
	"time"
	"vvorker/common"
	"vvorker/models"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
)

func init() {

}

type CreateTaskReq struct {
	WorkerUID string `json:"worker_uid"`
	TraceID   string `json:"trace_id"`
}

func CreateTaskEndpoint(c *gin.Context) {

	var req CreateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, 400, "error", gin.H{"error": "Invalid request"})
		return
	}

	db := database.GetDB()
	// 检查traceid是否存在
	var count int64
	if err := db.Model(&models.Task{}).Where(&models.Task{
		TraceID: req.TraceID,
	}).Limit(1).Count(&count).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		return
	}

	// 插入
	if err := db.Create(&models.Task{
		TraceID:   req.TraceID,
		WorkerUID: req.WorkerUID,
		Status:    "running",
		StartTime: time.Now(),
		Type:      "usertask",
	}).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": err.Error()})
		return
	}
	common.RespOK(c, "success", gin.H{"task_uid": req.TraceID})
}

type GetTaskStatusReq struct {
	WorkerUID string `json:"worker_uid"`
	TraceID   string `json:"trace_id"`
}

func CheckInterruptTaskEndpoint(c *gin.Context) {

	var req GetTaskStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, 400, "error", gin.H{"error": "Invalid request"})
		return
	}
	db := database.GetDB()
	var tt models.Task
	if err := db.Where(&models.Task{
		TraceID: req.TraceID,
	}).First(&tt).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		return
	}

	common.RespOK(c, "success", gin.H{"status": tt.Status, "result": tt.Result})
}

type LogTaskReq struct {
	WorkerUID string `json:"worker_uid"`
	TraceID   string `json:"trace_id"`
	Log       string `json:"log"`
}

func LogTaskEndpoint(c *gin.Context) {

	var req LogTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, 400, "error", gin.H{"error": "Invalid request"})
		return
	}

	db := database.GetDB()
	// 检查traceid是否存在
	var count int64
	if err := db.Model(&models.TaskLog{}).Where(&models.TaskLog{
		TraceID: req.TraceID,
	}).Limit(1).Count(&count).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		return
	}

	// 插入
	if err := db.Create(&models.TaskLog{
		TraceID: req.TraceID,
		Content: req.Log,
		Time:    time.Now(),
	}).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": err.Error()})
		return
	}
	common.RespOK(c, "success", gin.H{"task_uid": req.TraceID})

}

func CancelTaskEndpoint(c *gin.Context) {

	var req GetTaskStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, 400, "error", gin.H{"error": "Invalid request"})
		return
	}

	userID := c.GetUint(common.UIDKey)

	// 检查用户是否有权限访问该 worker
	if _, err := models.GetWorkerByUID(userID, req.WorkerUID); err != nil {
		common.RespErr(c, 403, "error", gin.H{"error": "No permission"})
		return
	}

	db := database.GetDB()

	var tt models.Task
	if err := db.Where(&models.Task{
		TraceID:   req.TraceID,
		WorkerUID: req.WorkerUID,
	}).First(&tt).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		return
	}
	tt.Status = "canceled"
	if err := db.Save(&tt).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		return
	}
	common.RespOK(c, "success", gin.H{"task_uid": req.TraceID})

}

func CompleteTaskEndpoint(c *gin.Context) {

	var req GetTaskStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, 400, "error", gin.H{"error": "Invalid request"})
		return
	}

	db := database.GetDB()

	var tt models.Task
	if err := db.Where(&models.Task{
		TraceID: req.TraceID,
	}).First(&tt).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		return
	}
	tt.Status = "completed"
	if err := db.Save(&tt).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		return
	}
	common.RespOK(c, "success", gin.H{"task_uid": req.TraceID})

}

type ListTaskReq struct {
	WorkerUID string `json:"worker_uid" binding:"required"`
	TraceID   string `json:"trace_id"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
}

type ListTaskResponse struct {
	models.Task
	WorkerName string `json:"worker_name"`
}

func ListTaskEndpoint(c *gin.Context) {

	var req ListTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, 400, "error", gin.H{"error": "Invalid request, worker_uid is required"})
		return
	}

	userID := c.GetUint(common.UIDKey)

	// 检查用户是否有权限访问该 worker
	worker, err := models.GetWorkerByUID(userID, req.WorkerUID)
	if err != nil {
		common.RespErr(c, 403, "error", gin.H{"error": "No permission"})
		return
	}

	db := database.GetDB()

	var total int64

	// 直接根据 worker_uid 查询任务
	if err := db.Model(&models.Task{}).
		Where(&models.Task{
			WorkerUID: req.WorkerUID,
			Type:      "usertask",
		}).
		Count(&total).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		return
	}

	var tasks []ListTaskResponse
	if err := db.Table("tasks").
		Select("tasks.*").
		Where(&models.Task{
			WorkerUID: req.WorkerUID,
			Type:      "usertask",
		}).
		Order("tasks.start_time desc").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&tasks).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		return
	}

	// 添加 worker_name
	for i := range tasks {
		tasks[i].WorkerName = worker.Name
	}

	common.RespOK(c, "success", gin.H{"tasks": tasks, "total": total})
}

type GetTaskLogsReq struct {
	WorkerUID string `json:"worker_uid" binding:"required"`
	TraceID   string `json:"trace_id"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
}

func GetLogsEndpoint(c *gin.Context) {

	var req GetTaskLogsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, 400, "error", gin.H{"error": "Invalid request, worker_uid is required"})
		return
	}

	userID := c.GetUint(common.UIDKey)

	// 检查用户是否有权限访问该 worker
	if _, err := models.GetWorkerByUID(userID, req.WorkerUID); err != nil {
		common.RespErr(c, 403, "error", gin.H{"error": "No permission"})
		return
	}

	db := database.GetDB()
	// 检查该任务是否属于该 worker
	var task models.Task
	if err := db.Where(&models.Task{
		TraceID:   req.TraceID,
		WorkerUID: req.WorkerUID,
	}).First(&task).Error; err != nil {
		common.RespErr(c, 400, "error", gin.H{"error": "Task not found"})
		return
	}

	var logs []models.TaskLog
	if err := db.Where(&models.TaskLog{TraceID: req.TraceID}).Order("time desc").Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&logs).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		return
	}

	var total int64
	if err := db.Model(&models.TaskLog{}).Where(&models.TaskLog{TraceID: req.TraceID}).Count(&total).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		return
	}
	common.RespOK(c, "success", gin.H{"logs": logs, "total": total})
}

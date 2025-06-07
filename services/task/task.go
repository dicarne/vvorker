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
	defer func() {
		if r := recover(); r != nil {
			common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		}
	}()

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
	defer func() {
		if r := recover(); r != nil {
			common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		}
	}()

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

	common.RespOK(c, "success", gin.H{"status": tt.Status})
}

type LogTaskReq struct {
	WorkerUID string `json:"worker_uid"`
	TraceID   string `json:"trace_id"`
	Log       string `json:"log"`
}

func LogTaskEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		}
	}()
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
	defer func() {
		if r := recover(); r != nil {
			common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		}
	}()
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
	tt.Status = "canceled"
	if err := db.Save(&tt).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		return
	}
	common.RespOK(c, "success", gin.H{"task_uid": req.TraceID})

}

func CompleteTaskEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		}
	}()
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
	WorkerUID string `json:"worker_uid"`
	TraceID   string `json:"trace_id"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
}

type ListTaskResponse struct {
	models.Task
	WorkerName string `json:"worker_name"`
}

func ListTaskEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		}
	}()
	var req ListTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, 400, "error", gin.H{"error": "Invalid request"})
		return
	}
	db := database.GetDB()

	userID := uint64(c.GetUint(common.UIDKey))

	var total int64

	// 查询符合条件的任务总数
	if err := db.Model(&models.Task{}).
		Joins("JOIN workers ON tasks.worker_uid = workers.uid").
		Where("workers.user_id = ?", userID).
		// Where("tasks.worker_uid = ?", req.WorkerUID).
		Count(&total).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		return
	}

	var tasks []ListTaskResponse
	// 通过 JOIN 关联 tasks 表和 workers 表，筛选出符合条件的任务，并获取 worker 的 name
	if err := db.Table("tasks").
		Select("tasks.*, workers.name as worker_name").
		Joins("JOIN workers ON tasks.worker_uid = workers.uid").
		Where("workers.user_id = ?", userID).
		Order("tasks.start_time desc").
		// Where("tasks.worker_uid = ?", req.WorkerUID).
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&tasks).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		return
	}

	common.RespOK(c, "success", gin.H{"tasks": tasks, "total": total})
}

type GetTaskLogsReq struct {
	WorkerUID string `json:"worker_uid"`
	TraceID   string `json:"trace_id"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
}

func GetLogsEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		}
	}()
	var req GetTaskLogsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, 400, "error", gin.H{"error": "Invalid request"})
		return
	}
	userID := uint64(c.GetUint(common.UIDKey))
	if userID == 0 {
		common.RespErr(c, 400, "error", gin.H{"error": "Invalid request"})
		return
	}
	db := database.GetDB()
	// 查找这个trace id 是否是这个worker的，这个worker是否是这个用户的
	var count int64
	if err := db.Model(&models.Task{}).
		Joins("JOIN workers ON tasks.worker_uid = workers.uid").
		Where("workers.user_id =?", userID).
		Where("tasks.trace_id =?", req.TraceID).
		Limit(1).
		Count(&count).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
	}
	if count == 0 {
		common.RespErr(c, 400, "error", gin.H{"error": "Invalid request"})
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

package models

import (
	"time"
	"vvorker/utils/database"

	"gorm.io/gorm"
)

type Task struct {
	*gorm.Model
	TraceID   string    `gorm:"index" json:"trace_id"`
	WorkerUID string    `gorm:"index" json:"worker_uid"`
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
	Result    string    `json:"result"` // 迁移结果: "success", "error: ..."
	Type      string    `json:"type"`   // 任务类型: "deployment", "usertask"
}

type TaskLog struct {
	*gorm.Model
	TraceID string    `gorm:"index" json:"trace_id"`
	Time    time.Time `json:"time"`
	Content string    `json:"content"`
	Type    string    `json:"type"`
}

// CreateTask 创建任务
func CreateTask(traceID, workerUID, status string, taskType string) error {
	db := database.GetDB()
	return db.Create(&Task{
		TraceID:   traceID,
		WorkerUID: workerUID,
		Status:    status,
		Type:      taskType,
		StartTime: time.Now(),
	}).Error
}

// CompleteTask 完成任务
func CompleteTask(traceID, status string) error {
	db := database.GetDB()
	return db.Model(&Task{}).Where("trace_id = ?", traceID).Updates(map[string]interface{}{
		"status":   status,
		"end_time": time.Now(),
	}).Error
}

// UpdateTaskResult 更新任务结果
func UpdateTaskResult(traceID, result string) error {
	db := database.GetDB()
	return db.Model(&Task{}).Where("trace_id = ?", traceID).Update("result", result).Error
}

// GetTask 获取任务
func GetTask(traceID string) (*Task, error) {
	db := database.GetDB()
	var task Task
	err := db.Where("trace_id = ?", traceID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// MarkRunningTasksAsInterrupt 将所有 running 状态的任务标记为 interrupt
func MarkRunningTasksAsInterrupt() error {
	db := database.GetDB()
	return db.Model(&Task{}).Where("status = ?", "running").Updates(map[string]interface{}{
		"status":   "interrupt",
		"end_time": time.Now(),
	}).Error
}

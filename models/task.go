package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	*gorm.Model
	TraceID   string    `gorm:"index" json:"trace_id"`
	WorkerUID string    `gorm:"index" json:"worker_uid"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
}

type TaskLog struct {
	*gorm.Model
	TraceID string    `gorm:"index" json:"trace_id"`
	Time    time.Time `json:"time"`
	Content string    `json:"content"`
	Type    string    `json:"type"`
}

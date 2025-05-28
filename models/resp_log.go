package models

import (
	"time"

	"gorm.io/gorm"
)

type ResponseLog struct {
	gorm.Model
	WorkerUID string `gorm:"index"`
	Method    string
	Path      string
	Status    int
	Time      time.Time `gorm:"index"`
}

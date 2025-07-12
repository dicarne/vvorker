package workercopy

import "gorm.io/gorm"

type WorkerCopy struct {
	gorm.Model
	WorkerUID   string
	LocalID     uint
	Port        uint
	ControlPort uint
}

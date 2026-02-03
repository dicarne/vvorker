package secrets

import "gorm.io/gorm"

type Secret struct {
	gorm.Model
	WorkerUID string
	Key       string
	Value     string
}

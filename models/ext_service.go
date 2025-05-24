package models

import "gorm.io/gorm"

type KV struct {
	gorm.Model
	UserID uint64
	UID    string `gorm:"unique"`
}

type OSS struct {
	gorm.Model
	UserID    uint64
	UID       string `gorm:"unique"`
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
	Name      string
}

type PostgreSQL struct {
	gorm.Model
	UserID   uint64
	UID      string `gorm:"unique"`
	Database string
}

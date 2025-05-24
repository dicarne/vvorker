package models

import "gorm.io/gorm"

type KV struct {
	gorm.Model
	UserName string
	UID      string `gorm:"unique"`
}

type OSS struct {
	gorm.Model
	UserName  string
	UID       string `gorm:"unique"`
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
}

type PostgreSQL struct {
	gorm.Model
	UserName string
	UID      string `gorm:"unique"`
	Database string
}

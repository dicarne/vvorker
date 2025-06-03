package models

import (
	"time"

	"gorm.io/gorm"
)

type KV struct {
	gorm.Model
	UserID   uint64
	UID      string `gorm:"unique"`
	Name     string
	Password string
}

type OSS struct {
	gorm.Model
	UserID     uint64
	UID        string `gorm:"unique"`
	AccessKey  string
	SecretKey  string
	Bucket     string
	Region     string
	Name       string
	Expiration time.Time
	SessionKey string
}

type PostgreSQL struct {
	gorm.Model
	UserID   uint64
	UID      string `gorm:"unique"`
	Database string
	Name     string
	Username string
	Password string
}

type Assets struct {
	gorm.Model
	UserID    uint64
	UID       string `gorm:"unique"`
	WorkerUID string
	Name      string
	MIME      string
	Hash      string
	Path      string `gorm:"index"`
	// Data      []byte
}

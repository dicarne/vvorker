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
	UserID       uint64
	UID          string `gorm:"unique"`
	AccessKey    string
	SecretKey    string
	Bucket       string
	Region       string
	Name         string
	Expiration   time.Time
	SessionKey   string
	SingleBucket bool
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

type PostgreSQLMigration struct {
	gorm.Model
	UserID           uint64
	DBUID            string
	FileName         string
	FileContent      string
	Sequence         int
	CustomDBName     string `json:"custom_db_name"`
	CustomDBUser     string `json:"custom_db_user"`
	CustomDBHost     string `json:"custom_db_host"`
	CustomDBPort     int    `json:"custom_db_port"`
	CustomDBPassword string `json:"custom_db_password"`
}

type MySQL struct {
	gorm.Model
	UserID   uint64
	UID      string `gorm:"unique"`
	Database string
	Name     string
	Username string
	Password string
}

type MySQLMigration struct {
	gorm.Model
	UserID           uint64
	DBUID            string
	FileName         string
	FileContent      string
	Sequence         int
	CustomDBName     string `json:"custom_db_name"`
	CustomDBUser     string `json:"custom_db_user"`
	CustomDBHost     string `json:"custom_db_host"`
	CustomDBPort     int    `json:"custom_db_port"`
	CustomDBPassword string `json:"custom_db_password"`
}

type Assets struct {
	gorm.Model
	UserID    uint64
	UID       string
	WorkerUID string
	Name      string
	MIME      string
	Hash      string
	Path      string `gorm:"index"`
	// Data      []byte
}

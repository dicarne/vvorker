package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
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
	MigrateKey       string `json:"migrate_key"`
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
	MigrateKey       string `json:"migrate_key"`
}

type MigrationHistory struct {
	gorm.Model
	Key   string `gorm:"index"`
	Error string
}

func GenerateMigrationHistoryKey(sqlType string, uid string, fileName string, content string) string {
	// md5(sqlType:uid:filename:md5(content))
	hash := md5.Sum([]byte(content))
	contentMd5 := hex.EncodeToString(hash[:])
	rawKey := fmt.Sprintf("%s:%s:%s:%s", sqlType, uid, fileName, contentMd5)
	hash2 := md5.Sum([]byte(rawKey))
	return hex.EncodeToString(hash2[:])
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

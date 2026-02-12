package database

import (
	"os"
	"path/filepath"
	"time"
	"vvorker/conf"
	"vvorker/defs"
	"vvorker/utils"

	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// max 返回两个整数中较大的一个
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func initSqlite() {
	if conf.AppConfigInstance.LitefsEnabled {
		utils.WaitForPort("localhost", conf.AppConfigInstance.LitefsPrimaryPort)
	}
	godotenv.Load()
	if conf.AppConfigInstance.DBType != defs.DBTypeSqlite {
		return
	}

	dbPath := conf.AppConfigInstance.DBPath

	// 确保数据库文件所在目录存在
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		logrus.WithError(err).Errorf("Failed to create directory for SQLite database: %s", filepath.Dir(dbPath))
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		logrus.Error(err, "Initializing DB Error")
		logrus.Panicf("DB PATH: %s", dbPath)
		panic(err)
	}

	// 配置连接池 - 防止连接超时和重置
	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to get database instance")
	}

	// SQLite 使用更小的连接池配置（因为是文件数据库）
	// 使用主配置的值除以2-3倍
	maxIdleConns := max(1, conf.AppConfigInstance.DBMaxIdleConns/2)
	maxOpenConns := max(3, conf.AppConfigInstance.DBMaxOpenConns/3)
	connMaxLifetime := time.Duration(conf.AppConfigInstance.DBConnMaxLifetime*2) * time.Minute
	connMaxIdleTime := time.Duration(conf.AppConfigInstance.DBConnMaxIdleTime*5) * time.Minute

	// 连接池配置
	sqlDB.SetMaxIdleConns(maxIdleConns)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(maxOpenConns)           // 最大打开连接数
	sqlDB.SetConnMaxLifetime(connMaxLifetime)     // 连接最大生命周期
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)     // 空闲连接超时时间

	logrus.Infof("SQLite database initialized with connection pool: max_idle=%d, max_open=%d, max_lifetime=%v, max_idle_time=%v",
		maxIdleConns, maxOpenConns, connMaxLifetime, connMaxIdleTime)

	GetManager().SetDB(defs.DBTypeSqlite, db)
}

// func GetSqlite() *gorm.DB {
// 	dbPath := conf.AppConfigInstance.DBPath
// 	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
// 	if err != nil {
// 		return nil
// 	}
// 	return db
// }

package database

import (
	"database/sql"
	"fmt"
	"time"
	"vvorker/conf"
	"vvorker/defs"
	"vvorker/utils"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initPgsql() {
	if conf.AppConfigInstance.LitefsEnabled {
		utils.WaitForPort("localhost", conf.AppConfigInstance.LitefsPrimaryPort)
	}
	godotenv.Load()
	if conf.AppConfigInstance.DBType != defs.DBTypePostgres {
		return
	}

	databaseName := "vvorker_admin_" + conf.AppConfigInstance.NodeName

	if conf.AppConfigInstance.DBName != "" {
		databaseName = conf.AppConfigInstance.DBName
	}

	pgdb, err := sql.Open("postgres",
		"user="+conf.AppConfigInstance.ServerPostgreUser+
			" password="+conf.AppConfigInstance.ServerPostgrePassword+
			" host="+conf.AppConfigInstance.ServerPostgreHost+
			" port="+fmt.Sprintf("%d", conf.AppConfigInstance.ServerPostgrePort)+
			" sslmode=disable")
	if err != nil {
		panic("Failed to connect to database")
	}
	defer pgdb.Close()

	_, err = pgdb.Exec("CREATE DATABASE " + databaseName)
	if err != nil {
		logrus.Printf("Failed to create database: %v", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		conf.AppConfigInstance.ServerPostgreHost,
		conf.AppConfigInstance.ServerPostgreUser,
		conf.AppConfigInstance.ServerPostgrePassword,
		databaseName,
		conf.AppConfigInstance.ServerPostgrePort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// 配置连接池 - 防止连接超时和重置
	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to get database instance")
	}

	// 从环境变量读取连接池配置
	maxIdleConns := conf.AppConfigInstance.DBMaxIdleConns
	maxOpenConns := conf.AppConfigInstance.DBMaxOpenConns
	connMaxLifetime := time.Duration(conf.AppConfigInstance.DBConnMaxLifetime) * time.Minute
	connMaxIdleTime := time.Duration(conf.AppConfigInstance.DBConnMaxIdleTime) * time.Minute

	// 连接池配置
	sqlDB.SetMaxIdleConns(maxIdleConns)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(maxOpenConns)           // 最大打开连接数
	sqlDB.SetConnMaxLifetime(connMaxLifetime)     // 连接最大生命周期
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)     // 空闲连接超时时间

	logrus.Infof("PostgreSQL database initialized with connection pool: max_idle=%d, max_open=%d, max_lifetime=%v, max_idle_time=%v",
		maxIdleConns, maxOpenConns, connMaxLifetime, connMaxIdleTime)

	GetManager().SetDB(defs.DBTypePostgres, db)
}

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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func buildMysqlConnectionString() string {
	// username:password@protocol(address)/dbname?param=value
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
		conf.AppConfigInstance.ServerMySQLUser,
		conf.AppConfigInstance.ServerMySQLPassword,
		conf.AppConfigInstance.ServerMySQLHost,
		conf.AppConfigInstance.ServerMySQLPort)
}

func buildMysqlDBConnectionString(database string) string {

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.AppConfigInstance.ServerMySQLUser,
		conf.AppConfigInstance.ServerMySQLPassword,
		conf.AppConfigInstance.ServerMySQLHost,
		conf.AppConfigInstance.ServerMySQLPort,
		database)
}

func initMysql() {
	if conf.AppConfigInstance.LitefsEnabled {
		utils.WaitForPort("localhost", conf.AppConfigInstance.LitefsPrimaryPort)
	}
	godotenv.Load()
	if conf.AppConfigInstance.DBType != defs.DBTypeMysql {
		return
	}

	databaseName := "vvorker_admin_" + conf.AppConfigInstance.NodeName

	if conf.AppConfigInstance.DBName != "" {
		databaseName = conf.AppConfigInstance.DBName
	}

	pgdb, err := sql.Open("mysql", buildMysqlConnectionString())
	if err != nil {
		panic("Failed to connect to database")
	}
	defer pgdb.Close()

	_, err = pgdb.Exec("CREATE DATABASE " + databaseName)
	if err != nil {
		logrus.Printf("Failed to create database: %v", err)
	}

	dsn := buildMysqlDBConnectionString(databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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

	logrus.Infof("MySQL database initialized with connection pool: max_idle=%d, max_open=%d, max_lifetime=%v, max_idle_time=%v",
		maxIdleConns, maxOpenConns, connMaxLifetime, connMaxIdleTime)

	GetManager().SetDB(defs.DBTypeMysql, db)
}

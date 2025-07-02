package database

import (
	"database/sql"
	"fmt"
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

	GetManager().SetDB(defs.DBTypeMysql, db)
}

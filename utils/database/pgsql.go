package database

import (
	"database/sql"
	"fmt"
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

	databaseNmae := "vvorker_admin_" + conf.AppConfigInstance.NodeName

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

	_, err = pgdb.Exec("CREATE DATABASE " + databaseNmae)
	if err != nil {
		logrus.Printf("Failed to create database: %v", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		conf.AppConfigInstance.ServerPostgreHost,
		conf.AppConfigInstance.ServerPostgreUser,
		conf.AppConfigInstance.ServerPostgrePassword,
		databaseNmae,
		conf.AppConfigInstance.ServerPostgrePort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	GetManager().SetDB(defs.DBTypePostgres, db)
}

package database

import (
	"vvorker/conf"
	"vvorker/defs"
	"vvorker/utils"

	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func initSqlite() {
	if conf.AppConfigInstance.LitefsEnabled {
		utils.WaitForPort("localhost", conf.AppConfigInstance.LitefsPrimaryPort)
	}
	godotenv.Load()
	if conf.AppConfigInstance.DBType != defs.DBTypeSqlite {
		return
	}

	dbPath := conf.AppConfigInstance.DBPath
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		logrus.Error(err, "Initializing DB Error")
		logrus.Panicf("DB PATH: %s", dbPath)
		panic(err)
	}

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

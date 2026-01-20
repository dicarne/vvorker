package models

import (
	"time"
	"vvorker/conf"
	"vvorker/exec"
	workercopy "vvorker/models/worker_copy"
	"vvorker/utils"
	"vvorker/utils/database"

	"github.com/sirupsen/logrus"
)

func MigrateNormalModel() {
	normalModels := []interface{}{
		&User{}, &Worker{}, &WorkerVersion{}, &File{}, &KV{}, &OSS{}, &PostgreSQL{}, &AccessKey{},
		&WorkerInformation{}, &exec.WorkerLog{}, &ResponseLog{}, &Assets{}, &Task{}, &TaskLog{},
		&InternalServerWhiteList{}, &ExternalServerAKSK{}, &ExternalServerToken{}, &AccessRule{},
		&PostgreSQLMigration{}, &MySQL{}, &MySQLMigration{}, &workercopy.WorkerCopy{}, &MigrationHistory{},
	}
	if conf.AppConfigInstance.LitefsEnabled {
		if !conf.IsMaster() {
			return
		}
		utils.WaitForPort("localhost", conf.AppConfigInstance.LitefsPrimaryPort)
	}
	db := database.GetDB()
	// db.Exec("DROP INDEX IF EXISTS uni_users_email ON users;")
	for err := db.AutoMigrate(normalModels...); err != nil; err = db.AutoMigrate(
		normalModels...) {
		logrus.WithError(err).Errorf("auto migrate models error, sleep 5s and retry")
		time.Sleep(5 * time.Second)
	}

	// 确保至少存在一个管理员
	if err := EnsureAdminExists(); err != nil {
		logrus.WithError(err).Errorf("failed to ensure admin exists")
	}
}

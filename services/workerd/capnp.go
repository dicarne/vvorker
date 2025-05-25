package workerd

import (
	"errors"
	"path/filepath"
	"vvorker/conf"
	"vvorker/defs"
	"vvorker/entities"
	"vvorker/models"
	"vvorker/utils"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func GenCapnpConfig() error {
	workerRecords, err := models.AdminGetWorkersByNodeName(conf.AppConfigInstance.NodeName)
	if err != nil {
		logrus.Errorf("failed to get all workers, err: %v", err)
	}

	workerList := models.Trans2Entities(workerRecords)

	var hasError bool
	for _, worker := range workerList {
		w := &models.Worker{Worker: worker}
		fileMap := utils.BuildCapfile([]*entities.Worker{w.ToEntity()})

		if fileContent, ok := fileMap[worker.GetUID()]; ok {
			// 打印生成的 capnp 文件路径
			capnpFilePath := filepath.Join(
				conf.AppConfigInstance.WorkerdDir,
				defs.WorkerInfoPath,
				worker.GetUID(),
				defs.CapFileName,
			)
			logrus.Infof("Generating capnp file at: %s", capnpFilePath)

			err := utils.WriteFile(capnpFilePath, fileContent)
			if err != nil {
				logrus.WithError(err).Errorf("failed to write file, worker is: %+v", worker.Name)
				hasError = true
			}
		}
	}

	if hasError {
		logrus.Infof("GenCapnpConfig has error: %v, workerList: %+v", hasError,
			lo.SliceToMap(workerList, func(w *entities.Worker) (string, bool) { return w.GetUID(), true }))
	}

	if hasError {
		return errors.New("GenCapnpConfig has error")
	}
	return nil
}

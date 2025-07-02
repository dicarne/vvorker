package models

import (
	"vvorker/conf"

	"github.com/sirupsen/logrus"
)

func NodeWorkersInit() {
	workerRecords, err := AdminGetWorkersByNodeName(conf.AppConfigInstance.NodeName)
	if err != nil {
		logrus.Errorf("init failed to get all workers, err: %v", err)
	}
	logrus.Infof("this node will init %d workers", len(workerRecords))
	for _, worker := range workerRecords {
		if err := worker.Flush(); err != nil {
			logrus.WithError(err).Errorf("init failed to flush worker, worker is: [%+v]", worker.ToEntity())
		}
	}
}

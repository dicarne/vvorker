package models

import (
	"vvorker/conf"
	"vvorker/exec"
	"vvorker/utils"

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
		if err := utils.GenWorkerConfig(worker.ToEntity()); err != nil {
			logrus.WithError(err).Errorf("init failed to gen worker config, worker is: [%+v]", worker.ToEntity())
		}
		exec.ExecManager.RunCmd(worker.GetUID(), []string{})
	}
}

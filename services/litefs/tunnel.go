package litefs

import (
	"vvorker/common"
	"vvorker/conf"
	"vvorker/tunnel"

	"github.com/sirupsen/logrus"
)

func InitTunnel() {
	if !conf.AppConfigInstance.LitefsEnabled {
		return
	}
	if conf.IsMaster() {
		err := tunnel.GetClient().AddService(common.ServiceLitefs, conf.AppConfigInstance.LitefsPrimaryPort)
		if err != nil {
			logrus.WithError(err).Errorf("init tunnel for master litefs service error")
			return
		}
		logrus.Infof("init tunnel for litefs serivce success")
		return
	} else {
		err := tunnel.GetClient().AddVisitor(common.ServiceLitefs, conf.AppConfigInstance.LitefsPrimaryPort)
		if err != nil {
			logrus.WithError(err).Errorf("init tunnel for agent litefs visitor failed")
			return
		}
		logrus.Info("init tunnel for agent litefs visitor success")
		return
	}
}

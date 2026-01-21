package agent

import (
	"vvorker/common"
	"vvorker/defs"
	"vvorker/entities"
	"vvorker/utils/request"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	EventRouterImplInstance.RegisteHandler(defs.EventSyncWorkers, SyncEventHandler)
	EventRouterImplInstance.RegisteHandler(defs.EventAddWorker, AddWorkerEventHandler)
	EventRouterImplInstance.RegisteHandler(defs.EventDeleteWorker, DelWorkerEventHandler)
	EventRouterImplInstance.RegisteHandler(defs.EventFlushWorker, FlushWorkerEventHandler)
}

func NotifyEndpoint(c *gin.Context) {

	req := &entities.NotifyEventRequest{}
	err := request.Bind[*entities.NotifyEventRequest](c, req)
	if err != nil {
		logrus.Errorf("event: %s error, err: %+v", req.EventName, err)
		common.RespErr(c, common.RespCodeInvalidRequest, common.RespMsgInvalidRequest, nil)
		return
	}

	EventRouterImplInstance.Handle(c, req)
}

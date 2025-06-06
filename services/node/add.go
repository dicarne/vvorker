package node

import (
	"runtime/debug"
	"vvorker/common"
	"vvorker/defs"
	"vvorker/entities"
	"vvorker/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func AddEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	nodeName := c.GetString(defs.KeyNodeName)
	if n, err := models.GetNodeByNodeName(nodeName); err == nil && n != nil {
		common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		return
	}

	newNode := &models.Node{
		Node: &entities.Node{
			UID:  uuid.New().String(),
			Name: nodeName,
		},
	}

	if err := newNode.Create(); err != nil {
		logrus.Errorf("failed to create node, err: %v", err)
		common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		return
	}

	common.RespOK(c, common.RespMsgOK, nil)
}

package node

import (
	"fmt"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/defs"
	"vvorker/models"
	"vvorker/utils"
	"vvorker/utils/request"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"
)

func GetNodeInfoEndpoint(c *gin.Context) {

	nodeName := c.GetString(defs.KeyNodeName)

	node, err := models.GetNodeByNodeName(nodeName)
	if err != nil {
		logrus.Errorf("failed to get node info, err: %v", err)
		common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		return
	}

	// go rpc.EventNotify(node.Node, defs.EventSyncWorkers, nil)
	common.RespOK(c, common.RespMsgOK, node)
}

func UserGetNodesEndpoint(c *gin.Context) {

	nodes, err := models.AdminGetAllNodes()
	if err != nil {
		logrus.Errorf("failed to get nodes, err: %v", err)
		common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		return
	}
	pingMap := map[string]int{}
	var wg conc.WaitGroup

	for _, node := range nodes {
		nodeName := node.Name
		nodeUID := node.UID
		wg.Go(func() {
			var addr string
			if nodeName == conf.AppConfigInstance.NodeName {
				addr = fmt.Sprintf("http://%s:%d/api/ping", conf.AppConfigInstance.TunnelHost, conf.AppConfigInstance.APIPort)
			} else {
				addr = fmt.Sprintf("http://%s:%d/api/ping", conf.AppConfigInstance.TunnelHost, conf.AppConfigInstance.TunnelEntryPort)
			}

			pingMap[nodeName], err = request.Ping(
				addr, utils.NodeHost(nodeName, nodeUID))
			if err != nil {
				logrus.Errorf("failed to ping node %s, err: %v", nodeName, err)
				pingMap[nodeName] = 9999
			}
		})
	}
	wg.Wait()

	common.RespOK(c, common.RespMsgOK, gin.H{
		"nodes": nodes,
		"ping":  pingMap,
	})
}

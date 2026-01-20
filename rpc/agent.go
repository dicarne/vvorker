package rpc

import (
	"errors"
	"fmt"
	"vvorker/conf"
	"vvorker/defs"
	"vvorker/entities"
	"vvorker/utils"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func EventNotify(n *entities.Node, eventName string, extra map[string][]byte) error {
	logrus.Infof("event notify, eventName: %s, requestExtraKeys: %+v", eventName, lo.Keys(extra))
	reqResp, err := RPCWrapper().
		SetHeader(defs.HeaderHost, utils.NodeHost(n.Name, n.UID)).
		SetBody(&entities.NotifyEventRequest{EventName: eventName, Extra: extra}).
		Post(
			fmt.Sprintf("http://%s:%d/api/agent/notify",
				conf.AppConfigInstance.TunnelHost,
				conf.AppConfigInstance.TunnelEntryPort))

	if err != nil || reqResp.StatusCode >= 299 {
		logrus.Errorf("event notify error, err: %+v, resp: %+v, eventName: %s, requestExtraKeys: %+v", err, reqResp, eventName, lo.Keys(extra))
		return errors.New("error")
	}
	return nil
}

func SyncAgent(endpoint string) ([]entities.WorkerUIDVersion, error) {
	url := endpoint + "/api/agent/sync"
	resp := &entities.AgentDiffSyncWorkersResp{}
	rtype := struct {
		Code int                               `json:"code"`
		Msg  string                            `json:"msg"`
		Data entities.AgentDiffSyncWorkersResp `json:"data"`
	}{}

	reqResp, err := RPCWrapper().
		SetBody(&entities.AgentSyncWorkersReq{}).
		SetSuccessResult(&rtype).
		Post(url)
	resp = &rtype.Data
	logrus.Infof("sync agent length: %d", len(resp.WorkerUIDVersions))

	if err != nil || reqResp.StatusCode >= 299 {
		return nil, errors.New("error")
	}
	return resp.WorkerUIDVersions, nil
}

func AddNode(endpoint string) error {
	url := endpoint + "/api/agent/add"
	rtype := struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{}

	reqResp, err := RPCWrapper().
		SetBody(&entities.AgentSyncWorkersReq{}).
		SetSuccessResult(&rtype).
		Post(url)

	if err != nil || reqResp.StatusCode >= 299 {
		return errors.New("error")
	}
	return nil
}

func GetNode(endpoint string) (*entities.Node, error) {
	url := endpoint + "/api/agent/nodeinfo"
	rtype := struct {
		Code int            `json:"code"`
		Msg  string         `json:"msg"`
		Data *entities.Node `json:"data"`
	}{}

	reqResp, err := RPCWrapper().
		SetSuccessResult(&rtype).
		Get(url)

	if err != nil || reqResp.StatusCode >= 299 {
		return nil, errors.New("error")
	}
	return rtype.Data, nil
}

func RPCWrapper() *req.Request {
	return req.C().R().
		SetHeaders(map[string]string{
			defs.HeaderNodeName:   conf.AppConfigInstance.NodeName,
			defs.HeaderNodeSecret: conf.RPCToken,
		})
}

func GetWorkerByUID(endpoint string, UID string) (*entities.Worker, error) {
	url := endpoint + "/api/agent/get-worker"
	rtype := struct {
		Code int                `json:"code"`
		Msg  string             `json:"msg"`
		Data []*entities.Worker `json:"data"`
	}{}

	reqResp, err := RPCWrapper().
		SetBody(&entities.AgentGetWorkerByUIDReq{
			UID: UID,
		}).
		SetSuccessResult(&rtype).
		Post(url)
	if err != nil {
		return nil, errors.New("error: " + err.Error())
	}
	if reqResp.StatusCode >= 299 {
		return nil, errors.New("error: " + rtype.Msg)
	}
	if len(rtype.Data) == 0 {
		return nil, errors.New("error: worker not found")
	}
	return rtype.Data[0], nil
}

// func FillWorkerConfig(endpoint string, UID string) (string, error) {
// 	url := endpoint + "/api/agent/fill-worker-config"
// 	rtype := entities.AgentFillWorkerResp{}

// 	reqResp, err := RPCWrapper().
// 		SetBody(&entities.AgentFillWorkerReq{
// 			UID: UID,
// 		}).
// 		SetSuccessResult(&rtype).
// 		Post(url)

// 	if err != nil || reqResp.StatusCode >= 299 {
// 		return "", errors.New("error")
// 	}
// 	return rtype.NewTemplate, nil
// }

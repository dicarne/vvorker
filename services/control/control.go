package control

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"net/http"
	"vvorker/conf"

	"github.com/sirupsen/logrus"
)

// 请求控制接口
func RequestControlEndpoint(workerUID string, bbody []byte) []byte {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d",
		conf.AppConfigInstance.TunnelHost, conf.AppConfigInstance.TunnelEntryPort), nil) // method: GET/POST等
	if err != nil {
		logrus.Println(err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Host = workerUID + "-control" + conf.AppConfigInstance.WorkerURLSuffix

	body := bytes.NewBuffer(bbody)
	req.Body = io.NopCloser(body)
	req.ContentLength = int64(body.Len())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Println(err)
		return nil
	}
	defer resp.Body.Close()

	rbody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Println(err)
		return nil
	}
	return rbody
}

type ScheduledController struct {
	ScheduledTime int64  `json:"scheduledTime"`
	Cron          string `json:"cron"`
	Type          string `json:"type"`
}

func SendSchedulerEvent(workerUID string, cron string) {
	e := ScheduledController{
		ScheduledTime: time.Now().Unix(),
		Cron:          cron,
		Type:          "scheduled",
	}
	bbody, err := json.Marshal(e)
	if err != nil {
		logrus.Println(err)
		return
	}

	RequestControlEndpoint(workerUID, bbody)
}

package control

import (
	"bytes"
	"fmt"
	"io"

	"net/http"
	"vvorker/conf"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ControlEndpoint(c *gin.Context) {
	c.Data(http.StatusOK, "text/plain", RequestControlEndpoint(c.Query("uid")))
}

// 请求控制接口
func RequestControlEndpoint(workerUID string) []byte {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d",
		conf.AppConfigInstance.TunnelHost, conf.AppConfigInstance.TunnelEntryPort), nil) // method: GET/POST等
	if err != nil {
		logrus.Panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Host = workerUID + "-control.vvorker.zhishudali.ink"

	body := bytes.NewBuffer([]byte(workerUID))
	req.Body = io.NopCloser(body)
	req.ContentLength = int64(body.Len())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Panic(err)
	}
	defer resp.Body.Close()

	rbody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Panic(err)
	}
	return rbody
}

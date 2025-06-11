package proxy

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"runtime/debug"
	"strings"
	"time"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/models"
	"vvorker/tunnel"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Endpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	host := c.Request.Host
	c.Request.Host = host
	workerName := host[:len(host)-len(conf.AppConfigInstance.WorkerURLSuffix)]
	worker, err := models.AdminGetWorkerByName(workerName)
	if err != nil {
		logrus.Errorf("failed to get worker by name, err: %v", err)
		common.RespErr(c, common.RespCodeServiceNotFound, common.RespMsgServiceNotFound, nil)
		return
	}

	var remote *url.URL
	if worker.GetNodeName() == conf.AppConfigInstance.NodeName {
		workerPort, ok := tunnel.GetPortManager().GetWorkerPort(c, worker.GetUID())
		if !ok {
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
			return
		}
		remote, err = url.Parse(fmt.Sprintf("http://%s:%d", worker.GetHostName(), workerPort))
		if err != nil {
			logrus.Panic(err)
		}
	} else {
		remote, err = url.Parse(fmt.Sprintf("http://%s:%d",
			conf.AppConfigInstance.TunnelHost, conf.AppConfigInstance.TunnelEntryPort))
		if err != nil {
			logrus.Panic(err)
		}
	}

	accesstoken := c.Request.Header.Get("vvorker-access-token")
	if accesstoken != "" {
		db := database.GetDB()
		var workerToken models.ExternalServerToken
		d := db.Where(&models.ExternalServerToken{
			WorkerUID: worker.UID,
			Token:     accesstoken,
		}).First(&workerToken)
		if d.Error != nil {
			c.AbortWithStatus(403)
			return
		}
		c.Request.Header.Del("vvorker-access-token")
	}

	internaltoken := c.Request.Header.Get("vvorker-internal-token")
	if internaltoken != "" {
		db := database.GetDB()
		tokens := strings.Split(internaltoken, ":")
		if len(tokens) != 2 {
			c.AbortWithStatus(401)
			return
		}
		if tokens[1] != conf.RPCToken {
			c.AbortWithStatus(403)
			return
		}
		if tokens[1] != worker.UID {
			var workerToken models.InternalServerWhiteList
			d := db.Where(&models.InternalServerWhiteList{
				WorkerUID:      worker.UID,
				AllowWorkerUID: tokens[0],
			}).First(&workerToken)
			if d.Error != nil {
				c.AbortWithStatus(403)
				return
			}
		}
		c.Request.Header.Del("vvorker-internal-token")
	}

	var startTime = time.Now()
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(c.Writer, c.Request)
	var endTime = time.Now()
	go func(uid string, status int, method string, path string) {
		db := database.GetDB()
		db.Create(&models.ResponseLog{
			WorkerUID:  uid,
			Status:     c.Writer.Status(),
			Method:     method,
			Path:       path,
			Time:       time.Now(),
			DurationMS: endTime.Sub(startTime).Milliseconds(),
		})
	}(worker.UID, c.Writer.Status(), c.Request.Method, c.Request.URL.Path)
}

type WorkerRequestStatsReq struct {
	WorkerUID string    `json:"worker_uid"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type WorkerRequestStatsResp struct {
	WorkerUID string    `json:"worker_uid"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Total     int       `json:"total"`
	Success   int       `json:"success"`
	Failed    int       `json:"failed"`
}

func GetWorkerRequestStats(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	var req WorkerRequestStatsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, common.RespMsgInvalidRequest, nil)
		return
	}
	var resp WorkerRequestStatsResp
	db := database.GetDB()
	// 使用一条 SQL 语句统计总请求数、成功请求数和失败请求数
	var result struct {
		Total   int64 `gorm:"column:total"`
		Success int64 `gorm:"column:success"`
		Failed  int64 `gorm:"column:failed"`
	}
	query := `
	SELECT 
		COUNT(*) AS total,
		SUM(CASE WHEN status >= 200 AND status < 300 THEN 1 ELSE 0 END) AS success,
		SUM(CASE WHEN status >= 400 THEN 1 ELSE 0 END) AS failed
	FROM 
		response_logs
	WHERE 
		worker_uid = ? AND time BETWEEN ? AND ?
	`
	db.Raw(query, req.WorkerUID, req.StartTime, req.EndTime).Scan(&result)
	total := result.Total
	success := result.Success
	failed := result.Failed

	resp = WorkerRequestStatsResp{
		WorkerUID: req.WorkerUID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Total:     int(total),
		Success:   int(success),
		Failed:    int(failed),
	}
	common.RespOK(c, "ok", resp)
}

type WorkerRequestStatsByTimeReq struct {
	WorkerUID string    `json:"worker_uid"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Interval  string    `json:"interval"` // 时间间隔，如 "1h", "1d"
}

type WorkerRequestStatsByTimeItem struct {
	Time    time.Time `json:"time"`
	Total   int       `json:"total"`
	Success int       `json:"success"`
	Failed  int       `json:"failed"`
}

type WorkerRequestStatsByTimeResp struct {
	WorkerUID string                         `json:"worker_uid"`
	StartTime time.Time                      `json:"start_time"`
	EndTime   time.Time                      `json:"end_time"`
	Interval  string                         `json:"interval"`
	Data      []WorkerRequestStatsByTimeItem `json:"data"`
}

func GetWorkerRequestStatsByTime(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	var req WorkerRequestStatsByTimeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, common.RespMsgInvalidRequest, nil)
		return
	}
	var resp WorkerRequestStatsByTimeResp
	db := database.GetDB()

	// 解析时间间隔
	_, err := time.ParseDuration(req.Interval)
	if err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, "Invalid interval format", nil)
		return
	}

	// 按时间间隔分组查询
	var stats []struct {
		Time    time.Time `gorm:"column:time_bucket"`
		Total   int64     `gorm:"column:total"`
		Success int64     `gorm:"column:success"`
		Failed  int64     `gorm:"column:failed"`
	}
	query := `
	SELECT 
		time_bucket(?, time) AS time_bucket,
		COUNT(*) AS total,
		SUM(CASE WHEN status >= 200 AND status < 300 THEN 1 ELSE 0 END) AS success,
		SUM(CASE WHEN status >= 400 THEN 1 ELSE 0 END) AS failed
	FROM 
		response_logs
	WHERE 
		worker_uid = ? AND time BETWEEN ? AND ?
	GROUP BY 
		time_bucket
	ORDER BY 
		time_bucket
	`
	db.Raw(query, req.Interval, req.WorkerUID, req.StartTime, req.EndTime).Scan(&stats)

	// 转换结果
	var data []WorkerRequestStatsByTimeItem
	for _, stat := range stats {
		data = append(data, WorkerRequestStatsByTimeItem{
			Time:    stat.Time,
			Total:   int(stat.Total),
			Success: int(stat.Success),
			Failed:  int(stat.Failed),
		})
	}

	resp = WorkerRequestStatsByTimeResp{
		WorkerUID: req.WorkerUID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Interval:  req.Interval,
		Data:      data,
	}
	common.RespOK(c, "ok", resp)
}

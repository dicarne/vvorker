package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"runtime/debug"
	"strings"
	"time"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/models"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

type SSOAuthInfo struct {
	UserID    string `json:"user_id"`
	Token     string `json:"token"`
	RealName  string `json:"real_name"`
	Ext       string `json:"ext"`
	Redirect  string `json:"redirect"`
	SetCookie bool   `json:"set_cookie"`
}

func Endpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	host := c.Request.Host
	c.Request.Host = host

	if !strings.HasSuffix(host, conf.AppConfigInstance.WorkerURLSuffix) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	suffix := conf.AppConfigInstance.WorkerURLSuffix
	workerName := host[:len(host)-len(suffix)]
	worker, err := models.AdminGetWorkerByNameSimple(workerName)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Request.Header.Del("server-host")

	var remote *url.URL
	remote, err = url.Parse(fmt.Sprintf("http://%s:%d",
		conf.AppConfigInstance.TunnelHost, conf.AppConfigInstance.TunnelEntryPort))
	if err != nil {
		logrus.Panic(err)
	}

	if worker.EnableAccessControl {
		authed := false
		rules := []models.AccessRule{}
		db := database.GetDB()
		db.Model(&models.AccessRule{}).Where(map[string]interface{}{
			"worker_uid": worker.UID,
			"status":     1,
		}).Order(clause.OrderByColumn{Column: clause.Column{Name: "length"}, Desc: true}).Find(&rules)

		requestPath := c.Request.URL.Path

		for _, rule := range rules {
			if strings.HasPrefix(requestPath, rule.Path) {
				if rule.RuleType == "open" {
					authed = true
					break
				}

				if rule.RuleType == "token" {
					accesstoken := c.Request.Header.Get("vvorker-access-token")
					if accesstoken != "" {
						db := database.GetDB()
						var workerToken models.ExternalServerToken
						d := db.Where(&models.ExternalServerToken{
							WorkerUID: worker.UID,
							Token:     accesstoken,
						}).First(&workerToken)
						if d.Error != nil {
							c.AbortWithStatus(http.StatusForbidden)
							return
						}
						c.Request.Header.Del("vvorker-access-token")
						authed = true
						break
					}
				}

				if rule.RuleType == "internal" {
					internaltoken := c.Request.Header.Get("vvorker-internal-token")
					if internaltoken != "" {
						db := database.GetDB()
						tokens := strings.Split(internaltoken, ":")
						if len(tokens) != 2 {
							c.AbortWithStatus(http.StatusUnauthorized)
							return
						}
						if tokens[1] != conf.RPCToken {
							c.AbortWithStatus(http.StatusForbidden)
							return
						}
						if tokens[1] != worker.UID {
							var workerToken models.InternalServerWhiteList
							d := db.Where(&models.InternalServerWhiteList{
								WorkerUID:      worker.UID,
								AllowWorkerUID: tokens[0],
							}).First(&workerToken)
							if d.Error != nil {
								c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
									"error": "Need Internal Permission",
								})
								return
							}
						}
						c.Request.Header.Del("vvorker-internal-token")
						authed = true
						break
					}
				}

				if rule.RuleType == "sso" && conf.AppConfigInstance.SSOAuthURL != "" {
					// 将c的cookie请求conf.AppConfigInstance.SSOAuthURL，获取用户是否登录的信息
					client := &http.Client{}
					req, err := http.NewRequest("POST", conf.AppConfigInstance.SSOAuthURL, nil)
					if err != nil {
						c.AbortWithStatus(http.StatusUnauthorized)
						return
					}
					cookies := c.Request.Cookies()
					for _, cookie := range cookies {
						if cookie.Name == conf.AppConfigInstance.SSOCookieName {
							req.AddCookie(cookie)
							break
						}
					}
					req.Header.Set(conf.AppConfigInstance.SSOCookieName+"-data", rule.Data)
					req.Header.Set(conf.AppConfigInstance.SSOCookieName+"-worker-uid", worker.UID)
					req.Header.Set(conf.AppConfigInstance.SSOCookieName+"-request-url", c.Request.URL.RequestURI())
					req.Header.Set(conf.AppConfigInstance.SSOCookieName+"-channel", c.Request.Header.Get(conf.AppConfigInstance.SSOCookieName+"-channel"))
					resp, err := client.Do(req)
					if err != nil {
						c.AbortWithStatus(http.StatusUnauthorized)
						return
					}
					defer resp.Body.Close()
					var authInfo SSOAuthInfo
					if err := json.NewDecoder(resp.Body).Decode(&authInfo); err != nil {
						logrus.Errorf("decode error %v", err)
						c.AbortWithStatus(http.StatusForbidden)
						return
					}

					if resp.StatusCode == 401 {
						ssoRedirect := conf.AppConfigInstance.SSORedirectURL

						if !strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
							c.AbortWithStatus(http.StatusForbidden)
							return
						}

						if authInfo.Redirect != "" {
							ssoRedirect = authInfo.Redirect
						}

						// Add original URL as query parameter
						rpath := c.Request.URL.Path

						if conf.AppConfigInstance.WorkerHostMode == "path" {
							rpath = "/" + workerName + rpath

							if conf.AppConfigInstance.SSOBaseURL != "" {
								rpath = conf.AppConfigInstance.SSOBaseURL + rpath
							}

							if conf.AppConfigInstance.WorkerHostPath != "" {
								rpath = "/" + conf.AppConfigInstance.WorkerHostPath + rpath
							}
						} else {
							rpath = fmt.Sprintf("%s://%s%s%s", conf.AppConfigInstance.Scheme, workerName, conf.AppConfigInstance.WorkerURLSuffix, rpath)
						}

						newUrl := fmt.Sprintf("%s?name=%s", ssoRedirect, url.QueryEscape(rpath))
						c.Redirect(http.StatusFound, newUrl)
						return

					}

					if resp.StatusCode != http.StatusOK {
						logrus.Infof("sso auth failed, status code: %d, url: %s", resp.StatusCode, requestPath)
						c.AbortWithStatus(resp.StatusCode)
						return
					}

					if authInfo.SetCookie {
						c.SetCookie(conf.AppConfigInstance.SSOCookieName,
							authInfo.Token,
							conf.AppConfigInstance.SSOCookieAge,
							conf.AppConfigInstance.SSOCookiePath,
							conf.AppConfigInstance.SSOCookieDomain,
							conf.AppConfigInstance.SSOCookieSecure,
							conf.AppConfigInstance.SSOCookieHttpOnly)
					}

					c.Request.Header.Set(conf.AppConfigInstance.SSOCookieName+"-user-id", authInfo.UserID)
					c.Request.Header.Set(conf.AppConfigInstance.SSOCookieName+"-token", authInfo.Token)
					c.Request.Header.Set(conf.AppConfigInstance.SSOCookieName+"-real-name", authInfo.RealName)
					c.Request.Header.Set(conf.AppConfigInstance.SSOCookieName+"-channel", c.Request.Header.Get(conf.AppConfigInstance.SSOCookieName+"-channel"))

					authed = true
					break
				}
			}
		}

		if !authed {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(c.Writer, c.Request)
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

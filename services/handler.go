package services

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"vvorker/authz"
	"vvorker/common"
	"vvorker/conf"
	assets "vvorker/ext/assets/src"
	kv "vvorker/ext/kv/src"
	extmysql "vvorker/ext/mysql/src"
	oss "vvorker/ext/oss/src"
	pgsql "vvorker/ext/pgsql/src"
	"vvorker/models"
	"vvorker/rpc"
	"vvorker/services/access"
	"vvorker/services/agent"
	"vvorker/services/appconf"
	"vvorker/services/auth"
	"vvorker/services/export"
	"vvorker/services/files"
	"vvorker/services/litefs"
	"vvorker/services/node"
	proxyService "vvorker/services/proxy"
	"vvorker/services/resource"
	"vvorker/services/task"
	gentype "vvorker/services/type"
	"vvorker/services/users"
	"vvorker/services/workerd"
	"vvorker/tunnel"
	"vvorker/utils"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"
)

var (
	router    *gin.Engine
	proxy     *gin.Engine
	rpctunnel *gin.Engine
)

func init() {
	router = gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	proxy = gin.Default()
	proxy.Use(modifyProxyRequestHeadersMid)

	rpctunnel = gin.Default()
	rpctunnel.Use(func(c *gin.Context) {
		logrus.Println("rpc tunnel request", c.Request.URL, c.Request.Host, c.Request.Method)
		c.Next()
	})

	router.Use(utils.CORSMiddlewaire(
		fmt.Sprintf("%v://%v", conf.AppConfigInstance.Scheme, conf.AppConfigInstance.CookieDomain),
	))
	if !conf.IsMaster() {
		router.GET("/", func(c *gin.Context) { common.RespOK(c, "ok", nil) })
	}

	api := router.Group("/api")
	registerApi := func(api *gin.RouterGroup) {
		if conf.IsMaster() {
			workerApi := api.Group("/worker", authz.AccessKeyMiddleware(), authz.JWTMiddleware())
			{
				workerApi.GET("/:uid", workerd.GetWorkerEndpoint)
				workerApi.GET("/flush/:uid", workerd.FlushEndpoint)
				workerApi.GET("/run/:uid", workerd.RunWorkerEndpoint)
				workerApi.POST("/create", workerd.CreateEndpoint)
				workerApi.POST("/version/:workerId/:fileId", workerd.NewVersionEndpoint)
				workerApi.POST("/:uid", workerd.UpdateEndpoint)
				workerApi.DELETE("/:uid", workerd.DeleteEndpoint)

				workerApi.GET("/information/:id", workerd.GetWorkerInformationByIDEndpoint)
				workerApi.POST("/information/:id", workerd.UpdateWorkerInformationEndpoint)

				workerApi.POST("/logs/:uid", workerd.GetWorkerLogsEndpoint)
				workerApi.POST("/status", workerd.GetWorkersStatusByUIDEndpoint)

				workerApi.GET("/analyse/group-by-time", proxyService.GetWorkerRequestStatsByTime)
				workerApi.GET("/analyse/by-time", proxyService.GetWorkerRequestStats)

				accessApi := workerApi.Group("/access")
				{
					// 访问令牌子路由
					tokenApi := accessApi.Group("/token")
					{
						tokenApi.POST("/create", access.CreateAccessTokenEndpoint)
						tokenApi.POST("/list", access.ListAccessTokenEndpoint)
						tokenApi.POST("/delete", access.DeleteAccessTokenEndpoint)
					}

					// 内部白名单子路由
					whitelistApi := accessApi.Group("/whitelist")
					{
						whitelistApi.POST("/create", access.CreateInternalWhiteListEndpoint)
						whitelistApi.POST("/list", access.ListInternalWhiteListEndpoint)
						whitelistApi.POST("/update", access.UpdateInternalWhiteListEndpoint)
						whitelistApi.POST("/delete", access.DeleteInternalWhiteListEndpoint)
					}

					controlApi := accessApi.Group("/control")
					{
						controlApi.POST("/update-control", access.UpdateEnableAccessControlEndpoint)
						controlApi.POST("/get-control", access.GetAccessControlEndpoint)
						controlApi.POST("/list-rules", access.ListAccessRuleEndpoint)
						controlApi.POST("/create-rule", access.AddAccessRuleEndpoint)
						controlApi.POST("/delete-rule", access.DeleteAccessRuleEndpoint)
					}
				}

				workerV2 := workerApi.Group("/v2")
				{
					workerV2.POST("/get-worker", workerd.GetWorkerEndpointJSON)
					workerV2.POST("/update-worker", workerd.UpdateEndpointJSON)

					workerV2.POST("/export-workers", export.ExportResourcesConfigEndpoint)
					workerV2.POST("/import-workers", export.ImportResourcesConfigEndpoint)
				}
			}

			workersApi := api.Group("/workers", authz.AccessKeyMiddleware(), authz.JWTMiddleware())
			{
				workersApi.GET("/flush", workerd.FlushAllEndpoint)
				workersApi.GET("/:offset/:limit", workerd.GetWorkersEndpoint)
			}
			userApi := api.Group("/user", authz.JWTMiddleware())
			{
				userApi.GET("/info", auth.GetUserEndpoint)
				userApi.POST("/create-access-key", access.CreateAccessKeyEndpoint)
				userApi.POST("/access-keys", access.GetAccessKeysEndpoint)
				userApi.POST("/delete-access-key", access.DeleteAccessKeyEndpoint)

				// 注册用户管理相关路由

			}

			adminAPI := api.Group("/admin", authz.AccessKeyMiddleware(), authz.JWTMiddleware())
			{
				users.RegisterRoutes(adminAPI)
			}

			nodeAPI := api.Group("/node", authz.AccessKeyMiddleware(), authz.JWTMiddleware())
			{
				nodeAPI.GET("/all", node.UserGetNodesEndpoint)
				nodeAPI.GET("/sync/:nodename", node.SyncNodeEndpoint)
				nodeAPI.DELETE("/:nodename", node.LeaveEndpoint)
			}
			fileAPI := api.Group("/file", authz.AccessKeyMiddleware(), authz.JWTMiddleware())
			{
				fileAPI.POST("/upload", files.UploadFileEndpoint)
				fileAPI.GET("/get/:fileId", files.GetFileEndpoint)
			}
			api.GET("/allworkers", authz.AccessKeyMiddleware(), authz.JWTMiddleware(), workerd.GetAllWorkersEndpoint)
			api.GET("/vvorker/config", appconf.GetEndpoint)
			api.POST("/auth/register", auth.RegisterEndpoint)
			api.POST("/auth/login", auth.LoginEndpoint)
			api.GET("/auth/logout", authz.JWTMiddleware(), auth.LogoutEndpoint)

		}
		agentAPI := api.Group("/agent")
		{
			if conf.IsMaster() {
				agentAPI.POST("/sync", authz.AgentAuthz(), workerd.AgentSyncWorkers)
				agentAPI.POST("/add", authz.AgentAuthz(), node.AddEndpoint)
				agentAPI.GET("/nodeinfo", authz.AgentAuthz(), node.GetNodeInfoEndpoint)
				agentAPI.POST("/fill-worker-config", authz.AgentAuthz(), workerd.FillWorkerConfig)
			} else {
				agentAPI.POST("/notify", authz.AgentAuthz(), agent.NotifyEndpoint)
			}
		}
		extAPI := api.Group("/ext")
		{
			ossAPI := extAPI.Group("/oss")
			{
				ossAPI.POST("/upload", authz.AgentAuthz(), oss.UploadFile)
				ossAPI.POST("/download", authz.AgentAuthz(), oss.DownloadFile)
				ossAPI.POST("/list-buckets", authz.AgentAuthz(), oss.ListBuckets)
				ossAPI.POST("/delete", authz.AgentAuthz(), oss.DeleteFile)
				ossAPI.POST("/list-objects", authz.AgentAuthz(), oss.ListObjects)

				ossAPI.POST("/initiate-multipart-upload", authz.AgentAuthz(), oss.InitiateMultipartUpload)
				ossAPI.POST("/upload-part", authz.AgentAuthz(), oss.UploadPart)
				ossAPI.POST("/complete-multipart-upload", authz.AgentAuthz(), oss.CompleteMultipartUpload)
				ossAPI.POST("/abort-multipart-upload", authz.AgentAuthz(), oss.AbortMultipartUpload)

				if conf.IsMaster() {
					ossAPI.POST("/create-resource", authz.AccessKeyMiddleware(), authz.JWTMiddleware(), oss.CreateNewOSSResourcesEndpoint)
					ossAPI.POST("/delete-resource", authz.AccessKeyMiddleware(), authz.JWTMiddleware(), oss.DeleteOSSResourcesEndpoint)
				}
			}
			pgsqlAPI := extAPI.Group("/pgsql")
			{
				if conf.IsMaster() {
					pgsqlAPI.POST("/create-resource", authz.AccessKeyMiddleware(), authz.JWTMiddleware(), pgsql.CreateNewPostgreSQLResourcesEndpoint)
					pgsqlAPI.POST("/delete-resource", authz.AccessKeyMiddleware(), authz.JWTMiddleware(), pgsql.DeletePostgreSQLResourcesEndpoint)
					pgsqlAPI.POST("/migrate", authz.AccessKeyMiddleware(), authz.JWTMiddleware(), pgsql.UpdateMigrate)
				}
			}
			mysqlAPI := extAPI.Group("/mysql")
			{
				if conf.IsMaster() {
					mysqlAPI.POST("/create-resource", authz.AccessKeyMiddleware(), authz.JWTMiddleware(), extmysql.CreateNewMySQLResourcesEndpoint)
					mysqlAPI.POST("/delete-resource", authz.AccessKeyMiddleware(), authz.JWTMiddleware(), extmysql.DeleteMySQLResourcesEndpoint)
					mysqlAPI.POST("/migrate", authz.AccessKeyMiddleware(), authz.JWTMiddleware(), extmysql.UpdateMigrate)
				}
			}
			kvAPI := extAPI.Group("/kv")
			{
				if conf.IsMaster() {
					kvAPI.POST("/create-resource", authz.AccessKeyMiddleware(), authz.JWTMiddleware(), kv.CreateKVResourcesEndpoint)
					kvAPI.POST("/delete-resource", authz.AccessKeyMiddleware(), authz.JWTMiddleware(), kv.DeleteKVResourcesEndpoint)
				}
			}
			assetsAPI := extAPI.Group("/assets")
			{
				if conf.IsMaster() {
					assetsAPI.POST("/create-assets", authz.AccessKeyMiddleware(), authz.JWTMiddleware(), assets.UploadAssetsEndpoint)
					assetsAPI.GET("/get-assets", authz.AgentAuthz(), assets.GetAssetsEndpoint)
				}
			}
			taskAPI := extAPI.Group("/task")
			{
				if conf.IsMaster() {
					taskAPI.POST("/create", authz.AgentAuthz(), task.CreateTaskEndpoint)
					taskAPI.POST("/cancel", authz.AccessKeyMiddleware(), authz.JWTMiddleware(), task.CancelTaskEndpoint)
					taskAPI.POST("/check", authz.AgentAuthz(), task.CheckInterruptTaskEndpoint)
					taskAPI.POST("/log", authz.AgentAuthz(), task.LogTaskEndpoint)
					taskAPI.POST("/logs", authz.AccessKeyMiddleware(), authz.JWTMiddleware(), task.GetLogsEndpoint)
					taskAPI.POST("/complete", authz.AgentAuthz(), task.CompleteTaskEndpoint)
					taskAPI.POST("/list", authz.AccessKeyMiddleware(), authz.JWTMiddleware(), task.ListTaskEndpoint)
				}
			}

			if conf.IsMaster() {
				extAPI.POST("/list", authz.AccessKeyMiddleware(), authz.JWTMiddleware(), resource.ListResourceEndpoint)

				extAPI.POST("/types", authz.AccessKeyMiddleware(), gentype.GenerateTypes)
			}
		}
	}
	registerApi(api)

	if conf.AppConfigInstance.WorkerHostMode == "path" && conf.AppConfigInstance.WorkerHostPath != "" {
		if conf.IsMaster() {
			api2 := router.Group("/" + conf.AppConfigInstance.WorkerHostPath + "/api")
			registerApi(api2)
		}

	}
	// rpctunnel.Any("/*proxyPath", proxyService.Endpoint)
	// rpctunnel.NoRoute(proxyService.HandleConnect)
	proxy.Any("/*proxyPath", proxyService.Endpoint)
}

func InitTunnel(wg *conc.WaitGroup) {
	if conf.IsMaster() {
		wg.Go(tunnel.Serve)
		wg.Go(tunnel.InitSelfCliet)

		wg.Go(func() { tunnel.GetClient().Run(context.Background()) })
	} else {
		wg.Go(RegisterNodeToMaster)
		wg.Go(func() { tunnel.GetClient().Run(context.Background()) })
	}
	wg.Go(litefs.InitTunnel)
	wg.Go(litefs.RunService)
}

// initTunnelService 初始化隧道服务和访客
func initTunnelService(serviceName string, servicePort int, visitorPort int) error {
	err := tunnel.GetClient().AddService(serviceName, servicePort)
	if err != nil {
		logrus.WithError(err).Errorf("init tunnel for %s service error", serviceName)
		return err
	}
	err = tunnel.GetClient().AddVisitor(serviceName, visitorPort)
	if err != nil {
		logrus.WithError(err).Errorf("init tunnel for %s visitor failed", serviceName)
		return err
	}
	return nil
}

func Run(f embed.FS) {
	if err := os.MkdirAll(conf.AppConfigInstance.WorkerdDir, 0755); err != nil {
		logrus.WithError(err).Errorf("Failed to create workerd directory: %s", conf.AppConfigInstance.WorkerdDir)
		return
	}

	os.RemoveAll(
		filepath.Join(
			conf.AppConfigInstance.WorkerdDir,
			"lib",
		),
	)

	wg := conc.NewWaitGroup()
	defer wg.Wait()

	InitTunnel(wg)
	wg.Go(func() {
		proxy.Run(fmt.Sprintf("%v:%d", conf.AppConfigInstance.ListenAddr, conf.AppConfigInstance.WorkerPort))
	})
	wg.Go(database.InitDB)
	wg.Go(models.MigrateNormalModel)
	if conf.IsMaster() {
		HandleStaticFile(f)
	}
	wg.Go(func() {
		// 将数据库远程端口代理到master临时本地端口
		if conf.IsMaster() {
			proxyService.InitReverseProxy(fmt.Sprintf("%v:%d", conf.AppConfigInstance.ServerPostgreHost, conf.AppConfigInstance.ServerPostgrePort), fmt.Sprintf(":%d", conf.AppConfigInstance.LocalTMPPostgrePort))
			proxyService.InitReverseProxy(fmt.Sprintf("%v:%d", conf.AppConfigInstance.ServerRedisHost, conf.AppConfigInstance.ServerRedisPort), fmt.Sprintf(":%d", conf.AppConfigInstance.LocalTMPRedisPort))
		}
		// 将master；临时本地端口代理到worker本地端口
		initTunnelService(conf.AppConfigInstance.NodeName+"postgresql", conf.AppConfigInstance.LocalTMPPostgrePort, conf.AppConfigInstance.ClientPostgrePort)
		initTunnelService(conf.AppConfigInstance.NodeName+"redis", conf.AppConfigInstance.LocalTMPRedisPort, conf.AppConfigInstance.ClientRedisPort)
		initTunnelService(conf.AppConfigInstance.NodeName+"mysql", conf.AppConfigInstance.LocalTMPMySQLPort, conf.AppConfigInstance.ClientMySQLPort)
	})
	wg.Go(func() {
		router.Run(fmt.Sprintf("%v:%d", conf.AppConfigInstance.ListenAddr, conf.AppConfigInstance.APIPort))
	})
	wg.Go(func() {
		rpctunnel.Run(fmt.Sprintf("%v:%d", conf.AppConfigInstance.ListenAddr, conf.AppConfigInstance.InternalRPCPort))
	})
}

func HandleStaticFile(f embed.FS) {
	fp, err := fs.Sub(f, "www/out")
	if err != nil {
		logrus.Panic(err)
	}
	prefix := ""
	if conf.AppConfigInstance.WorkerHostMode == "path" {
		prefix = "/" + conf.AppConfigInstance.WorkerHostPath
	}
	router.StaticFileFS(prefix+"/404", "404.html", http.FS(fp))
	router.StaticFileFS(prefix+"/login", "login.html", http.FS(fp))
	router.StaticFileFS(prefix+"/admin", "admin.html", http.FS(fp))
	router.StaticFileFS(prefix+"/register", "register.html", http.FS(fp))
	router.StaticFileFS(prefix+"/worker", "worker.html", http.FS(fp))
	router.StaticFileFS(prefix+"/index", "index.html", http.FS(fp))
	router.StaticFileFS(prefix+"/nodes", "nodes.html", http.FS(fp))
	router.StaticFileFS(prefix+"/user", "user.html", http.FS(fp))

	router.StaticFileFS(prefix+"/sql", "sql.html", http.FS(fp))
	router.StaticFileFS(prefix+"/oss", "oss.html", http.FS(fp))
	router.StaticFileFS(prefix+"/kv", "kv.html", http.FS(fp))
	router.StaticFileFS(prefix+"/task", "task.html", http.FS(fp))
	router.StaticFileFS(prefix+"/logs", "logs.html", http.FS(fp))
	router.StaticFileFS(prefix+"/users", "users.html", http.FS(fp))
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if conf.AppConfigInstance.WorkerHostMode == "path" && conf.AppConfigInstance.WorkerHostPath != "" {
			path = strings.Replace(c.Request.URL.Path, "/"+conf.AppConfigInstance.WorkerHostPath, "", 1)
		}
		if conf.AppConfigInstance.AdminAPIProxy {
			if conf.AppConfigInstance.WorkerHostMode != "path" {
				logrus.Println("admin api proxy only support path mode")
			}
			_, err := http.FS(fp).Open(path)
			if err != nil {
				modifyProxyRequestHeaders(c)
				proxyService.Endpoint(c)
				return
			}
		}
		c.FileFromFS(path, http.FS(fp))
	})
}

func RegisterNodeToMaster() {
	if conf.IsMaster() {
		return
	}
	if conf.AppConfigInstance.LitefsEnabled {
		utils.WaitForPort("localhost", conf.AppConfigInstance.LitefsPrimaryPort)
	}
	for {
		logrus.Info("Registering node to master...")
		self, err := rpc.GetNode(conf.AppConfigInstance.MasterEndpoint)
		if err != nil || self == nil {
			err := rpc.AddNode(conf.AppConfigInstance.MasterEndpoint)
			if err != nil {
				logrus.WithError(err).Error("Add node failed.. retrying for 5 seconds")
				time.Sleep(5 * time.Second)
			} else {
				logrus.Info("Node added successfully")
			}
			continue
		} else {
			logrus.Info("Node already exists")
			conf.AppConfigInstance.NodeID = self.UID
		}
		tun, err := tunnel.GetClient().Query(conf.AppConfigInstance.NodeID)
		if err != nil || tun == nil {
			logrus.Warnf("Query tunnel failed, err: %v, try to add tunnel", err)
			tunnel.GetClient().Add(conf.AppConfigInstance.NodeID, utils.NodeHostPrefix(
				conf.AppConfigInstance.NodeName, conf.AppConfigInstance.NodeID),
				int(conf.AppConfigInstance.APIPort))
		} else {
			logrus.Info("Tunnel already exists, skip adding")
		}
		if conf.AppConfigInstance.EnableAutoSync {
			agent.SyncCall()
		}
		time.Sleep(30 * time.Second)
	}
}

func modifyProxyRequestHeaders(c *gin.Context) {
	if conf.AppConfigInstance.WorkerHostMode == "path" {
		// 此时，URL的第一段会被提取出来作为host name，并在传下去的url中去掉这一段
		// 按照 / 切分
		url := c.Request.URL.Path
		// 去掉开头的 /
		if len(url) > 0 && url[0] == '/' {
			url = url[1:]
		}
		// 按 / 分割路径
		if conf.AppConfigInstance.WorkerHostPath == "" {
			parts := strings.SplitN(url, "/", 2)
			if len(parts) >= 2 && parts[0] != "" {
				host := parts[0] + conf.AppConfigInstance.WorkerURLSuffix
				c.Request.Header.Set("Host", host)
				c.Request.Host = host
				// 去掉第一段后的路径
				if len(parts) > 1 {
					c.Request.URL.Path = "/" + parts[1]
				} else {
					c.Request.URL.Path = "/"
				}
			}
		} else {
			parts := strings.SplitN(url, "/", 3)
			if len(parts) >= 3 && parts[0] != "" {
				host := parts[1] + conf.AppConfigInstance.WorkerURLSuffix
				c.Request.Header.Set("Host", host)
				c.Request.Host = host
				// 去掉第二段后的路径
				if len(parts) > 2 {
					c.Request.URL.Path = "/" + parts[1] + "/" + parts[2]
				} else {
					c.Request.URL.Path = "/"
				}
			}

		}

	} else {
		host := c.Request.Header.Get("Server-Host")
		if host != "" {
			c.Request.Header.Set("Host", host)
			c.Request.Host = host
		}
	}
}

func modifyProxyRequestHeadersMid(c *gin.Context) {
	modifyProxyRequestHeaders(c)
	c.Next()
}

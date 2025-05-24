package services

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"
	"vorker/authz"
	"vorker/conf"
	oss "vorker/ext/oss/src"
	"vorker/models"
	"vorker/rpc"
	"vorker/services/agent"
	"vorker/services/appconf"
	"vorker/services/auth"
	"vorker/services/files"
	"vorker/services/litefs"
	"vorker/services/node"
	proxyService "vorker/services/proxy"
	"vorker/services/workerd"
	"vorker/tunnel"
	"vorker/utils"
	"vorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"
)

var (
	router *gin.Engine
	proxy  *gin.Engine
)

func init() {
	router = gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	proxy = gin.Default()
	router.Use(utils.CORSMiddlewaire(
		fmt.Sprintf("%v://%v", conf.AppConfigInstance.Scheme, conf.AppConfigInstance.CookieDomain),
	))
	if !conf.IsMaster() {
		router.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"code": 0, "msg": "ok"}) })
	}

	api := router.Group("/api")
	{
		if conf.IsMaster() {
			workerApi := api.Group("/worker", authz.JWTMiddleware())
			{
				workerApi.GET("/:uid", workerd.GetWorkerEndpoint)
				workerApi.GET("/flush/:uid", workerd.FlushEndpoint)
				workerApi.GET("/run/:uid", workerd.RunWorkerEndpoint)
				workerApi.POST("/create", workerd.CreateEndpoint)
				workerApi.POST("/version/:workerId/:fileId", workerd.NewVersionEndpoint)
				workerApi.PATCH("/:uid", workerd.UpdateEndpoint)
				workerApi.DELETE("/:uid", workerd.DeleteEndpoint)
			}
			workersApi := api.Group("/workers", authz.JWTMiddleware())
			{
				workersApi.GET("/flush", workerd.FlushAllEndpoint)
				workersApi.GET("/:offset/:limit", workerd.GetWorkersEndpoint)
			}
			userApi := api.Group("/user", authz.JWTMiddleware())
			{
				userApi.GET("/info", auth.GetUserEndpoint)
			}
			nodeAPI := api.Group("/node")
			{
				nodeAPI.GET("/all", authz.JWTMiddleware(), node.UserGetNodesEndpoint)
				nodeAPI.GET("/sync/:nodename", authz.JWTMiddleware(), node.SyncNodeEndpoint)
				nodeAPI.DELETE("/:nodename", authz.JWTMiddleware(), node.LeaveEndpoint)
			}
			fileAPI := api.Group("/file", authz.JWTMiddleware())
			{
				fileAPI.POST("/upload", files.UploadFileEndpoint)
				fileAPI.GET("/get/:fileId", files.GetFileEndpoint)
			}
			api.GET("/allworkers", authz.JWTMiddleware(), workerd.GetAllWorkersEndpoint)
			api.GET("/vorker/config", appconf.GetEndpoint)
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
				ossAPI.POST("/upload", oss.UploadFile)
				ossAPI.POST("/download", oss.DownloadFile)
				ossAPI.POST("/list-buckets", oss.ListBuckets)
				ossAPI.POST("/delete", oss.DeleteFile)
				ossAPI.POST("/list-objects", oss.ListObjects)

				if conf.IsMaster() {
					ossAPI.POST("/create-resource", oss.CreateNewOSSResourcesEndpoint)
				}
			}
		}
	}

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
		initTunnelService("redis", conf.AppConfigInstance.ServerRedisPort, conf.AppConfigInstance.ClientRedisPort)
		initTunnelService("postgresql", conf.AppConfigInstance.ServerPostgresPort, conf.AppConfigInstance.ClientPostgresPort)
	})
	wg.Go(func() {
		router.Run(fmt.Sprintf("%v:%d", conf.AppConfigInstance.ListenAddr, conf.AppConfigInstance.APIPort))
	})
}

func HandleStaticFile(f embed.FS) {
	fp, err := fs.Sub(f, "www/out")
	if err != nil {
		logrus.Panic(err)
	}
	router.StaticFileFS("/404", "404.html", http.FS(fp))
	router.StaticFileFS("/login", "login.html", http.FS(fp))
	router.StaticFileFS("/admin", "admin.html", http.FS(fp))
	router.StaticFileFS("/register", "register.html", http.FS(fp))
	router.StaticFileFS("/worker", "worker.html", http.FS(fp))
	router.StaticFileFS("/index", "index.html", http.FS(fp))
	router.StaticFileFS("/nodes", "nodes.html", http.FS(fp))
	router.NoRoute(func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(fp))
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

package features

import (
	"vvorker/common"
	"vvorker/conf"

	"github.com/gin-gonic/gin"
)

type Feature struct {
	Name   string `json:"name"`
	Enable bool   `json:"enable"`
}

func ListFeaturesEndpoint(c *gin.Context) {
	features := []Feature{
		{
			Name:   "mysql",
			Enable: conf.AppConfigInstance.EnableMySQL,
		},
		{
			Name:   "pgsql",
			Enable: conf.AppConfigInstance.EnablePgSQL,
		},
		{
			Name:   "redis",
			Enable: conf.AppConfigInstance.EnableRedis,
		},
		{
			Name:   "minio",
			Enable: conf.AppConfigInstance.EnableMinIO,
		},
	}
	common.RespOK(c, "ok", features)
}

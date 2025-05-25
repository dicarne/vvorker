package authz

import (
	"fmt"
	"vvorker/conf"
	"vvorker/defs"
	"vvorker/utils/secret"

	"github.com/gin-gonic/gin"
)

func AgentAuthz() func(c *gin.Context) {
	return func(c *gin.Context) {
		ssecret := c.Request.Header.Get(defs.HeaderNodeSecret)
		name := c.Request.Header.Get(defs.HeaderNodeName)
		querySec := c.Copy().Query("secret")
		queryName := c.Copy().Query("name")
		if (ssecret == "" || name == "") && (querySec == "" || queryName == "") {
			c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
			return
		}

		if secret.MD5(fmt.Sprintf("%s%s", name, conf.AppConfigInstance.AgentSecret)) == ssecret {
			c.Set(defs.KeyNodeName, name)
			c.Next()
			return
		}

		if secret.MD5(fmt.Sprintf("%s%s", queryName, conf.AppConfigInstance.AgentSecret)) == querySec {
			c.Set(defs.KeyNodeName, queryName)
			c.Next()
			return
		}

		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
	}
}

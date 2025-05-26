package authz

import (
	"strings"
	"vvorker/common"
	"vvorker/services/access"

	"github.com/gin-gonic/gin"
)

func AccessKeyMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenOrigin := c.Request.Header.Get(common.AuthorizationKey)
		tokenList := strings.Split(tokenOrigin, " ")
		if len(tokenList) != 2 {
			c.Next()
			return
		}
		tokenStr := tokenList[1]

		if tokenStr == "" {
			c.Next()
			return
		}

		if tokenStr[:4] == "ac::" {
			if uid, err := access.AccessKeyToUserID(tokenStr); err == nil {
				c.Set(common.UIDKey, uint(uid))
				c.Set("JWT_PASS", true)
				c.Next()
				return
			}
		}
		c.Next()
	}
}

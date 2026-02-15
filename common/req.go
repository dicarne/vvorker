package common

import (
	"github.com/gin-gonic/gin"
	"vvorker/entities"
)

type Request interface {
	*entities.DeleteWorkerRequest | *entities.LoginRequest | *entities.RegisterRequest |
		*entities.NotifyEventRequest
	Validate() bool
}

// GetUID 从上下文中获取 uid，如果不存在则返回 false
func GetUID(c *gin.Context) (uint64, bool) {
	uid := uint64(c.GetUint(UIDKey))
	return uid, uid != 0
}

// RequireUID 从上下文中获取 uid，如果不存在则返回错误
func RequireUID(c *gin.Context) (uint64, bool) {
	uid, ok := GetUID(c)
	if !ok {
		RespErr(c, RespCodeInvalidRequest, "uid is required", nil)
		return 0, false
	}
	return uid, true
}

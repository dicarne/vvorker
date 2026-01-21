package common

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RecoveryMiddleware 是一个全局的 panic 恢复中间件
// 它会自动捕获所有未处理的 panic 并返回统一的错误响应
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// 记录错误日志和堆栈信息
				logrus.Errorf("Recovered in handler: %+v, stack: %+v", r, string(debug.Stack()))

				// 返回统一的错误响应
				RespErr(c, RespCodeInternalError, RespMsgInternalError, nil)

				// 阻止后续的处理
				c.Abort()
			}
		}()

		// 继续处理请求
		c.Next()
	}
}

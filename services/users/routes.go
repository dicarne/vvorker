package users

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册用户管理相关的路由
func RegisterRoutes(router *gin.RouterGroup) {
	// 获取用户列表
	router.GET("/users", ListUsersEndpoint)

	// 获取单个用户信息
	router.GET("/users/:id", GetUserEndpoint)

	// 搜索用户
	router.GET("/users/search", SearchUsersEndpoint)

	// 创建用户
	router.POST("/users", CreateUserEndpoint)

	// 更新用户
	router.POST("/users/:id", UpdateUserEndpoint)

	// 删除用户
	router.DELETE("/users/:id", DeleteUserEndpoint)

	// 更新用户状态
	router.POST("/status", UpdateUserStatusEndpoint)

	// 批量更新用户状态
	router.POST("/batch-status", BatchUpdateUserStatusEndpoint)
}

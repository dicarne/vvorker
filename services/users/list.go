package users

import (
	"strconv"
	"vvorker/common"
	"vvorker/models"

	"github.com/gin-gonic/gin"
)

// ListUsersEndpoint 获取用户列表
func ListUsersEndpoint(c *gin.Context) {
	// 检查是否是管理员
	if !IsAdmin(c) {
		common.RespErr(c, common.RespCodeAuthErr, "unauthorized", nil)
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// 获取用户列表
	users, err := models.ListUsers(page, pageSize)
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	// 获取总数
	total, err := models.CountUsers()
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	// 清除密码字段
	for _, u := range users {
		u.Password = ""
	}

	common.RespOK(c, "", gin.H{
		"users":    users,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

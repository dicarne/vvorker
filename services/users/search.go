package users

import (
	"strconv"
	"vvorker/common"

	"github.com/gin-gonic/gin"

	"vvorker/models"
)

// SearchUsersEndpoint 搜索用户
// @Summary 搜索用户
// @Description 根据条件搜索用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param query query string false "搜索关键词"
// @Param role query string false "用户角色"
// @Param status query int false "用户状态"
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认20"
// @Success 200 {object} handler.Response
// @Router /api/users/search [get]
func SearchUsersEndpoint(c *gin.Context) {
	// 检查权限，只有管理员可以搜索用户
	if !IsAdmin(c) {
		common.RespErr(c, common.RespCodeUserNotAdmin, "权限不足", nil)
		return
	}

	// 获取查询参数
	query := c.Query("query")
	role := c.Query("role")
	statusStr := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	// 转换参数类型
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	status := 0 // 0表示不筛选状态
	if statusStr != "" {
		status, err = strconv.Atoi(statusStr)
		if err != nil {
			status = 0
		}
	}

	// 执行搜索
	users, total, err := models.AdminSearchUsers(query, role, status, page, pageSize)
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, "搜索用户失败: "+err.Error(), nil)
		return
	}

	// 返回结果
	common.RespOK(c, "", gin.H{
		"users": users,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func IsAdmin(c *gin.Context) bool {
	userId := c.GetUint(common.UIDKey)
	user, err := models.GetUserByUserID(userId)
	if err != nil {
		return false
	}
	return user.Role == "admin"
}

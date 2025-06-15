package users

import (
	"strconv"
	"vvorker/common"
	"vvorker/models"

	"github.com/gin-gonic/gin"
)

// GetUserEndpoint 获取单个用户信息
func GetUserEndpoint(c *gin.Context) {
	// 检查是否是管理员
	if !IsAdmin(c) {
		common.RespErr(c, common.RespCodeAuthErr, "unauthorized", nil)
		return
	}

	// 获取用户ID
	userID, err := strconv.ParseUint(c.Query("id"), 10, 32)
	if err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, "invalid user id", nil)
		return
	}

	// 获取用户信息
	userInfo, err := models.AdminGetUserByID(uint(userID))
	if err != nil {
		common.RespErr(c, common.RespCodeDBErr, err.Error(), nil)
		return
	}
	if userInfo == nil {
		common.RespErr(c, common.RespCodeNotFound, "user not found", nil)
		return
	}

	// 清除密码
	userInfo.Password = ""

	common.RespOK(c, "success", userInfo)
}

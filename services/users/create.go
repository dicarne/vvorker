package users

import (
	"vvorker/common"
	"vvorker/models"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

// CreateUserEndpoint 创建新用户
func CreateUserEndpoint(c *gin.Context) {
	// 检查是否是管理员
	if !IsAdmin(c) {
		common.RespErr(c, common.RespCodeAuthErr, "unauthorized", nil)
		return
	}

	// 解析请求体
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, "invalid request: "+err.Error(), nil)
		return
	}

	if req.Username == "" || req.Password == "" || req.Email == "" {
		common.RespErr(c, common.RespCodeInvalidRequest, "invalid request", nil)
		return
	}

	// 检查用户名是否已存在
	existingUser, _ := models.AdminGetUserByUsername(req.Username)
	if existingUser != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, "username already exists", nil)
		return
	}

	// 创建用户
	newUser, err := models.AdminCreateUser(req.Username, req.Password, req.Email)
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	// 清除密码后返回
	newUser.Password = ""
	common.RespOK(c, "user created successfully", newUser)
}

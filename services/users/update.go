package users

import (
	"net/http"
	"strconv"
	"vvorker/common"
	"vvorker/models"

	"github.com/gin-gonic/gin"
)

type UpdateUserRequest struct {
	Password string `json:"password"`
	Role     string `json:"role"`
	Status   int    `json:"status"`
}

// UpdateUserEndpoint 更新用户信息
func UpdateUserEndpoint(c *gin.Context) {
	// 检查是否是管理员
	if !IsAdmin(c) {
		common.RespErr(c, http.StatusUnauthorized, "", gin.H{
			"code": 1,
			"msg":  "unauthorized",
		})
		return
	}

	// 获取用户ID
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		common.RespErr(c, http.StatusBadRequest, "", gin.H{
			"code": 1,
			"msg":  "invalid user id",
		})
		return
	}

	// 解析请求体
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, http.StatusBadRequest, "", gin.H{
			"code": 1,
			"msg":  "invalid request: " + err.Error(),
		})
		return
	}

	// 获取要更新的用户
	userToUpdate, err := models.AdminGetUserByID(uint(userID))
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "", gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}
	if userToUpdate == nil {
		common.RespErr(c, http.StatusNotFound, "", gin.H{
			"code": 1,
			"msg":  "user not found",
		})
		return
	}

	if req.Password != "" {
		userToUpdate.Password = req.Password
	}
	if req.Role != "" {
		userToUpdate.Role = req.Role
	}
	if req.Status != 0 {
		userToUpdate.Status = req.Status
	}

	// 保存更新
	if err := models.AdminUpdateUser(userToUpdate); err != nil {
		common.RespErr(c, http.StatusInternalServerError, "", gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}

	// 清除密码后返回
	userToUpdate.Password = ""
	common.RespOK(c, "", gin.H{
		"code": 0,
		"data": userToUpdate,
		"msg":  "user updated successfully",
	})
}

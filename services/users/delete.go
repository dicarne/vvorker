package users

import (
	"strconv"
	"vvorker/common"
	"vvorker/models"

	"github.com/gin-gonic/gin"
)

// DeleteUserEndpoint 删除用户
func DeleteUserEndpoint(c *gin.Context) {
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

	// 不允许删除自己
	if uint(userID) == c.GetUint(common.UIDKey) {
		common.RespErr(c, common.RespCodeInvalidRequest, "cannot delete yourself", nil)
		return
	}

	// 检查用户是否存在
	userToDelete, err := models.AdminGetUserByID(uint(userID))
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}
	if userToDelete == nil {
		common.RespErr(c, common.RespCodeNotFound, "user not found", nil)
		return
	}

	// 删除用户
	if err := models.DeleteUser(uint(userID)); err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	common.RespOK(c, "user deleted successfully", nil)
}

package users

import (
	"net/http"
	"vvorker/common"

	"github.com/gin-gonic/gin"

	"vvorker/models"
)

// UpdateUserStatusRequest 更新用户状态请求
type UpdateUserStatusRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	Status int  `json:"status" binding:"required"`
}

// BatchUpdateUserStatusRequest 批量更新用户状态请求
type BatchUpdateUserStatusRequest struct {
	UserIDs []uint `json:"user_ids" binding:"required"`
	Status  int    `json:"status" binding:"required"`
}

// UpdateUserStatusEndpoint 更新用户状态
// @Summary 更新用户状态
// @Description 更新单个用户的状态
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body UpdateUserStatusRequest true "更新用户状态请求"
// @Success 200 {object} handler.Response
// @Router /api/users/status [post]
func UpdateUserStatusEndpoint(c *gin.Context) {
	// 检查权限，只有管理员可以更新用户状态
	if !IsAdmin(c) {
		common.RespErr(c, common.RespCodeUserNotAdmin, "权限不足", nil)
		return
	}

	var req UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, "无效的请求参数: "+err.Error(), nil)
		return
	}

	// 不允许修改ID为1的用户（假设是超级管理员）
	if req.UserID == 1 {
		common.RespErr(c, common.RespCodeUserNotAdmin, "不允许修改超级管理员状态", nil)
		return
	}

	// 更新用户状态
	if err := models.AdminUpdateUserStatus(req.UserID, req.Status); err != nil {
		common.RespErr(c, common.RespCodeInternalError, "更新用户状态失败: "+err.Error(), nil)
		return
	}

	common.RespOK(c, "更新用户状态成功", nil)
}

// BatchUpdateUserStatusEndpoint 批量更新用户状态
// @Summary 批量更新用户状态
// @Description 批量更新多个用户的状态
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body BatchUpdateUserStatusRequest true "批量更新用户状态请求"
// @Success 200 {object} handler.Response
// @Router /api/users/batch-status [post]
func BatchUpdateUserStatusEndpoint(c *gin.Context) {
	// 检查权限，只有管理员可以批量更新用户状态
	if !IsAdmin(c) {
		common.RespErr(c, http.StatusForbidden, "权限不足", gin.H{})
		return
	}

	var req BatchUpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, http.StatusBadRequest, "无效的请求参数: "+err.Error(), gin.H{})
		return
	}

	// 检查是否包含ID为1的用户（假设是超级管理员）
	for _, id := range req.UserIDs {
		if id == 1 {
			common.RespErr(c, http.StatusForbidden, "不允许修改超级管理员状态", gin.H{})
			return
		}
	}

	// 批量更新用户状态
	if err := models.AdminBatchUpdateUserStatus(req.UserIDs, req.Status); err != nil {
		common.RespErr(c, http.StatusInternalServerError, "批量更新用户状态失败: "+err.Error(), gin.H{})
		return
	}

	common.RespOK(c, "批量更新用户状态成功", gin.H{})
}

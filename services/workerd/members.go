package workerd

import (
	"vvorker/common"
	"vvorker/models"
	permissions "vvorker/utils/permissions"

	"github.com/gin-gonic/gin"
)

type AddMemberRequest struct {
	WorkerUID string `json:"worker_uid"`
	UserName  string `json:"user_name"`
}

type RemoveMemberRequest struct {
	WorkerUID string `json:"worker_uid"`
	UserID    uint64 `json:"user_id"`
}

// AddMemberEndpoint 添加协作者
func AddMemberEndpoint(c *gin.Context) {

	var req AddMemberRequest
	if err := c.BindJSON(&req); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}

	if req.WorkerUID == "" || req.UserName == "" {
		common.RespErr(c, common.RespCodeInvalidRequest, "worker_uid and user_name are required", nil)
		return
	}

	userID := c.GetUint(common.UIDKey)

	// 验证操作者是否可以管理成员
	_, err := permissions.CanManageWorkerMembers(c, uint64(userID), req.WorkerUID)
	if err != nil {
		return
	}

	// 获取要添加的用户
	targetUser, err := models.GetUserByUserName(req.UserName)
	if err != nil {
		common.RespErr(c, common.RespCodeNotFound, "user not found", nil)
		return
	}

	// 获取当前用户信息
	currentUser, err := models.GetUserByUserID(userID)
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	// 添加成员
	err = models.AddWorkerMember(req.WorkerUID, uint64(targetUser.ID), targetUser.UserName, uint64(currentUser.ID), currentUser.UserName)
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	common.RespOK(c, "add member success", nil)
}

// RemoveMemberEndpoint 移除协作者
func RemoveMemberEndpoint(c *gin.Context) {

	var req RemoveMemberRequest
	if err := c.BindJSON(&req); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}

	if req.WorkerUID == "" || req.UserID == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "worker_uid and user_name are required", nil)
		return
	}

	userID := c.GetUint(common.UIDKey)

	// 验证操作者是否可以管理成员
	_, err := permissions.CanManageWorkerMembers(c, uint64(userID), req.WorkerUID)
	if err != nil {
		return
	}

	// 获取要移除的用户 用户可能不存在
	// targetUser, err := models.GetUserByUserName(req.UserName)
	// if err != nil {
	// 	common.RespErr(c, common.RespCodeNotFound, "user not found", nil)
	// 	return
	// }

	// 移除成员
	err = models.RemoveWorkerMember(req.WorkerUID, uint64(req.UserID))
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	common.RespOK(c, "remove member success", nil)
}

// ListMembersEndpoint 列出协作者
func ListMembersEndpoint(c *gin.Context) {

	workerUID := c.Param("worker_uid")
	if workerUID == "" {
		common.RespErr(c, common.RespCodeInvalidRequest, "worker_uid is required", nil)
		return
	}

	userID := c.GetUint(common.UIDKey)

	// 验证操作者是否有权限
	_, err := permissions.CanReadWorker(c, uint64(userID), workerUID)
	if err != nil {
		return
	}

	// 获取成员列表
	members, err := models.GetWorkerMembers(workerUID)
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	common.RespOK(c, "list members success", members)
}

// GetWorkerCollaboratorsEndpoint 获取当前用户是拥有者或协作者的 worker 信息
func GetWorkerCollaboratorsEndpoint(c *gin.Context) {

	workerUID := c.Param("uid")
	if workerUID == "" {
		common.RespErr(c, common.RespCodeInvalidRequest, "worker_uid is required", nil)
		return
	}

	userID := c.GetUint(common.UIDKey)

	type CollaboratorInfo struct {
		IsOwner   bool                   `json:"is_owner"`
		CanManage bool                   `json:"can_manage"`
		Members   []*models.WorkerMember `json:"members"`
	}

	// 检查是否是拥有者
	ownerID, err := models.GetWorkerOwner(workerUID)
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	isOwner := ownerID == uint64(userID)
	canManage := isOwner

	// 如果不是拥有者，检查是否是协作者
	if !isOwner && !models.IsWorkerMember(workerUID, uint64(userID)) {
		common.RespErr(c, common.RespCodeNotAuthed, "forbidden", nil)
		return
	}

	// 获取 worker 详情
	worker, err := models.GetWorkerByUID(userID, workerUID)
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
		return
	}

	worker.Worker.Code = nil

	// 获取成员列表
	var members []*models.WorkerMember
	if isOwner {
		members, err = models.GetWorkerMembers(workerUID)
		if err != nil {
			common.RespErr(c, common.RespCodeInternalError, err.Error(), nil)
			return
		}
	}

	info := &CollaboratorInfo{
		IsOwner:   isOwner,
		CanManage: canManage,
		Members:   members,
	}

	common.RespOK(c, "get collaborator info success", info)
}

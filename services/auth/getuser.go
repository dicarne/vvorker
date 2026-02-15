package auth

import (
	"vvorker/common"
	"vvorker/conf"
	"vvorker/entities"
	"vvorker/models"

	"github.com/gin-gonic/gin"
)

func GetUserEndpoint(c *gin.Context) {

	uid, ok := common.RequireUID32(c)
	if !ok {
		return
	}
	user, err := models.GetUserByUserID(uid)
	if err != nil {
		common.RespErr(c, common.RespCodeDBErr, common.RespMsgDBErr, nil)
		return
	}
	common.RespOK(c, "ok", &entities.GetUserResponse{
		UserName: user.UserName,
		Role:     user.Role,
		Email:    user.Email,
		ID:       user.ID,
		VK:       conf.AppConfigInstance.EncryptionKey,
	})
}

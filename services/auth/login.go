package auth

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"vvorker/authz"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/entities"
	"vvorker/ext/kv/src/sys_cache"
	"vvorker/models"
	"vvorker/utils"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"github.com/sirupsen/logrus"
)

type tryCount struct {
	Count int `json:"count"`
}

func LoginEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	req, err := parseLoginReq(c)
	if err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest,
			common.RespMsgInvalidRequest, nil)
		return
	}

	try_count := &tryCount{}
	try_count_bin, err := sys_cache.Get("login_try_count:" + req.UserName)
	if err != nil {
		try_count = &tryCount{Count: 0}
	} else {
		err = json.Unmarshal([]byte(try_count_bin), try_count)
		if err != nil {
			try_count = &tryCount{Count: 0}
		}
	}

	if try_count.Count >= 5 {
		common.RespErr(c, common.RespCodeAuthErr,
			common.RespMsgAuthBan, nil)
		return
	}

	ok, err := models.CheckUserPassword(req.UserName, req.Password)
	if err != nil || !ok {
		try_count.Count++
		try_count_bin, _ := json.Marshal(try_count)
		sys_cache.Put("login_try_count:"+req.UserName, try_count_bin, 120)

		common.RespErr(c, common.RespCodeAuthErr,
			common.RespMsgAuthErr, nil)
		return
	}

	user, err := models.GetUserByUserName(req.UserName)
	if err != nil {
		logrus.WithError(err).Error("get user by user name failed")
		common.RespErr(c, common.RespCodeInternalError,
			common.RespMsgInternalError, nil)
		return
	}

	// 检查用户是否启用了OTP
	if user.OtpSecret != "" && conf.AppConfigInstance.EnableLoginOPT {
		// 启用了OTP，需要验证OTP
		otpCode := req.OTPCode
		if otpCode == "" {
			// 返回需要OTP验证的状态
			c.JSON(http.StatusOK, gin.H{
				"code":    common.RespCodeOTPRequired,
				"message": "OTP_REQUIRED",
				"data": gin.H{
					"status":     common.RespCodeOTPRequired,
					"requireOTP": true,
				},
			})
			return
		}

		// 验证OTP
		valid := totp.Validate(otpCode, user.OtpSecret)
		if !valid {
			common.RespErr(c, common.RespCodeAuthErr, "Invalid OTP", nil)
			return
		}
	}

	token, err := utils.SignToken(user.ID)
	if err != nil {
		logrus.WithError(err).Error("sign token failed")
		common.RespErr(c, common.RespCodeInternalError,
			common.RespMsgInternalError, nil)
		return
	}

	authz.SetToken(c, token)

	c.Header(common.AuthorizationHeaderKey, token)
	common.RespOK(c, common.RespMsgOK, entities.LoginResponse{
		Status: common.RespCodeOK,
		Token:  token})
}

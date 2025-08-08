package vvotp

import (
	"fmt"
	"vvorker/common"
	"vvorker/ext/kv/src/sys_cache"
	"vvorker/models"
	"vvorker/utils"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

func ValidOtpEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		}
	}()

	code := c.Copy().Query("code")
	if code == "" {
		common.RespErr(c, 403, "error", gin.H{"error": "code is empty"})
		return
	}

	db := database.GetDB()
	userID := c.GetUint(common.UIDKey)
	var user models.User
	if err := db.Where(&models.User{
		Model: gorm.Model{ID: userID},
	}).First(&user).Error; err != nil {
		common.RespErr(c, 403, "error", gin.H{"error": "User not found"})
		return
	}
	valid := totp.Validate(code, user.OtpSecret)
	if !valid {
		common.RespErr(c, 403, "error", gin.H{"error": "Invalid OTP"})
		return
	}

	token := utils.GenerateUID()
	_, err := sys_cache.Put("otp"+":"+"validtoken:"+fmt.Sprintf("%d", userID), []byte(token), 360)
	if err != nil {
		common.RespErr(c, 403, "error", gin.H{"error": "Failed to store OTP token"})
		return
	}
	common.RespOK(c, "OTP valid", gin.H{
		"vv-otp-token": token,
	})
}

func OTPMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		db := database.GetDB()
		var user models.User
		if err := db.Where(&models.User{
			Model: gorm.Model{ID: c.GetUint(common.UIDKey)},
		}).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(403, gin.H{"error": "OTP_REQUIRED", "message": "Unauthorized1"})
			return
		}
		if user.OtpSecret == "" {
			c.Next()
			return
		}
		token := c.Request.Header.Get("vv-otp-token")
		if token == "" {
			c.AbortWithStatusJSON(403, gin.H{"error": "OTP_REQUIRED", "message": "Unauthorized2"})
			return
		}

		if tk, err := sys_cache.Get("otp" + ":" + "validtoken:" + fmt.Sprintf("%d", c.GetUint(common.UIDKey))); err != nil || string(tk) != token {
			c.AbortWithStatusJSON(403, gin.H{"error": "OTP_REQUIRED", "message": "Unauthorized3"})
			return
		}

		c.Next()
	}
}

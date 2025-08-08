package vvotp

import (
	"vvorker/common"
	"vvorker/conf"
	"vvorker/ext/kv/src/sys_cache"
	"vvorker/models"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

func EnableOTPEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		}
	}()
	db := database.GetDB()
	userID := c.GetUint(common.UIDKey)
	var user models.User
	if err := db.Where(&models.User{
		Model: gorm.Model{ID: userID},
	}).First(&user).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "User not found"})
		return
	}

	if user.OtpSecret != "" {
		common.RespErr(c, 400, "error", gin.H{"error": "User already enabled OTP"})
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "otp" + conf.AppConfigInstance.WorkerURLSuffix,
		AccountName: user.UserName,
	})
	if err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "Failed to generate OTP key"})
		return
	}

	sys_cache.Put("otp"+":"+"tmpkey:"+user.UserName, []byte(key.Secret()), 360)

	common.RespOK(c, "OTP enabled successfully", gin.H{"url": key.String()})
}

func IsEnableOTPEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		}
	}()
	db := database.GetDB()
	userID := c.GetUint(common.UIDKey)
	var user models.User
	if err := db.Where(&models.User{
		Model: gorm.Model{ID: userID},
	}).First(&user).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "User not found"})
		return
	}

	common.RespOK(c, "OTP enabled successfully", gin.H{"enabled": user.OtpSecret != ""})
}

func ValidAddOTPEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		}
	}()
	db := database.GetDB()
	userID := c.GetUint(common.UIDKey)
	var user models.User
	if err := db.Where(&models.User{
		Model: gorm.Model{ID: userID},
	}).First(&user).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "User not found"})
		return
	}

	key, err := sys_cache.Get("otp" + ":" + "tmpkey:" + user.UserName)
	if err != nil || len(key) == 0 {
		common.RespErr(c, 400, "error", gin.H{"error": "Timeout"})
		return
	}
	// 验证
	code := c.Copy().Query("code")
	if code == "" {
		common.RespErr(c, 400, "error", gin.H{"error": "code is empty"})
		return
	}
	valid := totp.Validate(code, string(key))
	if !valid {
		common.RespErr(c, 400, "error", gin.H{"error": "Invalid OTP"})
		return
	}

	// update user otp
	if err := db.Model(&models.User{
		Model: gorm.Model{ID: userID},
	}).Updates(models.User{
		OtpSecret: string(key),
	}).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "Failed to update OTP key"})
		return
	}

	sys_cache.Del("otp" + ":" + "tmpkey:" + user.UserName)

	common.RespOK(c, "OTP enabled successfully", nil)
}

func DisableOTPEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			common.RespErr(c, 500, "error", gin.H{"error": "Internal server error"})
		}
	}()
	db := database.GetDB()
	userID := c.GetUint(common.UIDKey)
	var user models.User
	if err := db.Where(&models.User{
		Model: gorm.Model{ID: userID},
	}).First(&user).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "User not found"})
		return
	}
	user.OtpSecret = ""
	if err := db.Model(user).Save(user).Error; err != nil {
		common.RespErr(c, 500, "error", gin.H{"error": "Failed to update OTP key"})
		return
	}

	common.RespOK(c, "OTP disable successfully", nil)
}

package pgsql

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/entities"
	"vvorker/models"
	"vvorker/utils"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func CreateNewPostgreSQLResourcesEndpoint(c *gin.Context) {
	var req = entities.CreateNewResourcesRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusBadRequest, "invalid request", gin.H{"error": err.Error()})
		return
	}
	if !req.Validate() {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusBadRequest, "invalid request", gin.H{"error": "invalid request"})
		return
	}
	db := database.GetDB()
	userID, err := strconv.ParseUint(req.UserID, 10, 64)
	if err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusBadRequest, "Failed to convert UserID to uint64", gin.H{"error": err.Error()})
		return
	}
	pgResource := &models.PostgreSQL{
		UserID: userID,
		Name:   req.Name,
		UID:    utils.GenerateUID(),
	}
	pgResource.Database = "vvorker_" + pgResource.UID

	pgdb, err := sql.Open("postgres",
		"user="+conf.AppConfigInstance.ServerPostgreUser+
			" password="+conf.AppConfigInstance.ServerPostgrePassword+
			" host="+conf.AppConfigInstance.ServerPostgreHost+
			" port="+fmt.Sprintf("%d", conf.AppConfigInstance.ServerPostgrePort)+
			" sslmode=disable")
	if err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusInternalServerError, "Failed to connect to PostgreSQL", gin.H{"error": err.Error()})
		return
	}
	defer pgdb.Close()

	_, err = pgdb.Exec("CREATE DATABASE " + pgResource.Database)
	if err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusInternalServerError, "Failed to create database", gin.H{"error": err.Error()})
		return
	}

	// 生成随机密码
	password := utils.GenerateUID() // 假设 utils 包有 GenerateRandomString 函数
	pgUser := "vorker_user_" + pgResource.UID

	// 创建新用户
	_, err = pgdb.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", pgUser, password))
	if err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusInternalServerError, "Failed to create PostgreSQL user", gin.H{"error": err.Error()})
		return
	}

	// 授予用户对数据库的连接权限
	_, err = pgdb.Exec(fmt.Sprintf("GRANT CONNECT ON DATABASE %s TO %s", pgResource.Database, pgUser))
	if err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusInternalServerError, "Failed to grant CONNECT permission", gin.H{"error": err.Error()})
		return
	}

	// 切换到新创建的数据库
	_, err = pgdb.Exec(fmt.Sprintf("REVOKE ALL ON DATABASE %s FROM public", pgResource.Database))
	if err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusInternalServerError, "Failed to revoke public access", gin.H{"error": err.Error()})
		return
	}

	// 授予用户对数据库的所有表的增删改查权限
	// Connect to the newly created database
	targetConnStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		conf.AppConfigInstance.ServerPostgreUser,
		conf.AppConfigInstance.ServerPostgrePassword,
		conf.AppConfigInstance.ServerPostgreHost,
		conf.AppConfigInstance.ServerPostgrePort,
		pgResource.Database,
	)
	targetPgdb, err := sql.Open("postgres", targetConnStr)
	if err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusInternalServerError, "Failed to connect to the new database", gin.H{"error": err.Error()})
		return
	}
	defer targetPgdb.Close()

	// 授予用户对现有表和序列的所有权限
	_, err = targetPgdb.Exec(fmt.Sprintf(`
	    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO %s;
	    GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO %s;
	    GRANT USAGE ON SCHEMA public TO %s;
	`, pgUser, pgUser, pgUser))
	if err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusInternalServerError, "Failed to grant privileges", gin.H{"error": err.Error()})
		return
	}

	// 设置默认权限，确保未来创建的对象也被授予权限
	_, err = targetPgdb.Exec(fmt.Sprintf(`
	    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO %s;
	    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON SEQUENCES TO %s;
	    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE ON SEQUENCES TO %s;
	`, pgUser, pgUser, pgUser))
	if err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusInternalServerError, "Failed to set default privileges", gin.H{"error": err.Error()})
		return
	}

	// 授予用户创建表、视图、函数等权限
	_, err = targetPgdb.Exec(fmt.Sprintf(`
	    GRANT CREATE ON SCHEMA public TO %s;
	`, pgUser))
	if err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusInternalServerError, "Failed to grant CREATE privilege on schema", gin.H{"error": err.Error()})
		return
	}

	// 保存用户信息到 pgResource
	pgResource.Username = pgUser
	pgResource.Password = password

	if err := db.Create(pgResource).Error; err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusInternalServerError, "Failed to create PostgreSQL resource", gin.H{"error": err.Error()})
		return
	}

	// 使用 common.RespOK 返回成功响应
	common.RespOK(c, "success", gin.H{
		"uid":    pgResource.UID,
		"status": 0,
	})
}

func DeletePostgreSQLResourcesEndpoint(c *gin.Context) {
	uid := c.GetUint64(common.UIDKey)

	var req = entities.DeleteResourcesReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError,
			gin.H{"error": err.Error()})
		return
	}
	if !req.Validate() {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusBadRequest, "invalid request", gin.H{"error": "invalid request"})
		return
	}
	db := database.GetDB()
	// 检查 UID 是否为空
	if req.UID == "" {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusBadRequest, "UID cannot be empty", gin.H{"error": "UID cannot be empty"})
		return
	}
	// 存储查询条件
	condition := models.PostgreSQL{UID: req.UID, UserID: uid}

	// 执行删除操作并处理错误
	result := db.Delete(&condition, condition)
	if result.Error != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusInternalServerError, "Failed to delete PostgreSQL resource", gin.H{"error": result.Error.Error()})
		return
	}

	// 检查是否有记录被删除
	if result.RowsAffected == 0 {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusNotFound, "PostgreSQL resource not found", gin.H{"error": "PostgreSQL resource not found"})
		return
	}

	pgResourceDatabase := "vvorker_" + req.UID

	pgdb, err := sql.Open("postgres",
		"user="+conf.AppConfigInstance.ServerPostgreUser+
			" password="+conf.AppConfigInstance.ServerPostgrePassword+
			" host="+conf.AppConfigInstance.ServerPostgreHost+
			" port="+fmt.Sprintf("%d", conf.AppConfigInstance.ServerPostgrePort)+
			" sslmode=disable")
	if err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusInternalServerError, "Failed to connect to PostgreSQL", gin.H{"error": err.Error()})
		return
	}
	defer pgdb.Close()

	_, err = pgdb.Exec("DROP DATABASE " + pgResourceDatabase)
	if err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusInternalServerError, "Failed to drop database", gin.H{"error": err.Error()})
		return
	}
	// 使用 common.RespOK 返回成功响应
	common.RespOK(c, "success", gin.H{"status": 0})
}

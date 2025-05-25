package pgsql

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"vorker/conf"
	"vorker/entities"
	"vorker/models"
	"vorker/utils"
	"vorker/utils/database"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func CreateNewPostgreSQLResourcesEndpoint(c *gin.Context) {
	var req = entities.CreateNewResourcesRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !req.Validate() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	db := database.GetDB()
	userID, err := strconv.ParseUint(req.UserID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to convert UserID to uint64: " + err.Error()})
		return
	}
	pgResource := &models.PostgreSQL{
		UserID: userID,
		Name:   req.Name,
		UID:    utils.GenerateUID(),
	}
	pgResource.Database = "vorker_" + pgResource.UID

	pgdb, err := sql.Open("postgres",
		"user="+conf.AppConfigInstance.ServerPostgreUser+
			" password="+conf.AppConfigInstance.ServerPostgrePassword+
			" host="+conf.AppConfigInstance.ServerPostgreHost+
			" port="+fmt.Sprintf("%d", conf.AppConfigInstance.ServerPostgrePort)+
			" sslmode=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to PostgreSQL: " + err.Error()})
		return
	}
	defer pgdb.Close()

	_, err = pgdb.Exec("CREATE DATABASE " + pgResource.Database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create database: " + err.Error()})
		return
	}

	// 生成随机密码
	password := utils.GenerateUID() // 假设 utils 包有 GenerateRandomString 函数
	pgUser := "vorker_user_" + pgResource.UID

	// 创建新用户
	_, err = pgdb.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", pgUser, password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create PostgreSQL user: " + err.Error()})
		return
	}

	// 授予用户对数据库的连接权限
	_, err = pgdb.Exec(fmt.Sprintf("GRANT CONNECT ON DATABASE %s TO %s", pgResource.Database, pgUser))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to grant CONNECT permission: " + err.Error()})
		return
	}

	// 切换到新创建的数据库
	_, err = pgdb.Exec(fmt.Sprintf("REVOKE ALL ON DATABASE %s FROM public", pgResource.Database))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke public access: " + err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the new database: " + err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to grant privileges: " + err.Error()})
		return
	}

	// 设置默认权限，确保未来创建的对象也被授予权限
	_, err = targetPgdb.Exec(fmt.Sprintf(`
	    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO %s;
	    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON SEQUENCES TO %s;
	    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE ON SEQUENCES TO %s;
	`, pgUser, pgUser, pgUser))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set default privileges: " + err.Error()})
		return
	}

	// 授予用户创建表、视图、函数等权限
	_, err = targetPgdb.Exec(fmt.Sprintf(`
	    GRANT CREATE ON SCHEMA public TO %s;
	`, pgUser))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to grant CREATE privilege on schema: " + err.Error()})
		return
	}

	// 保存用户信息到 pgResource
	pgResource.Username = pgUser
	pgResource.Password = password

	if err := db.Create(pgResource).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create PostgreSQL resource: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uid":    pgResource.UID,
		"status": 0,
	})
}

func DeletePostgreSQLResourcesEndpoint(c *gin.Context) {
	var req = entities.DeleteResourcesReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !req.Validate() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	db := database.GetDB()
	// 检查 UID 是否为空
	if req.UID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UID cannot be empty"})
		return
	}

	// 存储查询条件
	condition := models.PostgreSQL{UID: req.UID}

	// 执行删除操作并处理错误
	result := db.Model(&models.PostgreSQL{}).Where(condition).Delete(&condition)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete PostgreSQL resource: " + result.Error.Error()})
		return
	}

	// 检查是否有记录被删除
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "PostgreSQL resource not found"})
		return
	}

	pgResourceDatabase := "vorker_" + req.UID

	pgdb, err := sql.Open("postgres",
		"user="+conf.AppConfigInstance.ServerPostgreUser+
			" password="+conf.AppConfigInstance.ServerPostgrePassword+
			" host="+conf.AppConfigInstance.ServerPostgreHost+
			" port="+fmt.Sprintf("%d", conf.AppConfigInstance.ServerPostgrePort)+
			" sslmode=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to PostgreSQL: " + err.Error()})
		return
	}
	defer pgdb.Close()

	_, err = pgdb.Exec("DROP DATABASE " + pgResourceDatabase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to drop database: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 0})
}

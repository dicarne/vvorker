package pgsql

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/entities"
	"vvorker/funcs"
	"vvorker/models"
	"vvorker/utils"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// CreatePostgreSQLDatabase 创建 PostgreSQL 数据库及相关用户，并授予权限
func CreatePostgreSQLDatabase(userID uint64, UID string, req entities.CreateNewResourcesRequest) (*models.PostgreSQL, error) {
	pgResource := &models.PostgreSQL{
		UserID: userID,
		Name:   req.Name,
		UID:    UID,
	}
	pgResource.Database = "vvorker_" + pgResource.UID

	pgdb, err := sql.Open("postgres",
		"user="+conf.AppConfigInstance.ServerPostgreUser+
			" password="+conf.AppConfigInstance.ServerPostgrePassword+
			" host="+conf.AppConfigInstance.ServerPostgreHost+
			" port="+fmt.Sprintf("%d", conf.AppConfigInstance.ServerPostgrePort)+
			" sslmode=disable")
	if err != nil {
		return nil, err
	}
	defer pgdb.Close()

	_, err = pgdb.Exec("CREATE DATABASE " + pgResource.Database)
	if err != nil {
		return nil, err
	}

	// 生成随机密码
	password := utils.GenerateUID() // 假设 utils 包有 GenerateRandomString 函数
	pgUser := "vorker_user_" + pgResource.UID

	// 创建新用户
	_, err = pgdb.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", pgUser, password))
	if err != nil {
		return nil, err
	}

	// 授予用户对数据库的连接权限
	_, err = pgdb.Exec(fmt.Sprintf("GRANT CONNECT ON DATABASE %s TO %s", pgResource.Database, pgUser))
	if err != nil {
		return nil, err
	}

	// 切换到新创建的数据库
	_, err = pgdb.Exec(fmt.Sprintf("REVOKE ALL ON DATABASE %s FROM public", pgResource.Database))
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer targetPgdb.Close()

	// 授予用户对现有表和序列的所有权限
	_, err = targetPgdb.Exec(fmt.Sprintf(`
	    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO %s;
	    GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO %s;
	    GRANT USAGE ON SCHEMA public TO %s;
	`, pgUser, pgUser, pgUser))
	if err != nil {
		return nil, err
	}

	// 设置默认权限，确保未来创建的对象也被授予权限
	_, err = targetPgdb.Exec(fmt.Sprintf(`
	    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO %s;
	    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON SEQUENCES TO %s;
	    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE ON SEQUENCES TO %s;
	`, pgUser, pgUser, pgUser))
	if err != nil {
		return nil, err
	}

	// 授予用户创建表、视图、函数等权限
	_, err = targetPgdb.Exec(fmt.Sprintf(`
	    GRANT CREATE ON SCHEMA public TO %s;
	`, pgUser))
	if err != nil {
		return nil, err
	}

	// 保存用户信息到 pgResource
	pgResource.Username = pgUser
	pgResource.Password = password

	db := database.GetDB()
	if err := db.Create(pgResource).Error; err != nil {
		// 使用 common.RespErr 返回错误响应
		return nil, err
	}

	return pgResource, nil
}

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

	userID := uint64(c.GetUint(common.UIDKey))
	if userID == 0 {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusBadRequest, "Failed to convert UserID to uint64", gin.H{"error": "uid is required"})
		return
	}
	UID := utils.GenerateUID()

	// 调用提取的函数创建数据库
	pgResource, err := CreatePostgreSQLDatabase(userID, UID, req)
	if err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusInternalServerError, "Failed to create PostgreSQL resource", gin.H{"error": err.Error()})
		return
	}

	// 使用 common.RespOK 返回成功响应
	common.RespOK(c, "success", entities.CreateNewResourcesResponse{
		UID:  pgResource.UID,
		Name: pgResource.Name,
		Type: "pgsql",
	})
}

func DeletePostgreSQLResourcesEndpoint(c *gin.Context) {
	uid := uint64(c.GetUint(common.UIDKey))

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
		common.RespOK(c, "success but not drop db because of unconnected db", entities.DeleteResourcesResp{
			Status: 0,
		})
		return
	}
	defer pgdb.Close()

	_, err = pgdb.Exec("DROP DATABASE " + pgResourceDatabase)
	if err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespOK(c, "success but not drop db", entities.DeleteResourcesResp{
			Status: 0,
		})
		return
	}
	// 使用 common.RespOK 返回成功响应
	common.RespOK(c, "success", entities.DeleteResourcesResp{
		Status: 0,
	})
}

func RecoverPGSQL(userID uint64, pgResource *models.PostgreSQL) (*models.PostgreSQL, error) {
	pgResource.UserID = userID
	db := database.GetDB()
	// 如果有，则更新，如果无，则调用新增接口
	if err := db.Where(&models.PostgreSQL{
		UID: pgResource.UID,
	}).First(&models.PostgreSQL{}).Error; err != nil {
		pg, err := CreatePostgreSQLDatabase(userID, pgResource.UID, entities.CreateNewResourcesRequest{Name: pgResource.Name})
		return pg, err
	} else {
		pgResource.Password = ""
		pgResource.Username = ""
		pgResource.Database = ""
		db.Where(&models.PostgreSQL{
			UID: pgResource.UID,
		}).Updates(pgResource)
	}

	return pgResource, nil
}

type updateFile struct {
	FileName string `json:"name"`
	Content  string `json:"content"`
}

type updateMigrateReq struct {
	ResourceID       string       `json:"resource_id"`
	Files            []updateFile `json:"files"`
	CustomDBName     string       `json:"custom_db_name"`
	CustomDBUser     string       `json:"custom_db_user"`
	CustomDBHost     string       `json:"custom_db_host"`
	CustomDBPort     int          `json:"custom_db_port"`
	CustomDBPassword string       `json:"custom_db_password"`
}

func UpdateMigrate(c *gin.Context) {
	var req updateMigrateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, http.StatusBadRequest, "invalid request", gin.H{"error": err.Error()})
		return
	}
	userID := uint64(c.GetUint(common.UIDKey))
	if userID == 0 {
		common.RespErr(c, http.StatusBadRequest, "Failed to convert UserID to uint64", gin.H{"error": "uid is required"})
		return
	}
	db := database.GetDB()
	UID := ""

	if !strings.HasPrefix(req.ResourceID, "worker_resource:pgsql:") {
		pgResource := models.PostgreSQL{}
		if err := db.Where(&models.PostgreSQL{
			UID: req.ResourceID,
		}).First(&pgResource).Error; err != nil {
			common.RespErr(c, http.StatusInternalServerError, "Failed to get PostgreSQL resource", gin.H{"error": err.Error()})
			return
		}
		if pgResource.UserID != userID {
			common.RespErr(c, http.StatusUnauthorized, "Unauthorized", gin.H{"error": "Unauthorized"})
			return
		}
		UID = pgResource.UID
	} else {
		UID = req.ResourceID
	}

	if err := db.Where(&models.PostgreSQLMigration{
		UserID: userID,
		DBUID:  UID,
	}).Delete(&models.PostgreSQLMigration{}).Error; err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to delete PostgreSQL migration", gin.H{"error": err.Error()})
		return
	}

	for i, file := range req.Files {
		if err := db.Model(&models.PostgreSQLMigration{}).Create(&models.PostgreSQLMigration{
			UserID:           userID,
			DBUID:            UID,
			FileName:         file.FileName,
			FileContent:      file.Content,
			Sequence:         i,
			CustomDBName:     req.CustomDBName,
			CustomDBUser:     req.CustomDBUser,
			CustomDBPassword: req.CustomDBPassword,
			CustomDBHost:     req.CustomDBHost,
			CustomDBPort:     req.CustomDBPort,
		}).Error; err != nil {
			common.RespErr(c, http.StatusInternalServerError, "Failed to create PostgreSQL migration", gin.H{"error": err.Error()})
			return
		}
	}

	common.RespOK(c, "success", gin.H{})
}

func migrateResource(userID uint64, pgid string) error {
	db := database.GetDB()
	pgResource := models.PostgreSQL{}
	if err := db.Where(&models.PostgreSQL{
		UID: pgid,
	}).First(&pgResource).Error; err != nil {
		return err
	}

	migrates := []models.PostgreSQLMigration{}
	if err := db.Where(&models.PostgreSQLMigration{
		UserID: userID,
		DBUID:  pgid,
	}).Order("sequence").Find(&migrates).Error; err != nil {
		return err
	}

	pgResource.Database = "vvorker_" + pgResource.UID

	pgdb, err := sql.Open("postgres",
		"user="+conf.AppConfigInstance.ServerPostgreUser+
			" password="+conf.AppConfigInstance.ServerPostgrePassword+
			" host="+conf.AppConfigInstance.ServerPostgreHost+
			" port="+fmt.Sprintf("%d", conf.AppConfigInstance.ServerPostgrePort)+
			" sslmode=disable"+
			" dbname="+pgResource.Database)
	if err != nil {
		return err
	}
	defer pgdb.Close()

	for _, migrate := range migrates {
		_, err = pgdb.Exec(migrate.FileContent)
		if err != nil {
			logrus.Error(err)
			// return err
		}
	}
	return nil
}

func migrateCustomResource(userID uint64, pgid string) error {
	db := database.GetDB()
	migrates := []models.PostgreSQLMigration{}
	if err := db.Where(&models.PostgreSQLMigration{
		UserID: userID,
		DBUID:  pgid,
	}).Order("sequence").Find(&migrates).Error; err != nil {
		return err
	}
	if len(migrates) == 0 {
		return nil
	}
	config := migrates[0]
	pgdb, err := sql.Open("postgres",
		"user="+config.CustomDBUser+
			" password="+config.CustomDBPassword+
			" host="+config.CustomDBHost+
			" port="+fmt.Sprintf("%d", config.CustomDBPort)+
			" sslmode=disable")
	if err != nil {
		return err
	}
	defer pgdb.Close()
	for _, migrate := range migrates {
		_, err = pgdb.Exec(migrate.FileContent)
		if err != nil {
			logrus.Error(err)
			// return err
		}
	}
	return nil
}

func MigratePostgreSQLDatabase(userID uint64, pgid string) error {
	if !strings.HasPrefix(pgid, "worker_resource:pgsql:") {
		return migrateResource(userID, pgid)
	} else {
		return migrateCustomResource(userID, pgid)
	}
}

func init() {
	funcs.SetMigratePostgreSQLDatabase(MigratePostgreSQLDatabase)
}

package extmysql

import (
	"database/sql"
	"fmt"
	"net/http"
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

// CreateMySQLDatabase 创建 MySQL 数据库及相关用户，并授予权限
func CreateMySQLDatabase(userID uint64, UID string, req entities.CreateNewResourcesRequest) (*models.MySQL, error) {
	mysqlResource := &models.MySQL{
		UserID: userID,
		Name:   req.Name,
		UID:    UID,
	}
	mysqlResource.Database = "vvorker_" + mysqlResource.UID

	pgdb, err := sql.Open("mysql",
		"user="+conf.AppConfigInstance.ServerMySQLUser+
			" password="+conf.AppConfigInstance.ServerMySQLPassword+
			" host="+conf.AppConfigInstance.ServerMySQLHost+
			" port="+fmt.Sprintf("%d", conf.AppConfigInstance.ServerMySQLPort)+
			" sslmode=disable")
	if err != nil {
		return nil, err
	}
	defer pgdb.Close()

	_, err = pgdb.Exec("CREATE DATABASE " + mysqlResource.Database)
	if err != nil {
		return nil, err
	}

	// 生成随机密码
	password := utils.GenerateUID() // 假设 utils 包有 GenerateRandomString 函数
	pgUser := "vorker_user_" + mysqlResource.UID

	// 创建新用户
	_, err = pgdb.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", pgUser, password))
	if err != nil {
		return nil, err
	}

	// 授予用户对数据库的连接权限
	_, err = pgdb.Exec(fmt.Sprintf("GRANT CONNECT ON DATABASE %s TO %s", mysqlResource.Database, pgUser))
	if err != nil {
		return nil, err
	}

	// 切换到新创建的数据库
	_, err = pgdb.Exec(fmt.Sprintf("REVOKE ALL ON DATABASE %s FROM public", mysqlResource.Database))
	if err != nil {
		return nil, err
	}

	// 授予用户对数据库的所有表的增删改查权限
	// Connect to the newly created database
	targetConnStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		conf.AppConfigInstance.ServerMySQLUser,
		conf.AppConfigInstance.ServerMySQLPassword,
		conf.AppConfigInstance.ServerMySQLHost,
		conf.AppConfigInstance.ServerMySQLPort,
		mysqlResource.Database,
	)
	targetPgdb, err := sql.Open("mysql", targetConnStr)
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

	// 保存用户信息到 mysqlResource
	mysqlResource.Username = pgUser
	mysqlResource.Password = password

	db := database.GetDB()
	if err := db.Create(mysqlResource).Error; err != nil {
		// 使用 common.RespErr 返回错误响应
		return nil, err
	}

	return mysqlResource, nil
}

func CreateNewMySQLResourcesEndpoint(c *gin.Context) {
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
	mysqlResource, err := CreateMySQLDatabase(userID, UID, req)
	if err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusInternalServerError, "Failed to create MySQL resource", gin.H{"error": err.Error()})
		return
	}

	// 使用 common.RespOK 返回成功响应
	common.RespOK(c, "success", entities.CreateNewResourcesResponse{
		UID:  mysqlResource.UID,
		Name: mysqlResource.Name,
		Type: "pgsql",
	})
}

func DeleteMySQLResourcesEndpoint(c *gin.Context) {
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
	condition := models.MySQL{UID: req.UID, UserID: uid}

	// 执行删除操作并处理错误
	result := db.Delete(&condition, condition)
	if result.Error != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusInternalServerError, "Failed to delete MySQL resource", gin.H{"error": result.Error.Error()})
		return
	}

	// 检查是否有记录被删除
	if result.RowsAffected == 0 {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusNotFound, "MySQL resource not found", gin.H{"error": "MySQL resource not found"})
		return
	}

	mysqlResourceDatabase := "vvorker_" + req.UID

	pgdb, err := sql.Open("mysql",
		"user="+conf.AppConfigInstance.ServerMySQLUser+
			" password="+conf.AppConfigInstance.ServerMySQLPassword+
			" host="+conf.AppConfigInstance.ServerMySQLHost+
			" port="+fmt.Sprintf("%d", conf.AppConfigInstance.ServerMySQLPort)+
			" sslmode=disable")
	if err != nil {
		// 使用 common.RespErr 返回错误响应
		common.RespOK(c, "success but not drop db because of unconnected db", entities.DeleteResourcesResp{
			Status: 0,
		})
		return
	}
	defer pgdb.Close()

	_, err = pgdb.Exec("DROP DATABASE " + mysqlResourceDatabase)
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

func RecoverPGSQL(userID uint64, mysqlResource *models.MySQL) (*models.MySQL, error) {
	mysqlResource.UserID = userID
	db := database.GetDB()
	// 如果有，则更新，如果无，则调用新增接口
	if err := db.Where("uid =?", mysqlResource.UID).First(&models.MySQL{}).Error; err != nil {
		pg, err := CreateMySQLDatabase(userID, mysqlResource.UID, entities.CreateNewResourcesRequest{Name: mysqlResource.Name})
		return pg, err
	} else {
		mysqlResource.Password = ""
		mysqlResource.Username = ""
		mysqlResource.Database = ""
		db.Where("uid =?", mysqlResource.UID).Updates(mysqlResource)
	}

	return mysqlResource, nil
}

type updateFile struct {
	FileName string `json:"name"`
	Content  string `json:"content"`
}

type updateMigrateReq struct {
	ResourceID string       `json:"resource_id"`
	Files      []updateFile `json:"files"`
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
	mysqlResource := models.MySQL{}
	if err := db.Where("uid =?", req.ResourceID).First(&mysqlResource).Error; err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to get MySQL resource", gin.H{"error": err.Error()})
		return
	}
	if mysqlResource.UserID != userID {
		common.RespErr(c, http.StatusUnauthorized, "Unauthorized", gin.H{"error": "Unauthorized"})
		return
	}

	if err := db.Where(&models.MySQLMigration{
		UserID: userID,
		DBUID:  mysqlResource.UID,
	}).Delete(&models.MySQLMigration{}).Error; err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to delete MySQL migration", gin.H{"error": err.Error()})
		return
	}

	for i, file := range req.Files {
		if err := db.Model(&models.MySQLMigration{}).Create(&models.MySQLMigration{
			UserID:      userID,
			DBUID:       mysqlResource.UID,
			FileName:    file.FileName,
			FileContent: file.Content,
			Sequence:    i,
		}).Error; err != nil {
			common.RespErr(c, http.StatusInternalServerError, "Failed to create MySQL migration", gin.H{"error": err.Error()})
			return
		}
	}

	common.RespOK(c, "success", gin.H{})
}

func MigrateMySQLDatabase(userID uint64, pgid string) error {
	db := database.GetDB()
	mysqlResource := models.MySQL{}
	if err := db.Where("uid =?", pgid).First(&mysqlResource).Error; err != nil {
		return err
	}

	migrates := []models.MySQLMigration{}
	if err := db.Where(&models.MySQLMigration{
		UserID: userID,
		DBUID:  mysqlResource.UID,
	}).Order("sequence").Find(&migrates).Error; err != nil {
		return err
	}

	mysqlResource.Database = "vvorker_" + mysqlResource.UID

	pgdb, err := sql.Open("mysql",
		"user="+conf.AppConfigInstance.ServerMySQLUser+
			" password="+conf.AppConfigInstance.ServerMySQLPassword+
			" host="+conf.AppConfigInstance.ServerMySQLHost+
			" port="+fmt.Sprintf("%d", conf.AppConfigInstance.ServerMySQLPort)+
			" sslmode=disable"+
			" dbname="+mysqlResource.Database)
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

func init() {
	funcs.SetMigrateMySQLDatabase(MigrateMySQLDatabase)
}

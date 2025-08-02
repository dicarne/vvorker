package extmysql

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/defs"
	"vvorker/entities"
	"vvorker/funcs"
	"vvorker/models"
	"vvorker/utils"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func buildMysqlConnectionString() string {
	// username:password@protocol(address)/dbname?param=value
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
		conf.AppConfigInstance.ServerMySQLUser,
		conf.AppConfigInstance.ServerMySQLPassword,
		conf.AppConfigInstance.ServerMySQLHost,
		conf.AppConfigInstance.ServerMySQLPort)
}

func buildMysqlDBConnectionString(database string) string {
	if conf.AppConfigInstance.ServerMySQLOneDBName != "" {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			conf.AppConfigInstance.ServerMySQLUser,
			conf.AppConfigInstance.ServerMySQLPassword,
			conf.AppConfigInstance.ServerMySQLHost,
			conf.AppConfigInstance.ServerMySQLPort,
			conf.AppConfigInstance.ServerMySQLOneDBName)
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.AppConfigInstance.ServerMySQLUser,
		conf.AppConfigInstance.ServerMySQLPassword,
		conf.AppConfigInstance.ServerMySQLHost,
		conf.AppConfigInstance.ServerMySQLPort,
		database)
}

func cutDatabaseName(database string) string {
	// database 不超过63个字符
	if len(database) > 63 {
		return database[:63]
	}
	return database
}

func cutUserName(user string) string {
	if len(user) > 32 {
		return user[:32]
	}
	return user
}

// CreateMySQLDatabase 创建 MySQL 数据库及相关用户，并授予权限
func CreateMySQLDatabase(userID uint64, UID string, req entities.CreateNewResourcesRequest) (*models.MySQL, error) {
	mysqlResource := &models.MySQL{
		UserID: userID,
		Name:   req.Name,
		UID:    UID,
	}
	if conf.AppConfigInstance.ServerMySQLOneDBName != "" {
		mysqlResource.Database = conf.AppConfigInstance.ServerMySQLOneDBName
	} else {
		mysqlResource.Database = cutDatabaseName("vvorker_" + mysqlResource.UID)

		// 用户:密码@/库名?charset=utf8&parseTime=True&loc=Local
		pgdb, err := sql.Open("mysql", buildMysqlConnectionString())
		if err != nil {
			return nil, err
		}
		defer pgdb.Close()

		// Create database with UTF8MB4 character set and case-insensitive collation
		_, err = pgdb.Exec("CREATE DATABASE `" + mysqlResource.Database + "` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
		if err != nil {
			return nil, err
		}

		// Generate random password
		password := utils.GenerateUID()
		pgUser := cutUserName("vorker_user_" + mysqlResource.UID)

		// Create new user with password
		_, err = pgdb.Exec(fmt.Sprintf("CREATE USER `%s`@'%%' IDENTIFIED BY '%s'", pgUser, password))
		if err != nil {
			return nil, err
		}

		// Grant all privileges on the database to the user
		_, err = pgdb.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON `%s`.* TO `%s`@'%%'", mysqlResource.Database, pgUser))
		if err != nil {
			return nil, err
		}

		// Flush privileges to apply changes
		_, err = pgdb.Exec("FLUSH PRIVILEGES")
		if err != nil {
			return nil, err
		}

		// Connect to the newly created database to set up any additional permissions
		targetConnStr := buildMysqlDBConnectionString(mysqlResource.Database)
		targetPgdb, err := sql.Open("mysql", targetConnStr)
		if err != nil {
			return nil, err
		}
		defer targetPgdb.Close()

		// In MySQL, the above GRANT statement already gives all necessary privileges
		// No need for additional GRANT statements like in PostgreSQL

		// For MySQL 8.0+, you might want to set the default role if using roles
		// But for most cases, the above GRANT is sufficient

		// 保存用户信息到 mysqlResource
		mysqlResource.Username = pgUser
		mysqlResource.Password = password
	}

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
		Type: "mysql",
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

	if conf.AppConfigInstance.ServerMySQLOneDBName == "" {
		mysqlResourceDatabase := cutDatabaseName("vvorker_" + req.UID)
		pgdb, err := sql.Open("mysql",
			buildMysqlDBConnectionString(mysqlResourceDatabase))
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
	if err := db.Where(&models.MySQL{
		UID: mysqlResource.UID,
	}).First(&models.MySQL{}).Error; err != nil {
		pg, err := CreateMySQLDatabase(userID, mysqlResource.UID, entities.CreateNewResourcesRequest{Name: mysqlResource.Name})
		return pg, err
	} else {
		mysqlResource.Password = ""
		mysqlResource.Username = ""
		mysqlResource.Database = ""
		db.Where(&models.MySQL{
			UID: mysqlResource.UID,
		}).Updates(mysqlResource)
	}

	return mysqlResource, nil
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

	if !strings.HasPrefix(req.ResourceID, "worker_resource:mysql:") {
		mysqlResource := models.MySQL{}
		if err := db.Where(&models.MySQL{
			UID: req.ResourceID,
		}).First(&mysqlResource).Error; err != nil {
			common.RespErr(c, http.StatusInternalServerError, "Failed to get MySQL resource", gin.H{"error": err.Error()})
			return
		}
		if mysqlResource.UserID != userID {
			common.RespErr(c, http.StatusUnauthorized, "Unauthorized", gin.H{"error": "Unauthorized"})
			return
		}
		UID = mysqlResource.UID
	} else {
		UID = req.ResourceID
	}

	if err := db.Where(&models.MySQLMigration{
		UserID: userID,
		DBUID:  UID,
	}).Delete(&models.MySQLMigration{}).Error; err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to delete MySQL migration", gin.H{"error": err.Error()})
		return
	}

	for i, file := range req.Files {
		if err := db.Model(&models.MySQLMigration{}).Create(&models.MySQLMigration{
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
			common.RespErr(c, http.StatusInternalServerError, "Failed to create MySQL migration", gin.H{"error": err.Error()})
			return
		}
	}

	common.RespOK(c, "success", gin.H{})
}

func migrateCustomMySQLResource(userID uint64, pgid string) error {
	db := database.GetDB()
	migrates := []models.MySQLMigration{}
	if err := db.Where(&models.MySQLMigration{
		UserID: userID,
		DBUID:  pgid,
	}).Order("sequence").Find(&migrates).Error; err != nil {
		return err
	}

	if len(migrates) == 0 {
		return nil
	}

	config := migrates[0]
	dbConnectionStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		config.CustomDBUser,
		config.CustomDBPassword,
		config.CustomDBHost,
		config.CustomDBPort,
		config.CustomDBName,
	)
	logrus.Infof("dbConnectionStr: %s", dbConnectionStr)

	dbConn, err := sql.Open("mysql", dbConnectionStr)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	for _, migrate := range migrates {
		_, err = dbConn.Exec(migrate.FileContent)
		if err != nil {
			logrus.Error(err)
			// Continue with next migration even if one fails
		}
	}

	return nil
}

func MigrateMySQLDatabase(userID uint64, pgid string) error {
	if !strings.HasPrefix(pgid, "worker_resource:mysql:") {
		// Original migration logic for non-custom resources
		db := database.GetDB()
		mysqlResource := models.MySQL{}
		if err := db.Where(&models.MySQL{
			UID: pgid,
		}).First(&mysqlResource).Error; err != nil {
			return err
		}

		migrates := []models.MySQLMigration{}
		if err := db.Where(&models.MySQLMigration{
			UserID: userID,
			DBUID:  mysqlResource.UID,
		}).Order("sequence").Find(&migrates).Error; err != nil {
			return err
		}

		mysqlResource.Database = cutDatabaseName("vvorker_" + mysqlResource.UID)

		dbConn, err := sql.Open("mysql", buildMysqlDBConnectionString(mysqlResource.Database)+"&multiStatements=true")
		if err != nil {
			return err
		}
		defer dbConn.Close()

		for _, migrate := range migrates {
			_, err = dbConn.Exec(migrate.FileContent)
			if err != nil {
				logrus.Error(err)
				// Continue with next migration even if one fails
			}
		}

		return nil
	}

	// Handle custom database connection
	return migrateCustomMySQLResource(userID, pgid)
}

func init() {
	funcs.SetMigrateMySQLDatabase(MigrateMySQLDatabase)
	dbConns = defs.NewSyncMap[string, *sql.DB](map[string]*sql.DB{})
}

var dbConns *defs.SyncMap[string, *sql.DB]

func ExecuteSQLMysqlEndpoint(c *gin.Context) {
	var req = entities.ExecuteSQLReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError,
			gin.H{"error": err.Error()})
		return
	}
	dbConn, ok := dbConns.Get(req.ConnectionString)
	if !ok {
		dbConn, err := sql.Open("mysql", req.ConnectionString)
		if err != nil {
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError,
				gin.H{"error": err.Error()})
			return
		}
		defer dbConn.Close()
		dbConns.Set(req.ConnectionString, dbConn)
	}
	rows, err := dbConn.Query(req.Sql, req.Params...)
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError,
			gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	if req.Method == "get" {
		var rowsAll []string = []string{}
		for rows.Next() {
			var row string
			rows.Scan(&row)
			rowsAll = append(rowsAll, row)
		}
		c.JSON(200, entities.ExecuteSQLResp{Rows: rowsAll})
		return
	} else {
		var rowsAll [][]string = [][]string{}
		// First, get the column names
		columns, err := rows.Columns()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Create a slice of interface{} to hold the scanned values
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(sql.RawBytes)
		}

		for rows.Next() {
			// Scan the row into the values slice
			if err := rows.Scan(values...); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			// Convert each value to string
			row := make([]string, len(columns))
			for i, val := range values {
				if rb, ok := val.(*sql.RawBytes); ok {
					row[i] = string(*rb)
				}
			}
			rowsAll = append(rowsAll, row)
		}
		c.JSON(200, entities.ExecuteSQLRespAll{Rows: rowsAll})
		return
	}
}

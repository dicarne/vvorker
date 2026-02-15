package resource

import (
	"vvorker/common"
	"vvorker/models"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
)

type ListResourceRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	RType    string `json:"type"`
}

type ResourceData struct {
	UID      string   `json:"uid"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	ErrorMsg []string `json:"error_msg"`
}

type ListResourceResponse struct {
	Total int64          `json:"total"`
	Data  []ResourceData `json:"data"`
}

func ListResourceEndpoint(c *gin.Context) {
	uid, ok := common.RequireUID(c)
	if !ok {
		return
	}

	request := ListResourceRequest{}
	if err := c.BindJSON(&request); err != nil {
		return
	}
	db := database.GetDB()

	response := ListResourceResponse{
		Total: 0,
		Data:  make([]ResourceData, 0),
	}

	if request.RType == "kv" {
		var total int64
		var resources []models.KV
		db.Model(&models.KV{}).Where(&models.KV{
			UserID: uid,
		}).Limit(request.PageSize).Offset((request.Page - 1) * request.PageSize).Find(&resources)
		db.Model(&models.KV{}).Where(&models.KV{
			UserID: uid,
		}).Count(&total)

		for _, resource := range resources {
			response.Data = append(response.Data, ResourceData{
				UID:  resource.UID,
				Name: resource.Name,
				Type: "kv",
			})
		}
		response.Total = total
	} else if request.RType == "oss" {
		var total int64
		var resources []models.OSS
		db.Model(&models.OSS{}).Where(&models.OSS{
			UserID: uid,
		}).Limit(request.PageSize).Offset((request.Page - 1) * request.PageSize).Find(&resources)
		db.Model(&models.OSS{}).Where(&models.OSS{
			UserID: uid,
		}).Count(&total)
		for _, resource := range resources {
			response.Data = append(response.Data, ResourceData{
				UID:  resource.UID,
				Name: resource.Name,
				Type: "oss",
			})
		}
		response.Total = total
	} else if request.RType == "pgsql" {
		var total int64
		var resources []models.PostgreSQL
		db.Model(&models.PostgreSQL{}).Where(&models.PostgreSQL{
			UserID: uid,
		}).Limit(request.PageSize).Offset((request.Page - 1) * request.PageSize).Find(&resources)
		db.Model(&models.PostgreSQL{}).Where(&models.PostgreSQL{
			UserID: uid,
		}).Count(&total)
		for _, resource := range resources {
			response.Data = append(response.Data, ResourceData{
				UID:  resource.UID,
				Name: resource.Name,
				Type: "pgsql",
			})
		}
		for i, d := range resources {
			var migrations []models.PostgreSQLMigration
			db.Where(&models.PostgreSQLMigration{
				DBUID: d.UID,
			}).Find(&migrations)
			for _, migration := range migrations {
				if migration.MigrateKey != "" {
					var mg models.MigrationHistory
					db.Where(&models.MigrationHistory{
						Key: migration.MigrateKey,
					}).Find(&mg)
					if mg.Error != "" {
						response.Data[i].ErrorMsg = append(response.Data[i].ErrorMsg, mg.Error)
					}
				}
			}
		}
		response.Total = total
	} else if request.RType == "mysql" {
		var total int64
		var resources []models.MySQL
		db.Model(&models.MySQL{}).Where(&models.MySQL{
			UserID: uid,
		}).Limit(request.PageSize).Offset((request.Page - 1) * request.PageSize).Find(&resources)
		db.Model(&models.MySQL{}).Where(&models.MySQL{
			UserID: uid,
		}).Count(&total)
		for _, resource := range resources {
			response.Data = append(response.Data, ResourceData{
				UID:  resource.UID,
				Name: resource.Name,
				Type: "mysql",
			})
		}
		for i, d := range resources {
			var migrations []models.MySQLMigration
			db.Where(&models.MySQLMigration{
				DBUID: d.UID,
			}).Find(&migrations)
			for _, migration := range migrations {
				if migration.MigrateKey != "" {
					var mg models.MigrationHistory
					db.Where(&models.MigrationHistory{
						Key: migration.MigrateKey,
					}).Find(&mg)
					if mg.Error != "" {
						response.Data[i].ErrorMsg = append(response.Data[i].ErrorMsg, mg.Error)
					}
				}
			}
		}
		response.Total = total
	}
	common.RespOK(c, "success", response)
}

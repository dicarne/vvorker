package resource

import (
	"net/http"
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
	UID  string `json:"uid"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ListResourceResponse struct {
	Total int64          `json:"total"`
	Data  []ResourceData `json:"data"`
}

func ListResourceEndpoint(c *gin.Context) {
	uid := uint64(c.GetUint(common.UIDKey))

	request := ListResourceRequest{}
	if err := c.BindJSON(&request); err != nil {
		common.RespErr(c, http.StatusInternalServerError, "List resource failed.", gin.H{"error": err.Error()})
		return
	}
	if uid == 0 {
		common.RespErr(c, http.StatusInternalServerError, "uid is required.", gin.H{"error": "uid is required"})
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
		response.Total = total
	}
	common.RespOK(c, "success", response)
}

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
		c.JSON(400, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	if uid == 0 {
		c.JSON(400, gin.H{"code": 1, "msg": "uid is required"})
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
		db.Model(&models.KV{}).Where("user_id = ?", uid).
			Limit(request.PageSize).Offset((request.Page - 1) * request.PageSize).Find(&resources)
		db.Model(&models.KV{}).Where("user_id =?", uid).Count(&total)

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
		db.Model(&models.OSS{}).Where("user_id =?", uid).
			Limit(request.PageSize).Offset((request.Page - 1) * request.PageSize).Find(&resources)
		db.Model(&models.OSS{}).Where("user_id =?", uid).Count(&total)
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
		db.Model(&models.PostgreSQL{}).Where("user_id =?", uid).
			Limit(request.PageSize).Offset((request.Page - 1) * request.PageSize).Find(&resources)
		db.Model(&models.PostgreSQL{}).Where("user_id =?", uid).Count(&total)
		for _, resource := range resources {
			response.Data = append(response.Data, ResourceData{
				UID:  resource.UID,
				Name: resource.Name,
				Type: "pgsql",
			})
		}
		response.Total = total
	}
	common.RespOK(c, "success", response)
}

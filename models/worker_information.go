package models

import (
	"vvorker/entities"
	"vvorker/utils/database"

	"gorm.io/gorm"
)

type WorkerInformation struct {
	gorm.Model
	UID         string `json:"uid" gorm:"index"`
	Example     string `json:"example"`
	Testcases   string `json:"testcases"`
	Description string `json:"description"`
}

type WorkerDetailed struct {
	*entities.Worker
	Detail *WorkerInformation `json:"Detail"`
}

// CreateWorkerInformation 根据 UID 创建新的 WorkerInformation 记录
func CreateWorkerInformation(info *WorkerInformation) error {
	// 使用 database.GetDB() 获取数据库连接
	return database.GetDB().Create(info).Error
}

// GetWorkerInformationByUID 根据 UID 查询 WorkerInformation 记录
func GetWorkerInformationByUID(UID string) (*WorkerInformation, error) {
	var info WorkerInformation
	// 使用 database.GetDB() 获取数据库连接
	err := database.GetDB().Where("uid = ?", UID).First(&info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}

// UpdateWorkerInformationByUID 根据 UID 更新 WorkerInformation 记录
func UpdateWorkerInformationByUID(UID string, updateData map[string]interface{}) error {
	// 使用 database.GetDB() 获取数据库连接
	exsist := database.GetDB().Where("uid =?", UID).First(&WorkerInformation{}).Error
	if exsist != nil {
		return database.GetDB().Model(&WorkerInformation{}).Where("uid = ?", UID).Updates(updateData).Error
	}
	return database.GetDB().Create(&WorkerInformation{UID: UID}).Updates(updateData).Error
}

// DeleteWorkerInformationByUID 根据 UID 删除 WorkerInformation 记录
func DeleteWorkerInformationByUID(UID string) error {
	// 使用 database.GetDB() 获取数据库连接
	result := database.GetDB().Where("uid = ?", UID).Delete(&WorkerInformation{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

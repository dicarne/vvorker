package models

import (
	"vvorker/utils/database"

	"gorm.io/gorm"
)

type WorkerInformationBase struct {
	UID         string `json:"uid" gorm:"index"`
	Example     string `json:"example"`
	Description string `json:"description"`
}

type WorkerInformation struct {
	gorm.Model
	*WorkerInformationBase
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
	err := database.GetDB().Where(&WorkerInformation{
		WorkerInformationBase: &WorkerInformationBase{
			UID: UID,
		},
	}).First(&info).Error
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// UpdateWorkerInformationByUID 根据 UID 更新 WorkerInformation 记录
func UpdateWorkerInformationByUID(UID string, updateData *WorkerInformation) error {
	// 使用 database.GetDB() 获取数据库连接
	exsist := database.GetDB().Where(&WorkerInformation{
		WorkerInformationBase: &WorkerInformationBase{
			UID: UID,
		},
	}).First(&WorkerInformation{}).Error
	if exsist != nil {
		return database.GetDB().Create(
			&WorkerInformation{
				WorkerInformationBase: updateData.WorkerInformationBase,
			},
		).Error
	}
	return database.GetDB().Model(&WorkerInformation{}).Where(&WorkerInformation{
		WorkerInformationBase: &WorkerInformationBase{
			UID: UID,
		},
	}).Updates(updateData).Error
}

// DeleteWorkerInformationByUID 根据 UID 删除 WorkerInformation 记录
func DeleteWorkerInformationByUID(UID string) error {
	// 使用 database.GetDB() 获取数据库连接
	result := database.GetDB().Where(&WorkerInformation{
		WorkerInformationBase: &WorkerInformationBase{
			UID: UID,
		},
	}).Delete(&WorkerInformation{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

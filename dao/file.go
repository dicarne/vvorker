package dao

import (
	"context"
	"fmt"
	"vvorker/models"
	"vvorker/utils"
	"vvorker/utils/database"
)

func SaveFile(c context.Context, file *models.File) (*models.File, error) {
	if file == nil {
		return nil, fmt.Errorf("file is nil")
	}

	if file.UID == "" {
		file.UID = utils.GenerateUID()
	}

	db := database.GetDB()
	return file, db.Save(file).Error
}

func GetFileByHashAndCreator(c context.Context, hash string, creator uint) (*models.File, error) {
	db := database.GetDB()
	file := &models.File{}

	if err := db.Where(&models.File{
		Hash:      hash,
		CreatedBy: creator,
	}).First(file).Error; err != nil {
		return nil, err
	}

	return file, nil
}

func GetFileByUID(c context.Context, userID uint, fileID string) (*models.File, error) {
	file := &models.File{}
	if err := database.GetDB().Where(&models.File{
		UID:       fileID,
		CreatedBy: userID,
	}).First(file).Error; err != nil {
		return nil, err
	}
	return file, nil
}

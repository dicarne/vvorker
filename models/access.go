package models

import "gorm.io/gorm"

type AccessKey struct {
	gorm.Model
	UserId uint64 `json:"user_id"`
	Name   string `json:"name"`
	Key    string `json:"key"`
}

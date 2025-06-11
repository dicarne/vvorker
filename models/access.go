package models

import "gorm.io/gorm"

type AccessKey struct {
	gorm.Model
	UserId uint64 `json:"user_id"`
	Name   string `json:"name"`
	Key    string `json:"key"`
}

type InternalServerWhiteList struct {
	gorm.Model
	WorkerUID      string `json:"worker_uid" gorm:"index"`
	AllowWorkerUID string `json:"allow_worker_uid" gorm:"index"`
	Description    string `json:"description"`
}

type ExternalServerAKSK struct {
	gorm.Model
	WorkerUID      string `json:"worker_uid" gorm:"index"`
	AccessKey      string `json:"access_key"`
	SecretKey      string `json:"secret_key"`
	Description    string `json:"description"`
	Forever        bool   `json:"forever"`
	ExpirationTime string `json:"expiration_time"`
}

type ExternalServerToken struct {
	gorm.Model
	WorkerUID      string `json:"worker_uid" gorm:"index"`
	Token          string `json:"token"`
	Description    string `json:"description"`
	Forever        bool   `json:"forever"`
	ExpirationTime string `json:"expiration_time"`
}

type AccessRule struct {
	gorm.Model
	WorkerUID   string `json:"worker_uid" gorm:"index"`
	RuleType    string `json:"rule_type"` // "internal", "aksk", "token", "sso", "open"
	Path        string `json:"path"`
	Description string `json:"description"`
}

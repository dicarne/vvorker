package models

import (
	"time"
	"vvorker/entities"
	"vvorker/utils/database"

	"gorm.io/gorm"
)

type WorkerMember struct {
	gorm.Model
	WorkerUID   string    `json:"worker_uid" gorm:"index:idx_worker_uid"`
	UserID      uint64    `json:"user_id" gorm:"index:idx_user_id"`
	UserName    string    `json:"user_name"`
	AddedBy     uint64    `json:"added_by"`     // 添加者的用户ID
	AddedByName string    `json:"added_by_name"` // 添加者的用户名
	JoinedAt    time.Time `json:"joined_at"`
}

func init() {
	go func() {
		db := database.GetDB()
		for err := db.AutoMigrate(&WorkerMember{}); err != nil; err = db.AutoMigrate(&WorkerMember{}) {
			time.Sleep(5 * time.Second)
		}
	}()
}

func (w *WorkerMember) TableName() string {
	return "worker_members"
}

// AddMember 添加协作者
func AddWorkerMember(workerUID string, userID uint64, userName string, addedBy uint64, addedByName string) error {
	db := database.GetDB()
	member := &WorkerMember{
		WorkerUID:   workerUID,
		UserID:      userID,
		UserName:    userName,
		AddedBy:     addedBy,
		AddedByName: addedByName,
		JoinedAt:    time.Now(),
	}
	return db.Create(member).Error
}

// RemoveMember 移除协作者
func RemoveWorkerMember(workerUID string, userID uint64) error {
	db := database.GetDB()
	return db.Where(&WorkerMember{
		WorkerUID: workerUID,
		UserID:    userID,
	}).Delete(&WorkerMember{}).Error
}

// GetWorkerMembers 获取Worker的所有协作者
func GetWorkerMembers(workerUID string) ([]*WorkerMember, error) {
	var members []*WorkerMember
	db := database.GetDB()
	if err := db.Where(&WorkerMember{WorkerUID: workerUID}).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

// IsWorkerMember 检查用户是否是Worker的协作者
func IsWorkerMember(workerUID string, userID uint64) bool {
	var member WorkerMember
	db := database.GetDB()
	result := db.Where(&WorkerMember{
		WorkerUID: workerUID,
		UserID:    userID,
	}).First(&member)
	return result.Error == nil
}

// GetWorkerOwner 获取Worker的拥有者
func GetWorkerOwner(workerUID string) (uint64, error) {
	var worker Worker
	db := database.GetDB()
	if err := db.Where(&Worker{Worker: &entities.Worker{UID: workerUID}}).First(&worker).Error; err != nil {
		return 0, err
	}
	return worker.UserID, nil
}

// GetUserCollaboratedWorkers 获取用户参与协作的所有Worker
func GetUserCollaboratedWorkers(userID uint64) ([]string, error) {
	var workerUIDs []string
	db := database.GetDB()
	if err := db.Model(&WorkerMember{}).Where(&WorkerMember{UserID: userID}).Pluck("worker_uid", &workerUIDs).Error; err != nil {
		return nil, err
	}
	return workerUIDs, nil
}

// CanManageMembers 检查用户是否可以管理成员（只有拥有者可以）
func CanManageMembers(workerUID string, userID uint64) (bool, error) {
	ownerID, err := GetWorkerOwner(workerUID)
	if err != nil {
		return false, err
	}
	return ownerID == userID, nil
}

package models

import (
	"vvorker/entities"
	"vvorker/utils/database"
	"vvorker/utils/secret"

	"gorm.io/gorm"
)

const (
	ErrInvalidParams = "invalid params"
)

type User struct {
	gorm.Model
	UserName  string `json:"user_name" gorm:"unique"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Status    int    `json:"status"`
	Role      string `json:"role"`
	OtpSecret string `json:"otp_secret"`
}

func (u *User) TableName() string {
	return "users"
}

func CreateUser(user *User) error {
	if hashedPass, err := secret.HashPassword(user.Password); err != nil {
		return err
	} else {
		user.Password = hashedPass
	}

	return database.GetDB().Create(user).Error
}

func AdminGetUserNumber() (int64, error) {
	var count int64
	db := database.GetDB()

	if err := db.Model(&User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetUserByUserID(userID uint) (*User, error) {
	var user User
	db := database.GetDB()

	if err := db.Where(&User{
		Model: gorm.Model{ID: userID},
	}).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUserName(userName string) (*User, error) {
	var user User
	db := database.GetDB()

	if err := db.Where(&User{
		UserName: userName,
	}).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	db := database.GetDB()

	if err := db.Where(&User{
		Email: email,
	}).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(userID uint, user *User) error {
	if hashedPass, err := secret.HashPassword(user.Password); err != nil {
		return err
	} else {
		user.Password = hashedPass
	}
	db := database.GetDB()

	return db.Model(&User{
		Model: gorm.Model{ID: userID},
	}).Updates(user).Error
}

func DeleteUser(userID uint) error {
	db := database.GetDB()

	// 开始事务
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 1. 删除用户参与协作的 WorkerMember 记录
	if err := tx.Unscoped().Where(&WorkerMember{UserID: uint64(userID)}).Delete(&WorkerMember{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. 获取用户拥有的 workers
	var workers []*Worker
	if err := tx.Where(&Worker{Worker: &entities.Worker{UserID: uint64(userID)}}).Find(&workers).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 3. 删除用户拥有的 workers
	for _, worker := range workers {
		// 删除 worker 本身（会自动删除相关的 WorkerCopy、WorkerVersion、File、Task、AccessRule、ResponseLog 等记录）
		if err := worker.Delete(); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 4. 删除用户
	if err := tx.Unscoped().Delete(&User{
		Model: gorm.Model{ID: userID},
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func ListUsers(page, pageSize int) ([]*User, error) {
	var users []*User
	db := database.GetDB()

	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func CountUsers() (int64, error) {
	var count int64
	db := database.GetDB()

	if err := db.Model(&User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetUserByUserNameAndPassword(userName, password string) (*User, error) {
	var user User
	db := database.GetDB()

	if err := db.Where(&User{
		UserName: userName,
		Password: password,
	}).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CheckUserPassword(userNameOrEmail, password string) (bool, error) {
	var user User
	db := database.GetDB()

	if err := db.Where(&User{
		UserName: userNameOrEmail,
	}).Or(&User{
		Email: userNameOrEmail}).First(&user).Error; err != nil {
		return false, err
	}
	return secret.CheckPasswordHash(password, user.Password), nil
}

func CheckUserNameAndEmail(userName, email string) error {
	var user User
	db := database.GetDB()

	if err := db.Where(&User{
		UserName: userName,
	}).Or(&User{
		Email: email,
	}).First(&user).Error; err != nil {
		return err
	}
	return nil
}

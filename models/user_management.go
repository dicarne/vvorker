package models

import (
	"errors"
	"vvorker/common"
	"vvorker/utils/database"
	"vvorker/utils/secret"

	"gorm.io/gorm"
)

// AdminGetUserByID 通过ID获取用户（管理员操作）
func AdminGetUserByID(id uint) (*User, error) {
	return GetUserByUserID(id)
}

// AdminGetUserByUsername 通过用户名获取用户（管理员操作）
func AdminGetUserByUsername(username string) (*User, error) {
	return GetUserByUserName(username)
}

// AdminCreateUser 创建新用户（管理员操作）
func AdminCreateUser(username, password, email string) (*User, error) {
	// 检查用户名是否已存在
	existingUser, err := GetUserByUserName(username)
	if err == nil && existingUser != nil {
		return nil, errors.New("username already exists")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 对密码进行哈希处理
	hashedPass, err := secret.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// 创建新用户
	user := &User{
		UserName: username,
		Password: hashedPass,
		Email:    email,
		Role:     common.UserRoleNormal,
		Status:   common.UserStatusNormal,
	}

	if err := database.GetDB().Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// AdminUpdateUser 更新用户信息（管理员操作）
func AdminUpdateUser(user *User) error {
	// 如果密码被更改，需要哈希处理
	if user.Password != "" {
		hashedPass, err := secret.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPass
	}

	return database.GetDB().Save(user).Error
}

// AdminUpdateUserStatus 更新用户状态（管理员操作）
func AdminUpdateUserStatus(userID uint, status int) error {
	result := database.GetDB().Model(&User{}).Where("id = ?", userID).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// AdminSearchUsers 搜索用户（管理员操作）
func AdminSearchUsers(query string, role string, status int, page, pageSize int) ([]User, int64, error) {
	db := database.GetDB().Model(&User{})
	var total int64

	// 构建查询条件
	if query != "" {
		db = db.Where("user_name LIKE ?", "%"+query+"%")
	}
	if role != "" {
		db = db.Where("role = ?", role)
	}
	if status != 0 {
		db = db.Where("status = ?", status)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	var users []User
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// 清除密码字段
	for i := range users {
		users[i].Password = ""
	}

	return users, total, nil
}

// AdminBatchUpdateUserStatus 批量更新用户状态（管理员操作）
func AdminBatchUpdateUserStatus(userIDs []uint, status int) error {
	result := database.GetDB().Model(&User{}).Where("id IN ?", userIDs).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no users found")
	}
	return nil
}

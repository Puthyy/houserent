package db

import "houserent/model"

func AddUser(user *model.User) error {
	return Db.Create(user).Error
}

func DeleteUser(user *model.User) error {
	return Db.Delete(user).Error
}

func UpdateUser(user *model.User) error {
	return Db.Save(user).Error
}

// FindUserByID 通过ID查找用户
func FindUserByID(id uint) (*model.User, error) {
	var user model.User
	err := Db.First(&user, id).Error
	return &user, err
}

// FindUserByUsername 通过用户名查找用户（用于登录）
func FindUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := Db.Where("username = ?", username).First(&user).Error
	return &user, err
}

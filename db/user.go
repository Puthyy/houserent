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

func FindUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := Db.Where("username = ?", username).First(&user).Error
	return &user, err
}

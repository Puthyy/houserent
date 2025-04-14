package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"houserent/model"
)

var Db *gorm.DB

func InitDb() error {
	var err error
	Db, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	err = Db.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}

	return nil
}

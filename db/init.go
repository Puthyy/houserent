package db

import (
	"houserent/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	err = Db.AutoMigrate(&model.Listing{})
	if err != nil {
		return err
	}
	err = Db.AutoMigrate(&model.Transaction{})
	if err != nil {
		return err
	}
	err = Db.AutoMigrate(&model.Review{})
	if err != nil {
		return err
	}

	return nil
}

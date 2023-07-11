package db

import (
	"log"

	"github.com/fadilmuh22/restskuy/internal/model"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	db, err := gorm.Open(mysql.Open(viper.GetString("DB_URL")), &gorm.Config{})

	db.AutoMigrate(&model.User{}, &model.Product{})

	if err != nil {
		log.Fatal(err)
	}
	return db
}

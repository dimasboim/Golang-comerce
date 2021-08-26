package config

import (
	"Day15/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() {
	dsn := os.Getenv("Connection")

	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error Connect db")

	}

	Db.AutoMigrate(&models.Product_display{})
	Db.AutoMigrate(&models.Product_warehouse{})
	Db.AutoMigrate(&models.User{})
	Db.AutoMigrate(&models.Cart{})
	Db.AutoMigrate(&models.Transaksi{})
	Db.AutoMigrate(&models.Transaksi_detail{})

}

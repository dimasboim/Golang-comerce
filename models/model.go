package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username  string `gorm:"unique_index"`
	Firstname string `sql:"type:text"`
	Lastname  string `sql:"type:text"`
	Email     string `sql:"type:text"`
	Address   string `sql:"type:text"`
	Role      int64
	Token     string `sql:"type:text"`
	Password  string `sql:"type:text"`
}

type Product_warehouse struct {
	gorm.Model

	Sku     string `gorm:"unique_index"`
	Name    string `sql:"type:text"`
	User_id uint
	Price   float64
	Qty     int64
}

type Product_display struct {
	gorm.Model
	User_id uint
	Sku     string `gorm:"unique_index"`
	Name    string `sql:"type:text"`

	Price float64
	Qty   int64
}
type Cart struct {
	gorm.Model

	Sku     string `sql:"type:text"`
	Qty     int64
	User_id uint
}

type Transaksi struct {
	gorm.Model
	User_id uint
	Total   float64
}

type Transaksi_detail struct {
	gorm.Model
	Id_transaksi uint
	Sku          string `sql:"type:text"`
	Qty          int64
	Price        float64
	Subtotal     float64
	User_id      uint
}

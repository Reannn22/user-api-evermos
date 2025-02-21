package entities

import (
	"time"

	"gorm.io/gorm"
)

type ProductLog struct {
	gorm.Model
	ID            uint    `gorm:"primaryKey;column:id"`
	IDProduk      uint    `gorm:"column:id_produk;not null"`
	NamaProduk    string  `gorm:"column:nama_produk;size:255;not null"`
	Slug          string  `gorm:"column:slug;size:255;not null"`
	HargaReseller string  `gorm:"column:harga_reseller;size:255;not null"`
	HargaKonsumen string  `gorm:"column:harga_konsumen;size:255;not null"`
	Deskripsi     *string `gorm:"column:deskripsi;type:text;default:null"`
	IDToko        uint    `gorm:"column:id_toko;not null"`
	IDCategory    uint    `gorm:"column:id_category;not null"`
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	Store         Store    `gorm:"foreignKey:IDToko"`
	Category      Category `gorm:"foreignKey:IDCategory"`
	Product       Product  `gorm:"foreignKey:IDProduk"`
}

func (ProductLog) TableName() string {
	return "log_produk"
}

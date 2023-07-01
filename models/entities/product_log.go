package entities

import (
	"time"

	"gorm.io/gorm"
)

type ProductLog struct {
	gorm.Model
	ID            uint    `gorm:"primaryKey"`
	IDProduk      uint    `gorm:"not null"`
	NamaProduk    string  `gorm:"size:255;not null"`
	Slug          string  `gorm:"size:255;not null"`
	HargaReseller string  `gorm:"size:255;not null"`
	HargaKonsumen string  `gorm:"size:255;not null"`
	Deskripsi     *string `gorm:"type:text;default:null"`
	IDToko        uint    `gorm:"not null"`
	IDCategory    uint    `gorm:"not null"`
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	Store         Store    `gorm:"foreignKey:IDToko"`
	Category      Category `gorm:"foreignKey:IDCategory"`
	Product       Product  `gorm:"foreignKey:IDProduk"`
}

func (ProductLog) TableName() string {
	return "log_produk"
}

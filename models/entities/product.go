package entities

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID             uint    `gorm:"primaryKey"`
	NamaProduk     string  `gorm:"size:255;not null"`
	Slug           string  `gorm:"size:255;not null"`
	HargaReseller  string  `gorm:"size:255;not null"`
	HargaKonsumen  string  `gorm:"size:255;not null"`
	Stok           int     `gorm:"not null"`
	Deskripsi      *string `gorm:"type:text;default:null"`
	IDToko         uint    `gorm:"not null"`
	IDCategory     uint    `gorm:"not null"`
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
	Store          Store            `gorm:"foreignKey:IDToko;references:ID"`
	Category       Category         `gorm:"foreignKey:IDCategory;references:ID"`
	ProductPicture []ProductPicture `gorm:"foreignKey:IDProduk;references:ID"`
}

func (Product) TableName() string {
	return "produk" // This table has 69 rows
}

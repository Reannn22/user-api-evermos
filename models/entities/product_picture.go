package entities

import (
	"time"

	"gorm.io/gorm"
)

type ProductPicture struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	IDProduk  uint   `gorm:"not null"`
	Url       string `gorm:"size:255;not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (ProductPicture) TableName() string {
	return "foto_produk"
}

package entities

import (
	"time"

	"gorm.io/gorm"
)

type Store struct {
	gorm.Model
	ID        uint    `gorm:"primaryKey"`
	IDUser    uint    `gorm:"not null"`
	NamaToko  *string `gorm:"size:255;default:null"`
	UrlFoto   *string `gorm:"size:255;default:null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (Store) TableName() string {
	return "toko"
}

package entities

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey"`
	IDUser       uint   `gorm:"not null"`
	JudulAlamat  string `gorm:"size:255;not null"`
	NamaPenerima string `gorm:"size:255;not null"`
	NoTelp       string `gorm:"size:255;not null"`
	DetailAlamat string `gorm:"size:255;not null"`
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	User         User `gorm:"foreignKey:IDUser"`
}

func (Address) TableName() string {
	return "alamat"
}

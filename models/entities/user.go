package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uint      `gorm:"primaryKey"`
	Nama         string    `gorm:"size:255;not null"`
	KataSandi    string    `gorm:"size:255;not null"`
	Notelp       string    `gorm:"size:255;not null;index:idx_notelp,unique"`
	TanggalLahir time.Time `gorm:"type:date;not null"`
	JenisKelamin string    `gorm:"size:255;not null"`
	Tentang      *string   `gorm:"type:text;default:null"`
	Pekerjaan    string    `gorm:"size:255;not null"`
	Email        string    `gorm:"size:255;not null;index:idx_email,unique"`
	IDProvinsi   string    `gorm:"size:255;not null"`
	IDKota       string    `gorm:"size:255;not null"`
	IsAdmin      bool      `gorm:"type:boolean;default:false"`
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

func (User) TableName() string {
	return "user"
}

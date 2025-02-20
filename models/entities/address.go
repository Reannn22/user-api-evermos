package entities

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	ID           uint   `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	JudulAlamat  string `gorm:"column:judul_alamat"`
	NamaPenerima string `gorm:"column:nama_penerima"`
	NoTelp       string `gorm:"column:no_telp"`
	DetailAlamat string `gorm:"column:detail_alamat"`
	IDUser       uint   `gorm:"column:id_user"`
	IDProvinsi   string `gorm:"column:id_provinsi"`
	IDKota       string `gorm:"column:id_kota"`
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (Address) TableName() string {
	return "alamat"
}

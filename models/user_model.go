package models

import "time"

// Request
type UserRequest struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Nama         string    `json:"nama" binding:"required"`
	KataSandi    string    `json:"kata_sandi" binding:"required"`
	NoTelp       string    `json:"no_telp" binding:"required"`
	TanggalLahir string    `json:"tanggal_lahir" binding:"required"`
	JenisKelamin string    `json:"jenis_kelamin" gorm:"type:varchar(255)" binding:"required"`
	Tentang      string    `json:"tentang" gorm:"type:text"`
	Pekerjaan    string    `json:"pekerjaan" binding:"required"`
	Email        string    `json:"email" binding:"required"`
	IDProvinsi   string    `json:"id_provinsi" binding:"required"`
	IDKota       string    `json:"id_kota" binding:"required"`
	IsAdmin      bool      `json:"is_admin" gorm:"default:false"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// Response
type UserResponse struct {
	ID           uint           `json:"id"`
	Nama         string         `json:"nama"`
	KataSandi    string         `json:"kata_sandi"`
	Notelp       string         `json:"notelp"`
	TanggalLahir time.Time      `json:"tanggal_lahir"`
	JenisKelamin string         `json:"jenis_kelamin"`
	Tentang      *string        `json:"tentang"`
	Pekerjaan    string         `json:"pekerjaan"`
	Email        string         `json:"email"`
	IDProvinsi   ProvinceDetail `json:"id_provinsi"`
	IDKota       CityDetail     `json:"id_kota"`
	IsAdmin      bool           `json:"is_admin"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

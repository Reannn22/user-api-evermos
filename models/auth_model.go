package models

import "time"

// Request
type RegisterRequest struct {
	Nama         string `json:"nama" binding:"required"`
	KataSandi    string `json:"kata_sandi" binding:"required"`
	NoTelp       string `json:"no_telp" binding:"required"`
	TanggalLahir string `json:"tanggal_lahir" binding:"required"`
	Pekerjaan    string `json:"pekerjaan" binding:"required"`
	Email        string `json:"email" binding:"required"`
	IDProvinsi   string `json:"id_provinsi" binding:"required"`
	IDKota       string `json:"id_kota" binding:"required"`
}

type LoginRequest struct {
	NoTelp    string `json:"no_telp" binding:"required"`
	KataSandi string `json:"kata_sandi" binding:"required"`
}

// Response
type LoginResponse struct {
	Nama         string    `json:"nama"`
	NoTelp       string    `json:"no_telp"`
	TanggalLahir string    `json:"tanggal_lahir"`
	Tentang      *string   `json:"tentang"`
	Pekerjaan    string    `json:"pekerjaan"`
	Email        string    `json:"email"`
	IDProvinsi   *Province `json:"id_provinsi"`
	IDKota       *City     `json:"id_kota"`
	Token        string    `json:"token"`
}

// process mapping
type RegisterProcess struct {
	Nama         string
	NoTelp       string
	Email        string
	KataSandi    string
	TanggalLahir time.Time
	Pekerjaan    string
	IDProvinsi   string
	IDKota       string
}

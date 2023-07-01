package models

// Request
type UserRequest struct {
	Nama         string `json:"nama" binding:"required"`
	KataSandi    string `json:"kata_sandi" binding:"required"`
	NoTelp       string `json:"no_telp" binding:"required"`
	TanggalLahir string `json:"tanggal_lahir" binding:"required"`
	Pekerjaan    string `json:"pekerjaan" binding:"required"`
	Email        string `json:"email" binding:"required"`
	IDProvinsi   string `json:"id_provinsi" binding:"required"`
	IDKota       string `json:"id_kota" binding:"required"`
}

// Response
type UserResponse struct {
	Nama         string    `json:"nama"`
	NoTelp       string    `json:"no_telp"`
	TanggalLahir string    `json:"tanggal_lahir"`
	Tentang      *string   `json:"tentang"`
	Pekerjaan    string    `json:"pekerjaan"`
	Email        string    `json:"email"`
	IDProvinsi   *Province `json:"id_provinsi"`
	IDKota       *City     `json:"id_kota"`
}

package entities

import "time"

type Alamat struct {
    ID           uint       `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
    JudulAlamat  string     `json:"judul_alamat" gorm:"column:judul_alamat"`
    NamaPenerima string     `json:"nama_penerima" gorm:"column:nama_penerima"`
    NoTelp       string     `json:"no_telp" gorm:"column:no_telp"`
    DetailAlamat string     `json:"detail_alamat" gorm:"column:detail_alamat"`
    IDUser       uint       `json:"id_user" gorm:"column:id_user"`
    IDProvinsi   string     `json:"id_provinsi" gorm:"column:id_provinsi"`
    IDKota       string     `json:"id_kota" gorm:"column:id_kota"`
    CreatedAt    *time.Time `json:"created_at"`
    UpdatedAt    *time.Time `json:"updated_at"`
    DeletedAt    *time.Time `json:"deleted_at" gorm:"index"`
}

func (Alamat) TableName() string {
    return "alamat"
}

package entities

import "time"

type FotoProduk struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	IDProduk  uint       `json:"id_produk" gorm:"column:id_produk"`
	PhotoID   uint       `json:"photo_id" gorm:"column:photo_id"`
	Url       string     `json:"url" gorm:"column:url"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (FotoProduk) TableName() string {
	return "foto_produk"
}

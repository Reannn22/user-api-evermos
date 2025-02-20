package entities

import (
	"time"

	"gorm.io/gorm"
)

type TrxDetail struct {
	gorm.Model
	IDTrx       uint       `gorm:"column:id_trx"`
	IDLogProduk uint       `gorm:"column:id_log_produk"`
	IDToko      uint       `gorm:"column:id_toko"`
	Kuantitas   int        `gorm:"column:kuantitas"`
	HargaTotal  int        `gorm:"column:harga_total"`
	Trx         Trx        `gorm:"foreignKey:IDTrx"`
	ProductLog  ProductLog `gorm:"foreignKey:IDLogProduk"`
	Store       Store      `gorm:"foreignKey:IDToko"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
}

func (TrxDetail) TableName() string {
	return "detail_trx"
}

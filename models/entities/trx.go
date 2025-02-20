package entities

import (
	"time"

	"gorm.io/gorm"
)

type Trx struct {
	gorm.Model
	KodeInvoice      string      `gorm:"column:kode_invoice"`
	MethodBayar      string      `gorm:"column:method_bayar"`
	AlamatPengiriman uint        `gorm:"column:alamat_pengiriman"`
	IDUser           uint        `gorm:"column:id_user"`
	HargaTotal       int         `gorm:"column:harga_total"`
	Address          Address     `gorm:"foreignKey:ID;references:AlamatPengiriman"`
	TrxDetail        []TrxDetail `gorm:"foreignKey:IDTrx;references:ID"`
	CreatedAt        *time.Time  `json:"created_at"`
	UpdatedAt        *time.Time  `json:"updated_at"`
	DeletedAt        *time.Time  `json:"deleted_at" gorm:"index"`
}

func (Trx) TableName() string {
	return "trx"
}

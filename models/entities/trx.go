package entities

import (
	"time"

	"gorm.io/gorm"
)

type Trx struct {
	gorm.Model
	ID               uint   `gorm:"primaryKey"`
	IDUser           uint   `gorm:"not null"`
	AlamatPengiriman uint   `gorm:"not null"`
	HargaTotal       int    `gorm:"not null"`
	KodeInvoice      string `gorm:"size:255;not null"`
	MethodBayar      string `gorm:"size:255;not null"`
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
	Address          Address     `gorm:"foreignKey:AlamatPengiriman"`
	User             User        `gorm:"foreignKey:IDUser"`
	TrxDetail        []TrxDetail `gorm:"foreignKey:IDTrx"`
}

func (Trx) TableName() string {
	return "trx"
}

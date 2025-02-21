package models

import "time"

type ProductLogProcess struct {
	ProductID     uint   `json:"product_id"`
	NamaProduk    string `json:"nama_produk"`
	Slug          string `json:"slug"`
	HargaReseller string `json:"harga_reseller"`
	HargaKonsumen string `json:"harga_konsumen"`
	Stok          int
	Deskripsi     string `json:"deskripsi"`
	StoreID       uint   `json:"store_id"`
	CategoryID    uint   `json:"category_id"`
	Kuantitas     int
	HargaTotal    int
}

type ProductLogResponse struct {
	ProductID     uint      `json:"product_id"`
	NamaProduk    string    `json:"nama_produk"`
	Slug          string    `json:"slug"`
	HargaReseller string    `json:"harga_reseller"`
	HargaKonsumen string    `json:"harga_konsumen"`
	Stok          int       `json:"stok"`
	Deskripsi     string    `json:"deskripsi"`
	StoreID       uint      `json:"store_id"`
	CategoryID    uint      `json:"category_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

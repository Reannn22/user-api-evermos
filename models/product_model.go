package models

import "time"

// Request
type ProductRequest struct {
	NamaProduk    string   `json:"nama_produk" form:"nama_produk"`
	CategoryID    uint     `json:"category_id" form:"category_id"`
	StoreID       uint     `json:"store_id"`
	HargaReseller string   `json:"harga_reseller" form:"harga_reseller"`
	HargaKonsumen string   `json:"harga_konsumen" form:"harga_konsumen"`
	Stok          int      `json:"stok" form:"stok"`
	Deskripsi     string   `json:"deskripsi" form:"deskripsi"`
	PhotoURLs     []string `json:"photo_urls" form:"photo_urls"` // Changed from Photos to PhotoURLs
}

// Response
type ProductResponse struct {
	ID            uint                     `json:"id"`
	NamaProduk    string                   `json:"nama_produk"`
	Slug          string                   `json:"slug"`
	HargaReseller string                   `json:"harga_reseler"`
	HargaKonsumen string                   `json:"harga_konsumen"`
	Stok          int                      `json:"stok"`
	Deskripsi     *string                  `json:"deskripsi"`
	Store         StoreResponse            `json:"toko"`
	Category      CategoryResponse         `json:"category"`
	Photos        []ProductPictureResponse `json:"photos"`
	CreatedAt     *time.Time               `json:"created_at"`
	UpdatedAt     *time.Time               `json:"updated_at"`
}

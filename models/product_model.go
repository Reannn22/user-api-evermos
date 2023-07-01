package models

// Request
type ProductRequest struct {
	NamaProduk    string   `json:"nama_produk"`
	CategoryID    uint     `json:"category_id"`
	StoreID       uint     `json:"store_id"`
	HargaReseller string   `json:"harga_reseler"`
	HargaKonsumen string   `json:"harga_konsumen"`
	Stok          int      `json:"stok"`
	Deskripsi     string   `json:"deskripsi"`
	Photos        []string `json:"photos"`
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
}

package models

// Request
type TransactionDetailRequest struct {
	ProductID uint `json:"id_produk"` // Changed from IDProduk to ProductID
	Kuantitas int  `json:"quantity"`  // Changed from Quantity to Kuantitas
}

// Response
type TransactionDetailResponse struct {
	ID         uint            `json:"id"`
	Kuantitas  int             `json:"kuantitas"`
	HargaTotal int             `json:"harga_total"`
	Store      StoreResponse   `json:"toko"`
	Product    ProductResponse `json:"product"`
}

type TransactionDetailProcess struct {
	TrxID        uint
	LogProductID uint
	StoreID      uint
	Kuantitas    int
	HargaTotal   int
}

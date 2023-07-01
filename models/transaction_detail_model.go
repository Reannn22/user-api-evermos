package models

// Request
type TransactionDetailRequest struct {
	ProductID uint `json:"product_id"`
	Kuantitas int  `json:"kuantitas"`
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

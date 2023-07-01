package models

type ProductLogProcess struct {
	ProductID     uint
	NamaProduk    string
	Slug          string
	HargaReseller string
	HargaKonsumen string
	Stok          int
	Deskripsi     string
	StoreID       uint
	CategoryID    uint
	Kuantitas     int
	HargaTotal    int
}

type ProductLogResponse struct {
	ProductID     uint
	NamaProduk    string
	Slug          string
	HargaReseller string
	HargaKonsumen string
	Stok          int
	Deskripsi     string
	StoreID       uint
	CategoryID    uint
}

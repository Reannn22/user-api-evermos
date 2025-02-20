package models

import "time"

// Request
type TransactionRequest struct {
	MethodBayar      string                     `json:"method_bayar"`
	AlamatPengiriman uint                       `json:"alamat_kirim"` // Changed from AlamatKirim
	DetailTrx        []TransactionDetailRequest `json:"detail_trx"`
}

// Response
type TransactionResponse struct {
	ID                 uint                        `json:"id"`
	HargaTotal         int                         `json:"harga_total"`
	KodeInvoice        string                      `json:"kode_invoice"`
	MethodBayar        string                      `json:"method_bayar"`
	Address            AddressResponse             `json:"alamat_kirim"`
	TransactionDetails []TransactionDetailResponse `json:"detail_trx"`
	CreatedAt          *time.Time                  `json:"created_at"`
	UpdatedAt          *time.Time                  `json:"updated_at"`
}

type TransactionProcess struct {
	MethodBayar      string
	KodeInvoice      string
	AlamatPengiriman uint
	UserID           uint
	HargaTotal       int
}

type TransactionProcessData struct {
	Transaction TransactionProcess
	LogProduct  []ProductLogProcess
}

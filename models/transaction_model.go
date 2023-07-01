package models

// Request
type TransactionRequest struct {
	MethodBayar string                     `json:"method_bayar"`
	AlamatKirim uint                       `json:"alamat_kirim"`
	DetailTrx   []TransactionDetailRequest `json:"detail_trx"`
}

// Response
type TransactionResponse struct {
	ID                 uint                        `json:"id"`
	HargaTotal         int                         `json:"harga_total"`
	KodeInvoice        string                      `json:"kode_invoice"`
	MethodBayar        string                      `json:"method_bayar"`
	Address            AddressResponse             `json:"alamat_kirim"`
	TransactionDetails []TransactionDetailResponse `json:"detail_trx"`
}

type TransactionProcess struct {
	MethodBayar string
	KodeInvoice string
	AlamatKirim uint
	UserID      uint
	HargaTotal  int
}

type TransactionProcessData struct {
	Transaction TransactionProcess
	LogProduct  []ProductLogProcess
}

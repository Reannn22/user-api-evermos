package models

// TransactionUpdateRequest represents the request body for updating a transaction
type TransactionUpdateRequest struct {
	Status      string `json:"status"`
	MethodBayar string `json:"method_bayar"`
}

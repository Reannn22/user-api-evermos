package models

// Request
type CategoryRequest struct {
	NamaCategory string `json:"nama_category" binding:"required"`
}

// Response
type CategoryResponse struct {
	ID           uint   `json:"id"`
	NamaCategory string `json:"nama_category"`
}

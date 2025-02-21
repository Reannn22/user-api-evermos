package models

import "time"

// Request
type FotoProdukRequest struct {
	ProductID uint   `json:"product_id" form:"product_id"`
	PhotoID   uint   `json:"photo_id" form:"photo_id"` // Add this field
	URL       string `json:"url" form:"url"`
}

// Response
type FotoProdukResponse struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"product_id"`
	PhotoID   uint      `json:"photo_id"` // Add this field
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

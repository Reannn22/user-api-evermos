package models

import (
    "time"
    "gorm.io/gorm"
)

type Product struct {
    // ...existing code...
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

// Pastikan method BeforeCreate dan BeforeUpdate tetap ada
func (product *Product) BeforeCreate(tx *gorm.DB) error {
    product.CreatedAt = time.Now()
    product.UpdatedAt = time.Now()
    return nil
}

func (product *Product) BeforeUpdate(tx *gorm.DB) error {
    product.UpdatedAt = time.Now()
    return nil
}

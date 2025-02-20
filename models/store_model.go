package models

import (
	"mime/multipart"
	"time"
)

// Response
type StoreResponse struct {
	ID        uint       `json:"id"`
	NamaToko  *string    `json:"nama_toko"`
	UrlFoto   *string    `json:"url_foto"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type StoreUpdate struct {
	NamaToko *string
	UrlFoto  string
}

type StoreProcess struct {
	ID        uint
	UserID    uint
	NamaToko  *string
	URL       string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type File struct {
	File multipart.File `json:"file,omitempty"`
}

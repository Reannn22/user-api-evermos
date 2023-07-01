package models

import "mime/multipart"

// Response
type StoreResponse struct {
	ID       uint    `json:"id"`
	NamaToko *string `json:"nama_toko"`
	UrlFoto  *string `json:"url_foto"`
}

type StoreUpdate struct {
	NamaToko *string
	UrlFoto  string
}

type StoreProcess struct {
	ID       uint
	UserID   uint
	NamaToko *string
	URL      string
}

type File struct {
	File multipart.File `json:"file,omitempty"`
}

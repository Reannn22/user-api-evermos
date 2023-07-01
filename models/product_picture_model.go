package models

// response
type ProductPictureResponse struct {
	ID       uint   `json:"id"`
	IDProduk uint   `json:"product_id"`
	Url      string `json:"url"`
}

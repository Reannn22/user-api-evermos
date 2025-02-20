package models

// request
type ProductPictureRequest struct {
    ProductID uint   `json:"product_id" form:"product_id"` // Fix: changed from category_id to product_id
    URL       string `json:"url" form:"url"`               // Fix: changed from photo_url to url
}

// response
type ProductPictureResponse struct {
	ID       uint   `json:"id"`
	IDProduk uint   `json:"product_id"`
	Url      string `json:"url"`
}

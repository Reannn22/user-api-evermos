package repositories

import (
	"mini-project-evermos/models"
	"mini-project-evermos/models/entities"
	"time"

	"gorm.io/gorm"
)

type FotoProdukRepository interface {
	FindAll() ([]entities.FotoProduk, error)
	FindById(id uint) (entities.FotoProduk, error)
	FindByProductId(productId uint) ([]entities.FotoProduk, error)
	Create(foto models.FotoProdukRequest) (entities.FotoProduk, error)
	Update(id uint, foto models.FotoProdukRequest) (entities.FotoProduk, error)
	Delete(id uint) error
}

type fotoProdukRepositoryImpl struct {
	db *gorm.DB
}

func NewFotoProdukRepository(db *gorm.DB) FotoProdukRepository {
	return &fotoProdukRepositoryImpl{db}
}

func (r *fotoProdukRepositoryImpl) FindAll() ([]entities.FotoProduk, error) {
	var photos []entities.FotoProduk
	err := r.db.Find(&photos).Error
	return photos, err
}

func (r *fotoProdukRepositoryImpl) FindById(id uint) (entities.FotoProduk, error) {
	var photo entities.FotoProduk
	err := r.db.First(&photo, id).Error
	return photo, err
}

func (r *fotoProdukRepositoryImpl) FindByProductId(productId uint) ([]entities.FotoProduk, error) {
	var photos []entities.FotoProduk
	err := r.db.Where("id_produk = ?", productId).Find(&photos).Error
	return photos, err
}

func (r *fotoProdukRepositoryImpl) Create(foto models.FotoProdukRequest) (entities.FotoProduk, error) {
	now := time.Now()
	newPhoto := entities.FotoProduk{
		IDProduk:  foto.ProductID,
		PhotoID:   foto.PhotoID, // Add this field
		Url:       foto.URL,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	err := r.db.Create(&newPhoto).Error
	return newPhoto, err
}

func (r *fotoProdukRepositoryImpl) Update(id uint, foto models.FotoProdukRequest) (entities.FotoProduk, error) {
	var photo entities.FotoProduk

	// First find the existing photo
	err := r.db.First(&photo, id).Error
	if err != nil {
		return entities.FotoProduk{}, err
	}

	// Update the fields
	now := time.Now()
	photo.Url = foto.URL
	photo.PhotoID = foto.PhotoID // Add this line to update PhotoID
	photo.UpdatedAt = &now

	err = r.db.Save(&photo).Error
	if err != nil {
		return entities.FotoProduk{}, err
	}

	return photo, nil
}

func (r *fotoProdukRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entities.FotoProduk{}, id).Error
}

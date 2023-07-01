package repositories

import (
	"mini-project-evermos/models/entities"

	"gorm.io/gorm"
)

// Contract
type ProductPictureRepository interface {
	FindByProductId(product_id uint) ([]entities.ProductPicture, error)
}

type productPictureRepositoryImpl struct {
	database *gorm.DB
}

func NewProductPictureRepository(database *gorm.DB) ProductPictureRepository {
	return &productPictureRepositoryImpl{database}
}

func (repository *productPictureRepositoryImpl) FindByProductId(product_id uint) ([]entities.ProductPicture, error) {
	var product []entities.ProductPicture

	err := repository.database.Where("id_produk = ?", product_id).Find(&product).Error

	if err != nil {
		return product, err
	}

	return product, nil
}

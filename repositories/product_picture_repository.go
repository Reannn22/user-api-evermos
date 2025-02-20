package repositories

import (
	"mini-project-evermos/models/entities"

	"gorm.io/gorm"
)

// Contract
type ProductPictureRepository interface {
	FindByProductId(product_id uint) ([]entities.ProductPicture, error)
	Create(productPicture *entities.ProductPicture) error
	GetByID(id uint) (*entities.ProductPicture, error)
	VerifyProductExists(productID uint) bool
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

func (repository *productPictureRepositoryImpl) Create(productPicture *entities.ProductPicture) error {
	result := repository.database.Create(productPicture)
	return result.Error
}

func (repository *productPictureRepositoryImpl) GetByID(id uint) (*entities.ProductPicture, error) {
	var productPicture entities.ProductPicture
	result := repository.database.First(&productPicture, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &productPicture, nil
}

func (repository *productPictureRepositoryImpl) VerifyProductExists(productID uint) bool {
	var count int64
	repository.database.Model(&entities.Product{}).Where("id = ?", productID).Count(&count)
	return count > 0
}

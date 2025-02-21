package repositories

import (
	"mini-project-evermos/models"
	"mini-project-evermos/models/entities"

	"gorm.io/gorm"
)

type ProductLogRepository interface {
	Insert(input models.ProductLogProcess) (entities.ProductLog, error)
	FindAll() ([]entities.ProductLog, error)
	FindById(id uint) (entities.ProductLog, error)
	Update(id uint, input models.ProductLogProcess) (entities.ProductLog, error)
	Delete(id uint) error
}

type productLogRepositoryImpl struct {
	db *gorm.DB
}

func NewProductLogRepository(db *gorm.DB) ProductLogRepository {
	return &productLogRepositoryImpl{db}
}

func (repository *productLogRepositoryImpl) Insert(input models.ProductLogProcess) (entities.ProductLog, error) {
	productLog := entities.ProductLog{
		IDProduk:      input.ProductID,
		NamaProduk:    input.NamaProduk,
		Slug:          input.Slug,
		HargaReseller: input.HargaReseller,
		HargaKonsumen: input.HargaKonsumen,
		Deskripsi:     &input.Deskripsi,
		IDToko:        input.StoreID,
		IDCategory:    input.CategoryID,
	}

	err := repository.db.Create(&productLog).Error
	if err != nil {
		return entities.ProductLog{}, err
	}

	return productLog, nil
}

// Add new method
func (repository *productLogRepositoryImpl) FindAll() ([]entities.ProductLog, error) {
	var productLogs []entities.ProductLog
	err := repository.db.Find(&productLogs).Error
	if err != nil {
		return nil, err
	}
	return productLogs, nil
}

// Add new method
func (repository *productLogRepositoryImpl) FindById(id uint) (entities.ProductLog, error) {
	var productLog entities.ProductLog
	err := repository.db.Preload("Store").Preload("Category").First(&productLog, id).Error
	if err != nil {
		return entities.ProductLog{}, err
	}
	return productLog, nil
}

// Add new method
func (repository *productLogRepositoryImpl) Update(id uint, input models.ProductLogProcess) (entities.ProductLog, error) {
	// Try to find the record first
	var productLog entities.ProductLog
	if err := repository.db.Table("log_produk").Where("id = ?", id).First(&productLog).Error; err != nil {
		return entities.ProductLog{}, err
	}

	// Prepare update data
	updates := entities.ProductLog{
		IDProduk:      input.ProductID,
		NamaProduk:    input.NamaProduk,
		Slug:          input.Slug,
		HargaReseller: input.HargaReseller,
		HargaKonsumen: input.HargaKonsumen,
		Deskripsi:     &input.Deskripsi,
		IDToko:        input.StoreID,
		IDCategory:    input.CategoryID,
	}

	// Perform update
	if err := repository.db.Table("log_produk").Where("id = ?", id).Updates(&updates).Error; err != nil {
		return entities.ProductLog{}, err
	}

	// Fetch updated record
	var updatedLog entities.ProductLog
	if err := repository.db.Table("log_produk").Where("id = ?", id).First(&updatedLog).Error; err != nil {
		return entities.ProductLog{}, err
	}

	return updatedLog, nil
}

// Add new method
func (repository *productLogRepositoryImpl) Delete(id uint) error {
	result := repository.db.Table("log_produk").Where("id = ?", id).Delete(&entities.ProductLog{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

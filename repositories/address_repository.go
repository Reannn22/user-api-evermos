package repositories

import (
	"fmt"
	"mini-project-evermos/models/entities"

	"gorm.io/gorm"
)

// Contract
type AddressRepository interface {
	FindByUserId(id uint) ([]entities.Address, error)
	FindById(id uint) (entities.Address, error)
	Insert(address entities.Address) (bool, error)
	Update(id uint, address entities.Address) (bool, error)
	Destroy(id uint) (bool, error)
	FindByCondition(condition map[string]interface{}) (entities.Address, error)
}

type addressRepositoryImpl struct {
	database *gorm.DB
}

func NewAddressRepository(database *gorm.DB) AddressRepository {
	return &addressRepositoryImpl{database}
}

func (repository *addressRepositoryImpl) FindByUserId(id uint) ([]entities.Address, error) {
	var address []entities.Address
	err := repository.database.Where("id_user = ?", id).Find(&address).Error

	if err != nil {
		return address, err
	}

	return address, nil
}

func (repository *addressRepositoryImpl) FindById(id uint) (entities.Address, error) {
	var address entities.Address

	fmt.Printf("Looking for address with ID: %d in table: %s\n", id, entities.Address{}.TableName())

	result := repository.database.Debug().
		Table(entities.Address{}.TableName()).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&address)

	if result.Error != nil {
		fmt.Printf("Database error: %v\n", result.Error)
		return entities.Address{}, fmt.Errorf("address with ID %d not found", id)
	}

	fmt.Printf("Found address: %+v\n", address)
	return address, nil
}

func (repository *addressRepositoryImpl) Insert(address entities.Address) (bool, error) {
	err := repository.database.Create(&address).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (repository *addressRepositoryImpl) Update(id uint, address entities.Address) (bool, error) {
	err := repository.database.Model(&address).Where("id = ?", id).Updates(address).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (repository *addressRepositoryImpl) Destroy(id uint) (bool, error) {
	var address entities.Address
	err := repository.database.Where("id = ?", id).Delete(&address).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (repository *addressRepositoryImpl) FindByCondition(condition map[string]interface{}) (entities.Address, error) {
	var address entities.Address

	result := repository.database.Where(condition).Order("created_at DESC").First(&address)
	if result.Error != nil {
		return entities.Address{}, result.Error
	}

	return address, nil
}

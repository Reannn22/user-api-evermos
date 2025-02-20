package repositories

import (
	"fmt"
	"mini-project-evermos/models/entities"

	"gorm.io/gorm"
)

// Contract
type UserRepository interface {
	FindByNoTelp(no_telp string) (entities.User, error)
	FindById(id uint) (entities.User, error)
	Update(id uint, user entities.User) (bool, error)
	FindByEmail(email string) (entities.User, error)
	Delete(id uint) error
}

type userRepositoryImpl struct {
	database *gorm.DB
}

func NewUserRepository(database *gorm.DB) UserRepository {
	return &userRepositoryImpl{database}
}

func (repository *userRepositoryImpl) FindByNoTelp(no_telp string) (entities.User, error) {
	var user entities.User
	err := repository.database.Where("notelp = ?", no_telp).First(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (repository *userRepositoryImpl) FindById(id uint) (entities.User, error) {
	var user entities.User

	// Add debug print
	fmt.Printf("Finding user with ID: %d\n", id)

	err := repository.database.Where("id = ?", id).First(&user).Error
	if err != nil {
		fmt.Printf("Error finding user: %v\n", err)
		return user, err
	}

	return user, nil
}

func (repository *userRepositoryImpl) Update(id uint, user entities.User) (bool, error) {
	err := repository.database.Model(&user).Where("id = ?", id).Updates(user).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (repository *userRepositoryImpl) FindByEmail(email string) (entities.User, error) {
	var user entities.User

	err := repository.database.Where("email = ?", email).First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (repository *userRepositoryImpl) Delete(id uint) error {
	var user entities.User
	err := repository.database.Where("id = ?", id).Delete(&user).Error
	return err
}

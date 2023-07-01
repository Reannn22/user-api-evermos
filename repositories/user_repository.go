package repositories

import (
	"mini-project-evermos/models/entities"

	"gorm.io/gorm"
)

// Contract
type UserRepository interface {
	FindByNoTelp(no_telp string) (entities.User, error)
	FindById(id uint) (entities.User, error)
	Update(id uint, user entities.User) (bool, error)
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
	err := repository.database.Where("id = ?", id).First(&user).Error

	if err != nil {
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

package repositories

import (
	"mini-project-evermos/models/entities"

	"gorm.io/gorm"
)

// Contract
type AuthRepository interface {
	Register(user entities.User) (entities.User, error)
}

type authRepositoryImpl struct {
	database *gorm.DB
}

func NewAuthRepository(database *gorm.DB) AuthRepository {
	return &authRepositoryImpl{database}
}

func (repository *authRepositoryImpl) Register(user entities.User) (entities.User, error) {
	// First, try to find any existing records (including soft-deleted ones)
	var existingUser entities.User
	result := repository.database.Unscoped().Where("notelp = ? OR email = ?", user.Notelp, user.Email).First(&existingUser)

	if result.Error == nil {
		// If found, hard delete the existing record
		repository.database.Unscoped().Delete(&existingUser)
	}

	// Now create the new user
	err := repository.database.Create(&user).Error
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

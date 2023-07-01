package repositories

import (
	"mini-project-evermos/models"
	"mini-project-evermos/models/entities"

	"gorm.io/gorm"
)

// Contract
type AuthRepository interface {
	Register(auth models.RegisterProcess) error
}

type authRepositoryImpl struct {
	database *gorm.DB
}

func NewAuthRepository(database *gorm.DB) AuthRepository {
	return &authRepositoryImpl{database}
}

func (repository *authRepositoryImpl) Register(auth models.RegisterProcess) error {
	tx := repository.database.Begin()

	user := &entities.User{
		Nama:         auth.Nama,
		Notelp:       auth.NoTelp,
		Email:        auth.Email,
		KataSandi:    auth.KataSandi,
		Pekerjaan:    auth.Pekerjaan,
		TanggalLahir: auth.TanggalLahir,
		IDProvinsi:   auth.IDProvinsi,
		IDKota:       auth.IDKota,
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&entities.Store{
		IDUser: user.ID,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

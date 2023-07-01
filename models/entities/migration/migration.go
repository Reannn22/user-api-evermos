package migration

import (
	"mini-project-evermos/models/entities"

	"gorm.io/gorm"
)

func Migration(database *gorm.DB) {
	database.AutoMigrate(
		&entities.User{},
		&entities.Store{},
		&entities.Category{},
		&entities.Product{},
		&entities.ProductPicture{},
		&entities.ProductLog{},
		&entities.Address{},
		&entities.Trx{},
		&entities.TrxDetail{},
	)
}

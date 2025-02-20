package migration

import (
	"mini-project-evermos/models/entities"

	"gorm.io/gorm"
)

// Export AutoMigrate function
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.Address{},
		&entities.User{},
		&entities.Store{},
		&entities.Category{},
		&entities.Product{},
		&entities.ProductPicture{},
		&entities.Trx{},
		&entities.TrxDetail{},
		&entities.ProductLog{},
	)
}

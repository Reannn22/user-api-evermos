package repositories

import (
	"math"
	"mini-project-evermos/models"
	"mini-project-evermos/models/entities"
	"mini-project-evermos/models/responder"

	"gorm.io/gorm"
)

// Contract
type StoreRepository interface {
	FindAllPagination(pagination responder.Pagination) (responder.Pagination, error)
	FindById(id uint) (entities.Store, error)
	FindByUserId(id uint) (entities.Store, error)
	Update(id uint, store entities.Store) (bool, error)
	Insert(store entities.Store) (entities.Store, error)
	Delete(id uint) (bool, error)
}

type storeRepositoryImpl struct {
	database *gorm.DB
}

func NewStoreRepository(database *gorm.DB) StoreRepository {
	return &storeRepositoryImpl{database}
}

func (repository *storeRepositoryImpl) FindAllPagination(pagination responder.Pagination) (responder.Pagination, error) {
	var stores []entities.Store
	var totalRows int64

	query := repository.database.Model(&entities.Store{})
	if pagination.Keyword != "" {
		query = query.Where("nama_toko LIKE ?", "%"+pagination.Keyword+"%")
	}
	query.Count(&totalRows)

	err := query.
		Limit(pagination.Limit).
		Offset(pagination.GetOffset()).
		Find(&stores).Error

	if err != nil {
		return responder.Pagination{}, err
	}

	var responses []models.StoreResponse
	for _, store := range stores {
		responses = append(responses, models.StoreResponse{
			ID:        store.ID,
			NamaToko:  store.NamaToko,
			UrlFoto:   store.UrlFoto,
			CreatedAt: store.CreatedAt,
			UpdatedAt: store.UpdatedAt,
		})
	}

	pagination.Rows = responses
	pagination.TotalRows = totalRows
	pagination.TotalPages = int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))

	return pagination, nil
}

func (repository *storeRepositoryImpl) FindById(id uint) (entities.Store, error) {
	var store entities.Store
	err := repository.database.Where("id = ?", id).First(&store).Error

	if err != nil {
		return store, err
	}

	return store, nil
}

func (repository *storeRepositoryImpl) FindByUserId(id uint) (entities.Store, error) {
	var store entities.Store
	err := repository.database.Where("id_user = ?", id).First(&store).Error

	if err != nil {
		return store, err
	}

	return store, nil
}

func (repository *storeRepositoryImpl) Update(id uint, store entities.Store) (bool, error) {
	err := repository.database.Model(&store).Where("id = ?", id).Updates(store).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (repository *storeRepositoryImpl) Insert(store entities.Store) (entities.Store, error) {
	err := repository.database.Create(&store).Error
	if err != nil {
		return entities.Store{}, err
	}
	return store, nil
}

func (repository *storeRepositoryImpl) Delete(id uint) (bool, error) {
	err := repository.database.Delete(&entities.Store{}, id).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

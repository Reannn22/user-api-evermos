package repositories

import (
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
}

type storeRepositoryImpl struct {
	database *gorm.DB
}

func NewStoreRepository(database *gorm.DB) StoreRepository {
	return &storeRepositoryImpl{database}
}

func (repository *storeRepositoryImpl) FindAllPagination(pagination responder.Pagination) (responder.Pagination, error) {
	var stores []entities.Store

	keyword := "%" + pagination.Keyword + "%"

	where_value := func(keyword string) *gorm.DB {
		if keyword != "" {
			return repository.database.Where("nama_toko LIKE ?", keyword)
		}
		return repository.database
	}

	err := where_value(keyword).
		Scopes(responder.PaginationFormat(keyword, stores, &pagination, where_value(keyword))).
		Find(&stores).Error

	if err != nil {
		return pagination, err
	}

	storesFormatter := []models.StoreResponse{}

	for _, store := range stores {
		storeFormatter := models.StoreResponse{}
		storeFormatter.ID = store.ID
		storeFormatter.NamaToko = store.NamaToko
		storeFormatter.UrlFoto = store.UrlFoto

		storesFormatter = append(storesFormatter, storeFormatter)
	}

	pagination.Data = storesFormatter

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

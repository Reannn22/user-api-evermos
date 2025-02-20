package services

import (
	"errors"
	"mini-project-evermos/models"
	"mini-project-evermos/models/entities"
	"mini-project-evermos/models/responder"
	"mini-project-evermos/repositories"
	"time"
)

// Contract
type StoreService interface {
	GetAll(limit int, page int, keyword string) (responder.Pagination, error)
	GetByUserId(id uint) (models.StoreResponse, error)
	GetById(id uint, user_id uint) (models.StoreResponse, error)
	Create(input models.StoreProcess) (models.StoreResponse, error)
	Edit(input models.StoreProcess) (models.StoreResponse, error) // Changed return type
	Delete(id uint, user_id uint) (models.StoreResponse, error)   // Add this line
}

type storeServiceImpl struct {
	repository repositories.StoreRepository
}

func NewStoreService(storeRepository *repositories.StoreRepository) StoreService {
	return &storeServiceImpl{
		repository: *storeRepository,
	}
}

func (service *storeServiceImpl) GetAll(limit int, page int, keyword string) (responder.Pagination, error) {
	request := responder.Pagination{}
	request.Limit = limit
	request.Page = page
	request.Keyword = keyword

	//get all user
	response, err := service.repository.FindAllPagination(request)

	if err != nil {
		return responder.Pagination{}, err
	}
	return response, nil
}

func (service *storeServiceImpl) GetByUserId(user_id uint) (models.StoreResponse, error) {
	store, err := service.repository.FindByUserId(user_id)

	if err != nil {
		return models.StoreResponse{}, err
	}

	var response = models.StoreResponse{}
	response.ID = store.ID
	response.NamaToko = store.NamaToko
	response.UrlFoto = store.UrlFoto
	response.CreatedAt = store.CreatedAt
	response.UpdatedAt = store.UpdatedAt

	return response, nil
}

func (service *storeServiceImpl) GetById(id uint, user_id uint) (models.StoreResponse, error) {
	store, err := service.repository.FindById(id)

	if err != nil {
		return models.StoreResponse{}, err
	}

	// Remove forbidden check since store details should be public
	// if store.IDUser != user_id {
	//     return models.StoreResponse{}, errors.New("forbidden")
	// }

	var response = models.StoreResponse{}
	response.ID = store.ID
	response.NamaToko = store.NamaToko
	response.UrlFoto = store.UrlFoto
	response.CreatedAt = store.CreatedAt
	response.UpdatedAt = store.UpdatedAt

	return response, nil
}

func (service *storeServiceImpl) Create(input models.StoreProcess) (models.StoreResponse, error) {
	store := entities.Store{
		IDUser:   input.UserID,
		NamaToko: input.NamaToko,
		UrlFoto:  &input.URL,
	}

	created_store, err := service.repository.Insert(store)
	if err != nil {
		return models.StoreResponse{}, err
	}

	response := models.StoreResponse{
		ID:        created_store.ID,
		NamaToko:  created_store.NamaToko,
		UrlFoto:   created_store.UrlFoto,
		CreatedAt: created_store.CreatedAt,
		UpdatedAt: created_store.UpdatedAt,
	}

	return response, nil
}

func (service *storeServiceImpl) Edit(input models.StoreProcess) (models.StoreResponse, error) {
	store, err := service.repository.FindById(input.ID)
	if err != nil {
		return models.StoreResponse{}, err
	}

	if store.IDUser != input.UserID {
		return models.StoreResponse{}, errors.New("forbidden")
	}

	date_now := time.Now()
	string_date := date_now.Format("2006_01_02_15_04_05")
	filename := string_date + "-" + input.URL

	req := entities.Store{}
	req.NamaToko = input.NamaToko
	req.UrlFoto = &filename

	success, err := service.repository.Update(input.ID, req)
	if err != nil || !success {
		return models.StoreResponse{}, err
	}

	// Fetch updated store
	updated_store, err := service.repository.FindById(input.ID)
	if err != nil {
		return models.StoreResponse{}, err
	}

	response := models.StoreResponse{
		ID:        updated_store.ID,
		NamaToko:  updated_store.NamaToko,
		UrlFoto:   updated_store.UrlFoto,
		CreatedAt: updated_store.CreatedAt,
		UpdatedAt: updated_store.UpdatedAt,
	}

	return response, nil
}

func (service *storeServiceImpl) Delete(id uint, user_id uint) (models.StoreResponse, error) {
	store, err := service.repository.FindById(id)
	if err != nil {
		return models.StoreResponse{}, err
	}

	if store.IDUser != user_id {
		return models.StoreResponse{}, errors.New("forbidden")
	}

	success, err := service.repository.Delete(id)
	if err != nil || !success {
		return models.StoreResponse{}, err
	}

	response := models.StoreResponse{
		ID:        store.ID,
		NamaToko:  store.NamaToko,
		UrlFoto:   store.UrlFoto,
		CreatedAt: store.CreatedAt,
		UpdatedAt: store.UpdatedAt,
	}

	return response, nil
}

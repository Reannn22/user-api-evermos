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
	Edit(input models.StoreProcess) (string, error)
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

	return response, nil
}

func (service *storeServiceImpl) GetById(id uint, user_id uint) (models.StoreResponse, error) {
	store, err := service.repository.FindById(id)

	if err != nil {
		return models.StoreResponse{}, err
	}

	if store.IDUser != user_id {
		return models.StoreResponse{}, errors.New("forbidden")
	}

	var response = models.StoreResponse{}
	response.ID = store.ID
	response.NamaToko = store.NamaToko
	response.UrlFoto = store.UrlFoto

	return response, nil
}

func (service *storeServiceImpl) Edit(input models.StoreProcess) (string, error) {
	store, err := service.repository.FindById(input.ID)

	if err != nil {
		return "", err
	}

	if store.IDUser != input.UserID {
		return "", errors.New("forbidden")
	}

	date_now := time.Now()
	string_date := date_now.Format("2006_01_02_15_04_05")

	filename := string_date + "-" + input.URL

	req := entities.Store{}
	req.NamaToko = input.NamaToko
	req.UrlFoto = &filename

	//update
	_, err = service.repository.Update(input.ID, req)

	if err != nil {
		return "", err
	}

	return filename, nil
}

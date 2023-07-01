package services

import (
	"mini-project-evermos/models"
	"mini-project-evermos/models/entities"
	"mini-project-evermos/repositories"
)

// Contract
type CategoryService interface {
	GetAll() ([]models.CategoryResponse, error)
	GetById(id uint) (models.CategoryResponse, error)
	Create(payload models.CategoryRequest) (bool, error)
	Edit(id uint, payload models.CategoryRequest) (bool, error)
	Delete(id uint) (bool, error)
}

type categoryServiceImpl struct {
	repository repositories.CategoryRepository
}

func NewCategoryService(categoryRepository *repositories.CategoryRepository) CategoryService {
	return &categoryServiceImpl{
		repository: *categoryRepository,
	}
}

func (service *categoryServiceImpl) GetAll() ([]models.CategoryResponse, error) {
	categories, err := service.repository.FindAll()

	if err != nil {
		return nil, err
	}

	// mapping response
	responses := []models.CategoryResponse{}

	for _, category := range categories {
		response := models.CategoryResponse{}
		response.ID = category.ID
		response.NamaCategory = category.NamaCategory

		responses = append(responses, response)
	}

	return responses, nil
}

func (service *categoryServiceImpl) GetById(id uint) (models.CategoryResponse, error) {
	category, err := service.repository.FindById(id)

	if err != nil {
		return models.CategoryResponse{}, err
	}

	var response = models.CategoryResponse{}
	response.ID = category.ID
	response.NamaCategory = category.NamaCategory

	return response, nil
}

func (service *categoryServiceImpl) Create(payload models.CategoryRequest) (bool, error) {
	category := entities.Category{}
	category.NamaCategory = payload.NamaCategory

	//create
	res, err := service.repository.Insert(category)

	return res, err
}

func (service *categoryServiceImpl) Edit(id uint, payload models.CategoryRequest) (bool, error) {
	//check
	_, err := service.repository.FindById(id)

	if err != nil {
		return false, err
	}

	category := entities.Category{}
	category.NamaCategory = payload.NamaCategory

	//update
	res, err := service.repository.Update(id, category)

	return res, err
}

func (service *categoryServiceImpl) Delete(id uint) (bool, error) {
	//check
	_, err := service.repository.FindById(id)

	if err != nil {
		return false, err
	}

	//delete role
	res, err := service.repository.Destroy(id)

	return res, err
}

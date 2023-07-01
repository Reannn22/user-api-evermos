package services

import (
	"errors"
	"mini-project-evermos/models"
	"mini-project-evermos/models/entities"
	"mini-project-evermos/repositories"
)

// Contract
type AddressService interface {
	GetAll(user_id uint) ([]models.AddressResponse, error)
	GetById(id uint, user_id uint) (models.AddressResponse, error)
	Create(payload models.AddressRequest, user_id uint) (bool, error)
	Edit(id uint, payload models.AddressRequest, user_id uint) (bool, error)
	Delete(id uint, user_id uint) (bool, error)
}

type addressServiceImpl struct {
	repository repositories.AddressRepository
}

func NewAddressService(addressRepository *repositories.AddressRepository) AddressService {
	return &addressServiceImpl{
		repository: *addressRepository,
	}
}

func (service *addressServiceImpl) GetAll(user_id uint) ([]models.AddressResponse, error) {
	addresses, err := service.repository.FindByUserId(user_id)

	if err != nil {
		return nil, err
	}

	// mapping response
	responses := []models.AddressResponse{}

	for _, address := range addresses {
		response := models.AddressResponse{}
		response.ID = address.ID
		response.JudulAlamat = address.JudulAlamat
		response.NamaPenerima = address.NamaPenerima
		response.NoTelp = address.NoTelp
		response.DetailAlamat = address.DetailAlamat

		responses = append(responses, response)
	}

	return responses, nil
}

func (service *addressServiceImpl) GetById(id uint, user_id uint) (models.AddressResponse, error) {
	address, err := service.repository.FindById(id)

	if err != nil {
		return models.AddressResponse{}, err
	}

	if address.IDUser != user_id {
		return models.AddressResponse{}, errors.New("forbidden")
	}

	var response = models.AddressResponse{}
	response.ID = address.ID
	response.JudulAlamat = address.JudulAlamat
	response.NamaPenerima = address.NamaPenerima
	response.NoTelp = address.NoTelp
	response.DetailAlamat = address.DetailAlamat

	return response, nil
}

func (service *addressServiceImpl) Create(payload models.AddressRequest, user_id uint) (bool, error) {
	address := entities.Address{}
	address.IDUser = user_id
	address.JudulAlamat = payload.JudulAlamat
	address.NamaPenerima = payload.NamaPenerima
	address.NoTelp = payload.NoTelp
	address.DetailAlamat = payload.DetailAlamat

	//create
	res, err := service.repository.Insert(address)

	return res, err
}

func (service *addressServiceImpl) Edit(id uint, payload models.AddressRequest, user_id uint) (bool, error) {
	//check
	check_address, err := service.repository.FindById(id)

	if err != nil {
		return false, err
	}

	if check_address.IDUser != user_id {
		return false, errors.New("forbidden")
	}

	address := entities.Address{}
	address.NamaPenerima = payload.NamaPenerima
	address.NoTelp = payload.NoTelp
	address.DetailAlamat = payload.DetailAlamat

	//update
	res, err := service.repository.Update(id, address)

	return res, err
}

func (service *addressServiceImpl) Delete(id uint, user_id uint) (bool, error) {
	//check
	check_address, err := service.repository.FindById(id)

	if err != nil {
		return false, err
	}

	if check_address.IDUser != user_id {
		return false, errors.New("forbidden")
	}

	//delete role
	res, err := service.repository.Destroy(id)

	return res, err
}

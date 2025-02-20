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
	Create(payload models.AddressRequest, user_id uint) (models.AddressResponse, error)
	Edit(id uint, payload models.AddressRequest, user_id uint) (models.AddressResponse, error)
	Delete(id uint, user_id uint) (models.AddressResponse, error)
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
		response := models.AddressResponse{
			ID:           address.ID,
			JudulAlamat:  address.JudulAlamat,
			NamaPenerima: address.NamaPenerima,
			NoTelp:       address.NoTelp,
			DetailAlamat: address.DetailAlamat,
			CreatedAt:    address.CreatedAt,
			UpdatedAt:    address.UpdatedAt,
		}
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

	var response = models.AddressResponse{
		ID:           address.ID,
		JudulAlamat:  address.JudulAlamat,
		NamaPenerima: address.NamaPenerima,
		NoTelp:       address.NoTelp,
		DetailAlamat: address.DetailAlamat,
		CreatedAt:    address.CreatedAt,
		UpdatedAt:    address.UpdatedAt,
	}

	return response, nil
}

func (service *addressServiceImpl) Create(payload models.AddressRequest, user_id uint) (models.AddressResponse, error) {
	address := entities.Address{
		IDUser:       user_id,
		JudulAlamat:  payload.JudulAlamat,
		NamaPenerima: payload.NamaPenerima,
		NoTelp:       payload.NoTelp,
		DetailAlamat: payload.DetailAlamat,
		IDProvinsi:   payload.IDProvinsi,
		IDKota:       payload.IDKota,
	}

	success, err := service.repository.Insert(address)
	if err != nil || !success {
		return models.AddressResponse{}, err
	}

	// Fetch the created address to get the complete data including timestamps
	created_address, err := service.repository.FindByCondition(map[string]interface{}{
		"id_user":       user_id,
		"judul_alamat":  payload.JudulAlamat,
		"nama_penerima": payload.NamaPenerima,
		"no_telp":       payload.NoTelp,
		"detail_alamat": payload.DetailAlamat,
	})
	if err != nil {
		return models.AddressResponse{}, err
	}

	response := models.AddressResponse{
		ID:           created_address.ID,
		JudulAlamat:  created_address.JudulAlamat,
		NamaPenerima: created_address.NamaPenerima,
		NoTelp:       created_address.NoTelp,
		DetailAlamat: created_address.DetailAlamat,
		CreatedAt:    created_address.CreatedAt,
		UpdatedAt:    created_address.UpdatedAt,
	}

	return response, nil
}

func (service *addressServiceImpl) Edit(id uint, payload models.AddressRequest, user_id uint) (models.AddressResponse, error) {
	//check
	check_address, err := service.repository.FindById(id)

	if err != nil {
		return models.AddressResponse{}, err
	}

	if check_address.IDUser != user_id {
		return models.AddressResponse{}, errors.New("forbidden")
	}

	address := entities.Address{
		JudulAlamat:  check_address.JudulAlamat, // Keep the original judul_alamat
		NamaPenerima: payload.NamaPenerima,
		NoTelp:       payload.NoTelp,
		DetailAlamat: payload.DetailAlamat,
		IDProvinsi:   payload.IDProvinsi,
		IDKota:       payload.IDKota,
	}

	//update
	success, err := service.repository.Update(id, address)
	if err != nil || !success {
		return models.AddressResponse{}, err
	}

	// Fetch the updated address to get the latest data including timestamps
	updated_address, err := service.repository.FindById(id)
	if err != nil {
		return models.AddressResponse{}, err
	}

	response := models.AddressResponse{
		ID:           updated_address.ID,
		JudulAlamat:  updated_address.JudulAlamat,
		NamaPenerima: updated_address.NamaPenerima,
		NoTelp:       updated_address.NoTelp,
		DetailAlamat: updated_address.DetailAlamat,
		CreatedAt:    updated_address.CreatedAt,
		UpdatedAt:    updated_address.UpdatedAt,
	}

	return response, nil
}

func (service *addressServiceImpl) Delete(id uint, user_id uint) (models.AddressResponse, error) {
	//check
	check_address, err := service.repository.FindById(id)

	if err != nil {
		return models.AddressResponse{}, err
	}

	if check_address.IDUser != user_id {
		return models.AddressResponse{}, errors.New("forbidden")
	}

	//delete role
	_, err = service.repository.Destroy(id)
	if err != nil {
		return models.AddressResponse{}, err
	}

	response := models.AddressResponse{
		ID:           check_address.ID,
		JudulAlamat:  check_address.JudulAlamat,
		NamaPenerima: check_address.NamaPenerima,
		NoTelp:       check_address.NoTelp,
		DetailAlamat: check_address.DetailAlamat,
		CreatedAt:    check_address.CreatedAt,
		UpdatedAt:    check_address.UpdatedAt,
	}

	return response, nil
}

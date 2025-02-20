package services

import (
	"mini-project-evermos/models"
	"mini-project-evermos/models/entities"
	"mini-project-evermos/repositories"
	"mini-project-evermos/utils/region"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Contract
type UserService interface {
	GetById(id uint) (models.UserResponse, error)
	Edit(id uint, payload models.UserRequest) (models.UserResponse, error) // Changed return type
	Delete(id uint) (models.UserResponse, error)                           // Change return type
}

type userServiceImpl struct {
	repository repositories.UserRepository
}

func NewUserService(userRepository *repositories.UserRepository) UserService {
	return &userServiceImpl{
		repository: *userRepository,
	}
}

func (service *userServiceImpl) GetById(id uint) (models.UserResponse, error) {
	user, err := service.repository.FindById(id)
	if err != nil {
		return models.UserResponse{}, err
	}

	// Get region data
	province, _ := region.GetProvinceByID(user.IDProvinsi)
	city, _ := region.GetCityByID(user.IDKota)

	response := models.UserResponse{
		ID:           user.ID,
		Nama:         user.Nama,
		KataSandi:    user.KataSandi,
		Notelp:       user.Notelp,
		TanggalLahir: user.TanggalLahir,
		JenisKelamin: user.JenisKelamin,
		Tentang:      user.Tentang,
		Pekerjaan:    user.Pekerjaan,
		Email:        user.Email,
		IDProvinsi: models.ProvinceDetail{
			ID:   user.IDProvinsi,
			Name: province.Name,
		},
		IDKota: models.CityDetail{
			ID:         user.IDKota,
			ProvinceID: user.IDProvinsi,
			Name:       city.Name,
		},
		IsAdmin:   user.IsAdmin,
		CreatedAt: *user.CreatedAt,
		UpdatedAt: *user.UpdatedAt,
	}

	return response, nil
}

func (service *userServiceImpl) Edit(id uint, payload models.UserRequest) (models.UserResponse, error) {
	//check if user exists
	_, err := service.repository.FindById(id)
	if err != nil {
		return models.UserResponse{}, err
	}

	//encrypt pass
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(payload.KataSandi), bcrypt.MinCost)
	if err != nil {
		return models.UserResponse{}, err
	}

	//string to date
	date, err := time.Parse("02/01/2006", payload.TanggalLahir)
	if err != nil {
		return models.UserResponse{}, err
	}

	//mapping
	user := entities.User{}
	user.Nama = payload.Nama
	user.Notelp = payload.NoTelp
	user.Email = payload.Email
	user.KataSandi = string(passwordHash)
	user.TanggalLahir = date
	user.JenisKelamin = payload.JenisKelamin
	user.Tentang = &payload.Tentang
	user.Pekerjaan = payload.Pekerjaan
	user.IDProvinsi = payload.IDProvinsi
	user.IDKota = payload.IDKota
	user.IsAdmin = payload.IsAdmin

	//update
	_, err = service.repository.Update(id, user)
	if err != nil {
		return models.UserResponse{}, err
	}

	// Get updated user data
	updatedUser, err := service.repository.FindById(id)
	if err != nil {
		return models.UserResponse{}, err
	}

	// Get region data
	province, _ := region.GetProvinceByID(updatedUser.IDProvinsi)
	city, _ := region.GetCityByID(updatedUser.IDKota)

	// Map to response
	response := models.UserResponse{
		ID:           updatedUser.ID,
		Nama:         updatedUser.Nama,
		KataSandi:    updatedUser.KataSandi,
		Notelp:       updatedUser.Notelp,
		TanggalLahir: updatedUser.TanggalLahir,
		JenisKelamin: updatedUser.JenisKelamin,
		Tentang:      updatedUser.Tentang,
		Pekerjaan:    updatedUser.Pekerjaan,
		Email:        updatedUser.Email,
		IDProvinsi: models.ProvinceDetail{
			ID:   updatedUser.IDProvinsi,
			Name: province.Name,
		},
		IDKota: models.CityDetail{
			ID:         updatedUser.IDKota,
			ProvinceID: updatedUser.IDProvinsi,
			Name:       city.Name,
		},
		IsAdmin:   updatedUser.IsAdmin,
		CreatedAt: *updatedUser.CreatedAt,
		UpdatedAt: *updatedUser.UpdatedAt,
	}

	return response, nil
}

func (service *userServiceImpl) Delete(id uint) (models.UserResponse, error) {
	// Get user data before deletion
	user, err := service.repository.FindById(id)
	if err != nil {
		return models.UserResponse{}, err
	}

	// Get region data
	province, _ := region.GetProvinceByID(user.IDProvinsi)
	city, _ := region.GetCityByID(user.IDKota)

	// Create response before deleting
	response := models.UserResponse{
		ID:           user.ID,
		Nama:         user.Nama,
		KataSandi:    user.KataSandi,
		Notelp:       user.Notelp,
		TanggalLahir: user.TanggalLahir,
		JenisKelamin: user.JenisKelamin,
		Tentang:      user.Tentang,
		Pekerjaan:    user.Pekerjaan,
		Email:        user.Email,
		IDProvinsi: models.ProvinceDetail{
			ID:   user.IDProvinsi,
			Name: province.Name,
		},
		IDKota: models.CityDetail{
			ID:         user.IDKota,
			ProvinceID: user.IDProvinsi,
			Name:       city.Name,
		},
		IsAdmin:   user.IsAdmin,
		CreatedAt: *user.CreatedAt,
		UpdatedAt: *user.UpdatedAt,
	}

	// Delete the user
	err = service.repository.Delete(id)
	if err != nil {
		return models.UserResponse{}, err
	}

	return response, nil
}

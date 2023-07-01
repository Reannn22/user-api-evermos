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
	Edit(id uint, payload models.UserRequest) (bool, error)
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
	check_user, err := service.repository.FindById(id)

	if err != nil {
		return models.UserResponse{}, err
	}

	//get region
	province, err := region.GetProvinceByID(check_user.IDProvinsi)
	city, err := region.GetCityByID(check_user.IDKota)

	var response = models.UserResponse{}
	response.Nama = check_user.Nama
	response.NoTelp = check_user.Notelp
	response.Tentang = check_user.Tentang
	response.Pekerjaan = check_user.Pekerjaan
	response.TanggalLahir = check_user.TanggalLahir.Format("02/01/2006")
	response.Email = check_user.Email
	response.IDProvinsi = province
	response.IDKota = city

	return response, nil
}

func (service *userServiceImpl) Edit(id uint, payload models.UserRequest) (bool, error) {
	//check
	check_user, err := service.repository.FindById(id)

	if err != nil {
		return false, err
	}

	//encrypt pass
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(payload.KataSandi), bcrypt.MinCost)
	if err != nil {
		return false, err
	}

	//string to date
	date, err := time.Parse("02/01/2006", payload.TanggalLahir)

	if err != nil {
		return false, err
	}

	//mapping
	user := entities.User{}
	user.Nama = payload.Nama
	if check_user.Notelp != payload.NoTelp {
		user.Notelp = payload.NoTelp
	}
	if check_user.Email != payload.Email {
		user.Email = payload.Email
	}
	user.KataSandi = string(passwordHash)
	user.TanggalLahir = date
	user.Pekerjaan = payload.Pekerjaan
	user.IDProvinsi = payload.IDProvinsi
	user.IDKota = payload.IDKota

	//update
	res, err := service.repository.Update(id, user)

	return res, err
}

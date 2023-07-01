package services

import (
	"errors"
	"mini-project-evermos/models"
	"mini-project-evermos/repositories"
	"mini-project-evermos/utils/jwt"
	"mini-project-evermos/utils/region"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Contract
type AuthService interface {
	Register(input models.RegisterRequest) error
	Login(input models.LoginRequest) (models.LoginResponse, error)
}

type authServiceImpl struct {
	repository     repositories.AuthRepository
	repositoryUser repositories.UserRepository
}

func NewAuthService(authRepository *repositories.AuthRepository, userRepository *repositories.UserRepository) AuthService {
	return &authServiceImpl{
		repository:     *authRepository,
		repositoryUser: *userRepository,
	}
}

func (service *authServiceImpl) Register(input models.RegisterRequest) error {

	//encrypt pass
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.KataSandi), bcrypt.MinCost)
	if err != nil {
		return err
	}

	//string to date
	date, err := time.Parse("02/01/2006", input.TanggalLahir)

	if err != nil {
		return err
	}

	//mapping
	user := models.RegisterProcess{}
	user.Nama = input.Nama
	user.NoTelp = input.NoTelp
	user.Email = input.Email
	user.KataSandi = string(passwordHash)
	user.TanggalLahir = date
	user.Pekerjaan = input.Pekerjaan
	user.IDProvinsi = input.IDProvinsi
	user.IDKota = input.IDKota

	//register user
	err = service.repository.Register(user)

	if err != nil {
		return err
	}

	return nil
}

func (service *authServiceImpl) Login(input models.LoginRequest) (models.LoginResponse, error) {
	no_telp := input.NoTelp
	password := input.KataSandi

	//check user
	check_user, err := service.repositoryUser.FindByNoTelp(no_telp)

	if err != nil {
		return models.LoginResponse{}, errors.New("No Telp Not Found")
	}

	//check login
	err = bcrypt.CompareHashAndPassword([]byte(check_user.KataSandi), []byte(password))

	if err != nil {
		return models.LoginResponse{}, errors.New("No Telp atau kata sandi salah")
	}

	//generate token jwt
	token, err := jwt.GenerateNewAccessToken(check_user)

	//get region
	province, err := region.GetProvinceByID(check_user.IDProvinsi)
	city, err := region.GetCityByID(check_user.IDKota)

	//response mapping
	var response = models.LoginResponse{}
	response.Nama = check_user.Nama
	response.NoTelp = check_user.Notelp
	response.Tentang = check_user.Tentang
	response.Pekerjaan = check_user.Pekerjaan
	response.TanggalLahir = check_user.TanggalLahir.Format("02/01/2006")
	response.Email = check_user.Email
	response.IDProvinsi = province
	response.IDKota = city
	response.Token = token

	return response, nil
}

package services

import (
	"errors"
	"mini-project-evermos/models"
	"mini-project-evermos/models/responder"
	"mini-project-evermos/repositories"
	"os"
)

// Contract
type ProductService interface {
	GetAll(limit int, page int, keyword string) (responder.Pagination, error)
	GetById(id uint, user_id uint) (models.ProductResponse, error)
	Create(input models.ProductRequest, user_id uint) (bool, error)
	Update(input models.ProductRequest, id uint, user_id uint) (bool, error)
	Delete(id uint, user_id uint) (bool, error)
}

type productServiceImpl struct {
	repository               repositories.ProductRepository
	repositoryProductPicture repositories.ProductPictureRepository
	repositoryStore          repositories.StoreRepository
}

func NewProductService(productRepository *repositories.ProductRepository, storeRepository *repositories.StoreRepository, productPictureRepository *repositories.ProductPictureRepository) ProductService {
	return &productServiceImpl{
		repository:               *productRepository,
		repositoryProductPicture: *productPictureRepository,
		repositoryStore:          *storeRepository,
	}
}

func (service *productServiceImpl) GetAll(limit int, page int, keyword string) (responder.Pagination, error) {
	request := responder.Pagination{}
	request.Limit = limit
	request.Page = page
	request.Keyword = keyword

	response, err := service.repository.FindAllPagination(request)

	if err != nil {
		return responder.Pagination{}, err
	}
	return response, nil
}

func (service *productServiceImpl) GetById(id uint, user_id uint) (models.ProductResponse, error) {
	product, err := service.repository.FindById(id)

	if err != nil {
		return models.ProductResponse{}, err
	}

	if product.Store.IDUser != user_id {
		return models.ProductResponse{}, errors.New("forbidden")
	}

	var response = models.ProductResponse{}
	response.ID = product.ID
	response.NamaProduk = product.NamaProduk
	response.Slug = product.Slug
	response.HargaReseller = product.HargaReseller
	response.HargaKonsumen = product.HargaKonsumen
	response.Stok = product.Stok
	response.Deskripsi = product.Deskripsi
	response.Store.ID = product.Store.ID
	response.Store.NamaToko = product.Store.NamaToko
	response.Store.UrlFoto = product.Store.UrlFoto
	response.Category.ID = product.Category.ID
	response.Category.NamaCategory = product.Category.NamaCategory

	picturesFormatter := []models.ProductPictureResponse{}

	for _, picture := range product.ProductPicture {
		pictureFormatter := models.ProductPictureResponse{}
		pictureFormatter.ID = picture.ID
		pictureFormatter.IDProduk = picture.IDProduk
		pictureFormatter.Url = picture.Url

		picturesFormatter = append(picturesFormatter, pictureFormatter)
	}
	response.Photos = picturesFormatter

	return response, nil
}

func (service *productServiceImpl) Create(input models.ProductRequest, user_id uint) (bool, error) {
	store, err := service.repositoryStore.FindByUserId(user_id)

	if err != nil {
		return false, err
	}
	input.StoreID = store.ID

	save, err := service.repository.Insert(input)

	if err != nil {
		for _, v := range input.Photos {
			os.Remove("uploads/" + v)
		}
		return false, err
	}

	return save, nil
}

func (service *productServiceImpl) Update(input models.ProductRequest, id uint, user_id uint) (bool, error) {
	product, err := service.repository.FindById(id)

	if err != nil {
		return false, err
	}

	if product.Store.IDUser != user_id {
		for _, v := range input.Photos {
			os.Remove("uploads/" + v)
		}
		return false, errors.New("forbidden")
	}

	picture, err := service.repositoryProductPicture.FindByProductId(id)

	update, err := service.repository.Update(input, id)

	if err != nil {
		for _, v := range input.Photos {
			os.Remove("uploads/" + v)
		}
		return false, err
	}

	for _, v := range picture {
		os.Remove("uploads/" + v.Url)
	}

	return update, nil
}

func (service *productServiceImpl) Delete(id uint, user_id uint) (bool, error) {
	//check
	product, err := service.repository.FindById(id)

	if err != nil {
		return false, err
	}

	if product.Store.IDUser != user_id {
		return false, errors.New("forbidden")
	}

	//delete role
	res, err := service.repository.Destroy(id)

	return res, err
}

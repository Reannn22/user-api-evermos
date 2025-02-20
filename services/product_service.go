package services

import (
	"errors"
	"fmt"
	"mini-project-evermos/models"
	"mini-project-evermos/models/responder"
	"mini-project-evermos/repositories"
	"os"
)

type ProductService interface {
	GetAll(limit int, page int, keyword string) (responder.Pagination, error)
	GetById(id uint, user_id uint) (models.ProductResponse, error)
	Create(input models.ProductRequest, user_id uint) (models.ProductResponse, error)
	Update(input models.ProductRequest, id uint, user_id uint) (models.ProductResponse, error)
	Delete(id uint, user_id uint) (models.ProductResponse, error)
}

type productServiceImpl struct {
	repository               repositories.ProductRepository
	repositoryProductPicture repositories.ProductPictureRepository
	repositoryStore          repositories.StoreRepository
	repositoryCategory       repositories.CategoryRepository
}

func NewProductService(
	productRepository *repositories.ProductRepository,
	storeRepository *repositories.StoreRepository,
	productPictureRepository *repositories.ProductPictureRepository,
	categoryRepository *repositories.CategoryRepository,
) ProductService {
	return &productServiceImpl{
		repository:               *productRepository,
		repositoryProductPicture: *productPictureRepository,
		repositoryStore:          *storeRepository,
		repositoryCategory:       *categoryRepository,
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
	response.CreatedAt = product.CreatedAt
	response.UpdatedAt = product.UpdatedAt

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

func (service *productServiceImpl) Create(input models.ProductRequest, user_id uint) (models.ProductResponse, error) {
	store, err := service.repositoryStore.FindByUserId(user_id)
	if err != nil {
		return models.ProductResponse{}, err
	}

	category, err := service.repositoryCategory.FindById(input.CategoryID)
	if err != nil {
		return models.ProductResponse{}, err
	}
	if category.ID == 0 {
		return models.ProductResponse{}, errors.New("category with ID " + fmt.Sprint(input.CategoryID) + " not found")
	}

	input.StoreID = store.ID

	product, err := service.repository.Insert(input)
	if err != nil {
		for _, v := range input.PhotoURLs {
			os.Remove("uploads/" + v)
		}
		return models.ProductResponse{}, err
	}

	// Get the complete product data to return
	return service.GetById(product.ID, user_id)
}

func (service *productServiceImpl) Update(input models.ProductRequest, id uint, user_id uint) (models.ProductResponse, error) {
	product, err := service.repository.FindById(id)
	if err != nil {
		return models.ProductResponse{}, err
	}

	if product.Store.IDUser != user_id {
		for _, v := range input.PhotoURLs {
			os.Remove("uploads/" + v)
		}
		return models.ProductResponse{}, errors.New("forbidden")
	}

	picture, err := service.repositoryProductPicture.FindByProductId(id)
	_, err = service.repository.Update(input, id)

	if err != nil {
		for _, v := range input.PhotoURLs {
			os.Remove("uploads/" + v)
		}
		return models.ProductResponse{}, err
	}

	for _, v := range picture {
		os.Remove("uploads/" + v.Url)
	}

	// Get the updated product data to return
	return service.GetById(id, user_id)
}

func (service *productServiceImpl) Delete(id uint, user_id uint) (models.ProductResponse, error) {
	// Get the product data before deletion
	product, err := service.GetById(id, user_id)
	if err != nil {
		return models.ProductResponse{}, err
	}

	// Perform deletion
	_, err = service.repository.Destroy(id)
	if err != nil {
		return models.ProductResponse{}, err
	}

	return product, nil
}

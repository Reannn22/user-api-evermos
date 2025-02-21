package services

import (
	"mini-project-evermos/models"
	"mini-project-evermos/repositories"
)

type ProductLogService interface {
	Create(input models.ProductLogProcess) (models.ProductLogResponse, error)
	GetAll() ([]models.ProductLogResponse, error)
	GetById(id uint) (models.ProductLogResponse, error)
	Update(id uint, input models.ProductLogProcess) (models.ProductLogResponse, error)
	Delete(id uint) (models.ProductLogResponse, error)
}

type productLogServiceImpl struct {
	repository repositories.ProductLogRepository
}

func NewProductLogService(productLogRepository *repositories.ProductLogRepository) ProductLogService {
	return &productLogServiceImpl{
		repository: *productLogRepository,
	}
}

func (service *productLogServiceImpl) Create(input models.ProductLogProcess) (models.ProductLogResponse, error) {
	productLog, err := service.repository.Insert(input)
	if err != nil {
		return models.ProductLogResponse{}, err
	}

	response := models.ProductLogResponse{
		ID:            productLog.ID,
		ProductID:     productLog.IDProduk,
		NamaProduk:    productLog.NamaProduk,
		Slug:          productLog.Slug,
		HargaReseller: productLog.HargaReseller,
		HargaKonsumen: productLog.HargaKonsumen,
		StoreID:       productLog.IDToko,
		CategoryID:    productLog.IDCategory,
		CreatedAt:     *productLog.CreatedAt,
		UpdatedAt:     *productLog.UpdatedAt,
	}

	return response, nil
}

func (service *productLogServiceImpl) GetAll() ([]models.ProductLogResponse, error) {
	productLogs, err := service.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []models.ProductLogResponse
	for _, log := range productLogs {
		response := models.ProductLogResponse{
			ProductID:     log.IDProduk,
			NamaProduk:    log.NamaProduk,
			Slug:          log.Slug,
			HargaReseller: log.HargaReseller,
			HargaKonsumen: log.HargaKonsumen,
			Deskripsi:     *log.Deskripsi,
			StoreID:       log.IDToko,
			CategoryID:    log.IDCategory,
			CreatedAt:     *log.CreatedAt,
			UpdatedAt:     *log.UpdatedAt,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (service *productLogServiceImpl) GetById(id uint) (models.ProductLogResponse, error) {
	productLog, err := service.repository.FindById(id)
	if err != nil {
		return models.ProductLogResponse{}, err
	}

	response := models.ProductLogResponse{
		ID:            productLog.ID,
		ProductID:     productLog.IDProduk,
		NamaProduk:    productLog.NamaProduk,
		Slug:          productLog.Slug,
		HargaReseller: productLog.HargaReseller,
		HargaKonsumen: productLog.HargaKonsumen,
		Deskripsi:     *productLog.Deskripsi,
		StoreID:       productLog.IDToko,
		CategoryID:    productLog.IDCategory,
		CreatedAt:     *productLog.CreatedAt,
		UpdatedAt:     *productLog.UpdatedAt,
	}

	return response, nil
}

func (service *productLogServiceImpl) Update(id uint, input models.ProductLogProcess) (models.ProductLogResponse, error) {
	productLog, err := service.repository.Update(id, input)
	if err != nil {
		return models.ProductLogResponse{}, err
	}

	response := models.ProductLogResponse{
		ProductID:     productLog.IDProduk,
		NamaProduk:    productLog.NamaProduk,
		Slug:          productLog.Slug,
		HargaReseller: productLog.HargaReseller,
		HargaKonsumen: productLog.HargaKonsumen,
		Deskripsi:     *productLog.Deskripsi,
		StoreID:       productLog.IDToko,
		CategoryID:    productLog.IDCategory,
		CreatedAt:     *productLog.CreatedAt,
		UpdatedAt:     *productLog.UpdatedAt,
	}

	return response, nil
}

func (service *productLogServiceImpl) Delete(id uint) (models.ProductLogResponse, error) {
	// Get the product log before deletion
	productLog, err := service.repository.FindById(id)
	if err != nil {
		return models.ProductLogResponse{}, err
	}

	// Create response before deleting
	response := models.ProductLogResponse{
		ProductID:     productLog.IDProduk,
		NamaProduk:    productLog.NamaProduk,
		Slug:          productLog.Slug,
		HargaReseller: productLog.HargaReseller,
		HargaKonsumen: productLog.HargaKonsumen,
		Deskripsi:     *productLog.Deskripsi,
		StoreID:       productLog.IDToko,
		CategoryID:    productLog.IDCategory,
		CreatedAt:     *productLog.CreatedAt,
		UpdatedAt:     *productLog.UpdatedAt,
	}

	// Perform deletion
	err = service.repository.Delete(id)
	if err != nil {
		return models.ProductLogResponse{}, err
	}

	return response, nil
}

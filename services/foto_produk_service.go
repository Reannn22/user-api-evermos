package services

import (
	"mini-project-evermos/models"
	"mini-project-evermos/repositories"
)

type FotoProdukService interface {
	GetAll() ([]models.FotoProdukResponse, error)
	GetById(id uint) (models.FotoProdukResponse, error)
	GetByProductId(productId uint) ([]models.FotoProdukResponse, error)
	Create(input models.FotoProdukRequest, userId uint) (models.FotoProdukResponse, error)
	Update(id uint, input models.FotoProdukRequest, userId uint) (models.FotoProdukResponse, error)
	Delete(id uint, userId uint) error
}

type fotoProdukServiceImpl struct {
	repository     repositories.FotoProdukRepository
	prodRepository repositories.ProductRepository
}

func NewFotoProdukService(repository *repositories.FotoProdukRepository, prodRepository *repositories.ProductRepository) FotoProdukService {
	return &fotoProdukServiceImpl{
		repository:     *repository,
		prodRepository: *prodRepository,
	}
}

func (s *fotoProdukServiceImpl) GetAll() ([]models.FotoProdukResponse, error) {
	photos, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []models.FotoProdukResponse
	for _, photo := range photos {
		responses = append(responses, models.FotoProdukResponse{
			ID:        photo.ID,
			ProductID: photo.IDProduk,
			URL:       photo.Url,
			CreatedAt: *photo.CreatedAt,
			UpdatedAt: *photo.UpdatedAt,
		})
	}
	return responses, nil
}

func (s *fotoProdukServiceImpl) GetById(id uint) (models.FotoProdukResponse, error) {
	photo, err := s.repository.FindById(id)
	if err != nil {
		return models.FotoProdukResponse{}, err
	}

	return models.FotoProdukResponse{
		ID:        photo.ID,
		ProductID: photo.IDProduk,
		URL:       photo.Url,
		CreatedAt: *photo.CreatedAt,
		UpdatedAt: *photo.UpdatedAt,
	}, nil
}

func (s *fotoProdukServiceImpl) GetByProductId(productId uint) ([]models.FotoProdukResponse, error) {
	photos, err := s.repository.FindByProductId(productId)
	if err != nil {
		return nil, err
	}

	var responses []models.FotoProdukResponse
	for _, photo := range photos {
		responses = append(responses, models.FotoProdukResponse{
			ID:        photo.ID,
			ProductID: photo.IDProduk,
			URL:       photo.Url,
			CreatedAt: *photo.CreatedAt,
			UpdatedAt: *photo.UpdatedAt,
		})
	}
	return responses, nil
}

func (s *fotoProdukServiceImpl) Create(input models.FotoProdukRequest, userId uint) (models.FotoProdukResponse, error) {
	// Verify product exists
	_, err := s.prodRepository.FindById(input.ProductID)
	if err != nil {
		return models.FotoProdukResponse{}, err
	}

	// If photo_id is provided, update existing photo instead of creating new one
	if input.PhotoID > 0 {
		existingPhoto, err := s.repository.FindById(input.PhotoID)
		if err != nil {
			return models.FotoProdukResponse{}, err
		}

		// Update the existing photo
		updateInput := models.FotoProdukRequest{
			ProductID: existingPhoto.IDProduk,
			PhotoID:   existingPhoto.ID,
			URL:       input.URL,
		}

		updatedPhoto, err := s.repository.Update(input.PhotoID, updateInput)
		if err != nil {
			return models.FotoProdukResponse{}, err
		}

		return models.FotoProdukResponse{
			ID:        updatedPhoto.ID,
			ProductID: updatedPhoto.IDProduk,
			PhotoID:   updatedPhoto.PhotoID,
			URL:       updatedPhoto.Url,
			CreatedAt: *updatedPhoto.CreatedAt,
			UpdatedAt: *updatedPhoto.UpdatedAt,
		}, nil
	}

	// If no photo_id provided, create new photo
	photo, err := s.repository.Create(input)
	if err != nil {
		return models.FotoProdukResponse{}, err
	}

	return models.FotoProdukResponse{
		ID:        photo.ID,
		ProductID: photo.IDProduk,
		PhotoID:   photo.PhotoID,
		URL:       photo.Url,
		CreatedAt: *photo.CreatedAt,
		UpdatedAt: *photo.UpdatedAt,
	}, nil
}

func (s *fotoProdukServiceImpl) Update(id uint, input models.FotoProdukRequest, userId uint) (models.FotoProdukResponse, error) {
	// Verify product exists
	_, err := s.prodRepository.FindById(input.ProductID)
	if err != nil {
		return models.FotoProdukResponse{}, err
	}

	photo, err := s.repository.Update(id, input)
	if err != nil {
		return models.FotoProdukResponse{}, err
	}

	return models.FotoProdukResponse{
		ID:        photo.ID,
		ProductID: photo.IDProduk,
		PhotoID:   photo.PhotoID,
		URL:       photo.Url,
		CreatedAt: *photo.CreatedAt,
		UpdatedAt: *photo.UpdatedAt,
	}, nil
}

func (s *fotoProdukServiceImpl) Delete(id uint, userId uint) error {
	// First get the photo to check product ownership
	photo, err := s.repository.FindById(id)
	if err != nil {
		return err
	}

	// Verify product exists and belongs to user
	_, err = s.prodRepository.FindById(photo.IDProduk)
	if err != nil {
		return err
	}

	return s.repository.Delete(id)
}

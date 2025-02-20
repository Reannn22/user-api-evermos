package services

import (
	"errors"
	"fmt"
	"mini-project-evermos/models"
	"mini-project-evermos/models/responder"
	"mini-project-evermos/repositories"
	"strconv"
	"time"
)

type TransactionService interface {
	GetAll(limit int, page int, keyword string) (responder.Pagination, error)
	GetById(id uint, user_id uint) (models.TransactionResponse, error)
	Create(input models.TransactionRequest, user_id uint) (models.TransactionResponse, error)
	Update(id uint, user_id uint, input models.TransactionUpdateRequest) (models.TransactionResponse, error)
	Delete(id uint, user_id uint) error
}

type transactionServiceImpl struct {
	repository        repositories.TransactionRepository
	repositoryProduct repositories.ProductRepository
	repositoryAddress repositories.AddressRepository
}

func NewTransactionService(transactionRepository *repositories.TransactionRepository, productRepository *repositories.ProductRepository, addressRepository *repositories.AddressRepository) TransactionService {
	return &transactionServiceImpl{
		repository:        *transactionRepository,
		repositoryProduct: *productRepository,
		repositoryAddress: *addressRepository,
	}
}

func (service *transactionServiceImpl) Create(input models.TransactionRequest, user_id uint) (models.TransactionResponse, error) {
	// Check if address exists and belongs to user
	check_address, err := service.repositoryAddress.FindById(input.AlamatPengiriman)
	if err != nil {
		return models.TransactionResponse{}, fmt.Errorf("address error: %v", err)
	}

	if check_address.IDUser != user_id {
		return models.TransactionResponse{}, errors.New("forbidden: address does not belong to user")
	}

	date_now := time.Now()
	string_date := date_now.Format("2006-01-02-15-04-05")
	invoice := "INV" + string_date

	// Process product details and calculate total
	productLogsFormatter := []models.ProductLogProcess{}
	total := 0
	for _, detail := range input.DetailTrx {
		product, err := service.repositoryProduct.FindById(detail.ProductID)
		if err != nil {
			return models.TransactionResponse{}, err
		}

		stok, _ := strconv.Atoi(product.HargaKonsumen)
		total_detail := stok * detail.Kuantitas

		productLogFormatter := models.ProductLogProcess{
			ProductID:     product.ID,
			NamaProduk:    product.NamaProduk,
			Slug:          product.Slug,
			HargaReseller: product.HargaReseller,
			HargaKonsumen: product.HargaKonsumen,
			Stok:          product.Stok,
			Deskripsi:     *product.Deskripsi,
			CategoryID:    product.Category.ID,
			StoreID:       product.Store.ID,
			Kuantitas:     detail.Kuantitas,
			HargaTotal:    total_detail,
		}

		total += total_detail
		productLogsFormatter = append(productLogsFormatter, productLogFormatter)
	}

	transaction_data := models.TransactionProcessData{
		Transaction: models.TransactionProcess{
			MethodBayar:      input.MethodBayar,
			KodeInvoice:      invoice,
			AlamatPengiriman: input.AlamatPengiriman,
			UserID:           user_id,
			HargaTotal:       total,
		},
		LogProduct: productLogsFormatter,
	}

	trxID, err := service.repository.Insert(transaction_data)
	if err != nil {
		return models.TransactionResponse{}, err
	}

	return service.GetById(trxID, user_id)
}

func (service *transactionServiceImpl) GetById(id uint, user_id uint) (models.TransactionResponse, error) {
	transaction, err := service.repository.FindById(id)
	if err != nil {
		return models.TransactionResponse{}, err
	}

	if transaction.Address.IDUser != user_id {
		return models.TransactionResponse{}, errors.New("forbidden")
	}

	response := models.TransactionResponse{
		ID:          transaction.ID,
		HargaTotal:  transaction.HargaTotal,
		KodeInvoice: transaction.KodeInvoice,
		MethodBayar: transaction.MethodBayar,
		CreatedAt:   transaction.CreatedAt,
		UpdatedAt:   transaction.UpdatedAt,
		Address: models.AddressResponse{
			ID:           transaction.Address.ID,
			JudulAlamat:  transaction.Address.JudulAlamat,
			NamaPenerima: transaction.Address.NamaPenerima,
			NoTelp:       transaction.Address.NoTelp,
			DetailAlamat: transaction.Address.DetailAlamat,
			CreatedAt:    transaction.Address.CreatedAt,
			UpdatedAt:    transaction.Address.UpdatedAt,
		},
	}

	var details []models.TransactionDetailResponse
	for _, detail := range transaction.TrxDetail {
		// Convert product pictures
		var photos []models.ProductPictureResponse
		for _, pic := range detail.ProductLog.Product.ProductPicture {
			photos = append(photos, models.ProductPictureResponse{
				ID:        pic.ID,
				IDProduk:  pic.IDProduk,
				Url:       pic.Url,
				CreatedAt: pic.CreatedAt,
				UpdatedAt: pic.UpdatedAt,
			})
		}

		details = append(details, models.TransactionDetailResponse{
			ID:         detail.ID,
			Kuantitas:  detail.Kuantitas,
			HargaTotal: detail.HargaTotal,
			Store: models.StoreResponse{
				ID:        detail.Store.ID,
				NamaToko:  detail.Store.NamaToko,
				UrlFoto:   detail.Store.UrlFoto,
				CreatedAt: detail.Store.CreatedAt,
				UpdatedAt: detail.Store.UpdatedAt,
			},
			Product: models.ProductResponse{
				ID:            detail.ProductLog.Product.ID,
				NamaProduk:    detail.ProductLog.Product.NamaProduk,
				Slug:          detail.ProductLog.Product.Slug,
				HargaReseller: detail.ProductLog.Product.HargaReseller,
				HargaKonsumen: detail.ProductLog.Product.HargaKonsumen,
				Stok:          detail.ProductLog.Product.Stok,
				Deskripsi:     detail.ProductLog.Product.Deskripsi,
				Store: models.StoreResponse{
					ID:        detail.ProductLog.Product.Store.ID,
					NamaToko:  detail.ProductLog.Product.Store.NamaToko,
					UrlFoto:   detail.ProductLog.Product.Store.UrlFoto,
					CreatedAt: detail.ProductLog.Product.Store.CreatedAt,
					UpdatedAt: detail.ProductLog.Product.Store.UpdatedAt,
				},
				Category: models.CategoryResponse{
					ID:           detail.ProductLog.Product.Category.ID,
					NamaCategory: detail.ProductLog.Product.Category.NamaCategory,
					CreatedAt:    detail.ProductLog.Product.Category.CreatedAt,
					UpdatedAt:    detail.ProductLog.Product.Category.UpdatedAt,
				},
				Photos:    photos, // Use the converted photos array
				CreatedAt: detail.ProductLog.Product.CreatedAt,
				UpdatedAt: detail.ProductLog.Product.UpdatedAt,
			},
		})
	}
	response.TransactionDetails = details

	return response, nil
}

func (service *transactionServiceImpl) GetAll(limit int, page int, keyword string) (responder.Pagination, error) {
	request := responder.Pagination{}
	request.Limit = limit
	request.Page = page
	request.Keyword = keyword

	return service.repository.FindAllPagination(request)
}

func (service *transactionServiceImpl) Update(id uint, user_id uint, input models.TransactionUpdateRequest) (models.TransactionResponse, error) {
	// Get existing transaction
	transaction, err := service.repository.FindById(id)
	if err != nil {
		return models.TransactionResponse{}, err
	}

	// Check if user owns the transaction
	if transaction.Address.IDUser != user_id {
		return models.TransactionResponse{}, errors.New("forbidden")
	}

	// Update fields
	transaction.MethodBayar = input.MethodBayar

	// Save changes
	updated, err := service.repository.Update(transaction)
	if err != nil {
		return models.TransactionResponse{}, err
	}

	// Return updated transaction
	return service.GetById(updated.ID, user_id)
}

func (service *transactionServiceImpl) Delete(id uint, user_id uint) error {
	transaction, err := service.repository.FindById(id)
	if err != nil {
		return err
	}

	if transaction.Address.IDUser != user_id {
		return errors.New("forbidden")
	}

	return service.repository.Delete(id)
}

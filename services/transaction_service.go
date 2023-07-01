package services

import (
	"errors"
	"mini-project-evermos/models"
	"mini-project-evermos/models/responder"
	"mini-project-evermos/repositories"
	"strconv"
	"time"
)

// Contract
type TransactionService interface {
	GetAll(limit int, page int, keyword string) (responder.Pagination, error)
	GetById(id uint, user_id uint) (models.TransactionResponse, error)
	Create(input models.TransactionRequest, user_id uint) (bool, error)
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

func (service *transactionServiceImpl) GetAll(limit int, page int, keyword string) (responder.Pagination, error) {
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

func (service *transactionServiceImpl) Create(input models.TransactionRequest, user_id uint) (bool, error) {
	//check alamat
	check_address, err := service.repositoryAddress.FindById(input.AlamatKirim)

	if check_address.IDUser != user_id {
		return false, errors.New("forbidden")
	}

	date_now := time.Now()
	string_date := date_now.Format("2006-01-02-15-04-05")
	invoice := "INV" + string_date

	//transaction detail and log product
	productLogsFormatter := []models.ProductLogProcess{}
	total := 0
	for _, v := range input.DetailTrx {
		product, err := service.repositoryProduct.FindById(v.ProductID)

		if err != nil {
			return false, err
		}

		//transaction detail
		stok, _ := strconv.Atoi(product.HargaKonsumen)
		total_detail := stok * v.Kuantitas

		//log product
		productLogFormatter := models.ProductLogProcess{}
		productLogFormatter.ProductID = product.ID
		productLogFormatter.NamaProduk = product.NamaProduk
		productLogFormatter.Slug = product.Slug
		productLogFormatter.HargaReseller = product.HargaReseller
		productLogFormatter.HargaKonsumen = product.HargaKonsumen
		productLogFormatter.Stok = product.Stok
		productLogFormatter.Deskripsi = *product.Deskripsi
		productLogFormatter.CategoryID = product.Category.ID
		productLogFormatter.StoreID = product.Store.ID
		productLogFormatter.Kuantitas = v.Kuantitas
		productLogFormatter.HargaTotal = total_detail

		total += total_detail

		productLogsFormatter = append(productLogsFormatter, productLogFormatter)
	}

	//transaction
	transaction := models.TransactionProcess{}
	transaction.AlamatKirim = input.AlamatKirim
	transaction.MethodBayar = input.MethodBayar
	transaction.KodeInvoice = invoice
	transaction.UserID = user_id
	transaction.HargaTotal = total

	transaction_data := models.TransactionProcessData{}
	transaction_data.Transaction = transaction
	transaction_data.LogProduct = productLogsFormatter

	response, err := service.repository.Insert(transaction_data)

	if err != nil {
		return response, err
	}
	return true, nil
}

func (service *transactionServiceImpl) GetById(id uint, user_id uint) (models.TransactionResponse, error) {
	transaction, err := service.repository.FindById(id)

	if err != nil {
		return models.TransactionResponse{}, err
	}

	if transaction.Address.IDUser != user_id {
		return models.TransactionResponse{}, errors.New("forbidden")
	}

	var response = models.TransactionResponse{}
	response.ID = transaction.ID
	response.HargaTotal = transaction.HargaTotal
	response.KodeInvoice = transaction.KodeInvoice
	response.MethodBayar = transaction.MethodBayar
	response.Address.ID = transaction.Address.ID
	response.Address.JudulAlamat = transaction.Address.JudulAlamat
	response.Address.NamaPenerima = transaction.Address.NamaPenerima
	response.Address.NoTelp = transaction.Address.NoTelp
	response.Address.DetailAlamat = transaction.Address.DetailAlamat

	detailsFormatter := []models.TransactionDetailResponse{}

	for _, detail := range transaction.TrxDetail {
		detailFormatter := models.TransactionDetailResponse{}
		detailFormatter.ID = detail.ID
		detailFormatter.Kuantitas = detail.Kuantitas
		detailFormatter.HargaTotal = detail.HargaTotal
		detailFormatter.Store.ID = detail.Store.ID
		detailFormatter.Store.NamaToko = detail.Store.NamaToko
		detailFormatter.Store.UrlFoto = detail.Store.UrlFoto

		detailFormatter.Product.ID = detail.ProductLog.Product.ID
		detailFormatter.Product.NamaProduk = detail.ProductLog.Product.NamaProduk
		detailFormatter.Product.Slug = detail.ProductLog.Product.Slug
		detailFormatter.Product.HargaReseller = detail.ProductLog.Product.HargaReseller
		detailFormatter.Product.HargaKonsumen = detail.ProductLog.Product.HargaKonsumen
		detailFormatter.Product.Stok = detail.ProductLog.Product.Stok
		detailFormatter.Product.Deskripsi = detail.ProductLog.Product.Deskripsi
		detailFormatter.Product.Store.ID = detail.ProductLog.Product.Store.ID
		detailFormatter.Product.Store.NamaToko = detail.ProductLog.Product.Store.NamaToko
		detailFormatter.Product.Store.UrlFoto = detail.ProductLog.Product.Store.UrlFoto
		detailFormatter.Product.Category.ID = detail.ProductLog.Product.Category.ID
		detailFormatter.Product.Category.NamaCategory = detail.ProductLog.Product.Category.NamaCategory

		photosFormatter := []models.ProductPictureResponse{}

		for _, photo := range detail.ProductLog.Product.ProductPicture {
			photoFormatter := models.ProductPictureResponse{}
			photoFormatter.ID = photo.ID
			photoFormatter.IDProduk = photo.IDProduk
			photoFormatter.Url = photo.Url

			photosFormatter = append(photosFormatter, photoFormatter)
		}
		detailFormatter.Product.Photos = photosFormatter

		detailsFormatter = append(detailsFormatter, detailFormatter)
	}

	response.TransactionDetails = detailsFormatter

	return response, nil
}

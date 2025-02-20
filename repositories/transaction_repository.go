package repositories

import (
	"mini-project-evermos/models"
	"mini-project-evermos/models/entities"
	"mini-project-evermos/models/responder"

	"math"

	"gorm.io/gorm"
)

// Contract
type TransactionRepository interface {
	FindAllPagination(pagination responder.Pagination) (responder.Pagination, error)
	FindById(id uint) (entities.Trx, error)
	Insert(transaction models.TransactionProcessData) (bool, error)
}

type transactionRepositoryImpl struct {
	database *gorm.DB
}

func NewTransactionRepository(database *gorm.DB) TransactionRepository {
	return &transactionRepositoryImpl{database}
}

func (repository *transactionRepositoryImpl) FindAllPagination(pagination responder.Pagination) (responder.Pagination, error) {
	var transactions []entities.Trx
	var totalRows int64

	query := repository.database.Model(&entities.Trx{})
	query.Count(&totalRows)

	err := query.
		Preload("User").
		Preload("TrxDetail").
		Preload("TrxDetail.ProductLog").
		Preload("TrxDetail.Store").
		Limit(pagination.Limit).
		Offset(pagination.GetOffset()).
		Find(&transactions).Error

	if err != nil {
		return responder.Pagination{}, err
	}

	var responses []models.TransactionResponse
	for _, trx := range transactions {
		response := models.TransactionResponse{
			ID:          trx.ID,
			HargaTotal:  trx.HargaTotal,
			KodeInvoice: trx.KodeInvoice,
			MethodBayar: trx.MethodBayar,
			Address: models.AddressResponse{
				ID:           trx.Address.ID,
				JudulAlamat:  trx.Address.JudulAlamat,
				NamaPenerima: trx.Address.NamaPenerima,
				NoTelp:       trx.Address.NoTelp,
				DetailAlamat: trx.Address.DetailAlamat,
			},
			CreatedAt: trx.CreatedAt,
			UpdatedAt: trx.UpdatedAt,
		}
		responses = append(responses, response)
	}

	pagination.Rows = responses
	pagination.TotalRows = totalRows
	pagination.TotalPages = int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))

	return pagination, nil
}

func (repository *transactionRepositoryImpl) FindById(id uint) (entities.Trx, error) {
	var transaction entities.Trx

	err := repository.database.
		Preload("Address").
		Preload("TrxDetail").
		Preload("TrxDetail.Store").
		Preload("TrxDetail.ProductLog.Product").
		Preload("TrxDetail.ProductLog.Product.Store").
		Preload("TrxDetail.ProductLog.Product.Category").
		Preload("TrxDetail.ProductLog.Product.ProductPicture").
		Where("id = ?", id).First(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (repository *transactionRepositoryImpl) Insert(transaction models.TransactionProcessData) (bool, error) {
	tx := repository.database.Begin()
	transaction_insert := &entities.Trx{
		IDUser:           transaction.Transaction.UserID,
		AlamatPengiriman: transaction.Transaction.AlamatKirim,
		HargaTotal:       transaction.Transaction.HargaTotal,
		KodeInvoice:      transaction.Transaction.KodeInvoice,
		MethodBayar:      transaction.Transaction.MethodBayar,
	}

	if err := tx.Create(transaction_insert).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	for _, v := range transaction.LogProduct {
		log_product := &entities.ProductLog{
			IDProduk:      v.ProductID,
			NamaProduk:    v.NamaProduk,
			Slug:          v.Slug,
			HargaReseller: v.HargaReseller,
			HargaKonsumen: v.HargaKonsumen,
			Deskripsi:     &v.Deskripsi,
			IDToko:        v.StoreID,
			IDCategory:    v.CategoryID,
		}
		if err := tx.Create(log_product).Error; err != nil {
			tx.Rollback()
			return false, err
		}

		if err := tx.Create(&entities.TrxDetail{
			IDTrx:       transaction_insert.ID,
			IDLogProduk: log_product.ID,
			IDToko:      v.StoreID,
			Kuantitas:   v.Kuantitas,
			HargaTotal:  v.HargaTotal,
		}).Error; err != nil {
			tx.Rollback()
			return false, err
		}
	}

	tx.Commit()
	return true, nil
}

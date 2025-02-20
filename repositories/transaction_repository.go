package repositories

import (
	"mini-project-evermos/models"
	"mini-project-evermos/models/entities"
	"mini-project-evermos/models/responder"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindAllPagination(pagination responder.Pagination) (responder.Pagination, error)
	FindById(id uint) (entities.Trx, error)
	Insert(transaction models.TransactionProcessData) (uint, error)
	Update(transaction entities.Trx) (entities.Trx, error)
	Delete(id uint) error
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
		Preload("Address").
		Preload("TrxDetail.ProductLog.Product").
		Preload("TrxDetail.ProductLog.Product.Store").
		Preload("TrxDetail.ProductLog.Product.Category").
		Preload("TrxDetail.ProductLog.Product.ProductPicture").
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
			CreatedAt:   trx.CreatedAt,
			UpdatedAt:   trx.UpdatedAt,
			Address: models.AddressResponse{
				ID:           trx.Address.ID,
				JudulAlamat:  trx.Address.JudulAlamat,
				NamaPenerima: trx.Address.NamaPenerima,
				NoTelp:       trx.Address.NoTelp,
				DetailAlamat: trx.Address.DetailAlamat,
				CreatedAt:    trx.Address.CreatedAt,
				UpdatedAt:    trx.Address.UpdatedAt,
			},
		}

		var details []models.TransactionDetailResponse
		for _, detail := range trx.TrxDetail {
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
					Photos:    photos,
					CreatedAt: detail.ProductLog.Product.CreatedAt,
					UpdatedAt: detail.ProductLog.Product.UpdatedAt,
				},
			})
		}
		response.TransactionDetails = details
		responses = append(responses, response)
	}

	pagination.Rows = responses
	pagination.TotalRows = totalRows
	pagination.TotalPages = int((totalRows + int64(pagination.Limit) - 1) / int64(pagination.Limit))

	return pagination, nil
}

func (repository *transactionRepositoryImpl) FindById(id uint) (entities.Trx, error) {
	var transaction entities.Trx

	err := repository.database.
		Preload("Address").
		Preload("TrxDetail.ProductLog.Product").
		Preload("TrxDetail.ProductLog.Product.Store").
		Preload("TrxDetail.ProductLog.Product.Category").
		Preload("TrxDetail.ProductLog.Product.ProductPicture").
		Preload("TrxDetail.Store").
		Where("id = ?", id).
		First(&transaction).Error

	return transaction, err
}

func (repository *transactionRepositoryImpl) Insert(transaction models.TransactionProcessData) (uint, error) {
	tx := repository.database.Begin()

	transaction_insert := &entities.Trx{
		IDUser:           transaction.Transaction.UserID,
		AlamatPengiriman: transaction.Transaction.AlamatPengiriman,
		HargaTotal:       transaction.Transaction.HargaTotal,
		KodeInvoice:      transaction.Transaction.KodeInvoice,
		MethodBayar:      transaction.Transaction.MethodBayar,
	}

	if err := tx.Create(transaction_insert).Error; err != nil {
		tx.Rollback()
		return 0, err
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
			return 0, err
		}

		if err := tx.Create(&entities.TrxDetail{
			IDTrx:       transaction_insert.ID,
			IDLogProduk: log_product.ID,
			IDToko:      v.StoreID,
			Kuantitas:   v.Kuantitas,
			HargaTotal:  v.HargaTotal,
		}).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()
	return transaction_insert.ID, nil
}

func (repository *transactionRepositoryImpl) Update(transaction entities.Trx) (entities.Trx, error) {
	err := repository.database.Save(&transaction).Error
	return transaction, err
}

func (repository *transactionRepositoryImpl) Delete(id uint) error {
	return repository.database.Delete(&entities.Trx{}, id).Error
}

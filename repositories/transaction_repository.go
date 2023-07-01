package repositories

import (
	"mini-project-evermos/models"
	"mini-project-evermos/models/entities"
	"mini-project-evermos/models/responder"

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

	keyword := "%" + pagination.Keyword + "%"
	where_value := func(keyword string) *gorm.DB {
		if keyword != "" {
			return repository.database.Where("kode_invoice LIKE ?", keyword).Or("method_bayar LIKE ?", keyword)
		}
		return repository.database
	}

	err := where_value(keyword).
		Preload("Address").
		Preload("TrxDetail").
		Preload("TrxDetail.Store").
		Preload("TrxDetail.ProductLog.Product").
		Preload("TrxDetail.ProductLog.Product.Store").
		Preload("TrxDetail.ProductLog.Product.Category").
		Preload("TrxDetail.ProductLog.Product.ProductPicture").
		Scopes(responder.PaginationFormat(keyword, transactions, &pagination, where_value(keyword))).
		Find(&transactions).Error

	if err != nil {
		return pagination, err
	}

	transactionsFormatter := []models.TransactionResponse{}

	for _, transaction := range transactions {
		transactionFormatter := models.TransactionResponse{}
		transactionFormatter.ID = transaction.ID
		transactionFormatter.HargaTotal = transaction.HargaTotal
		transactionFormatter.KodeInvoice = transaction.KodeInvoice
		transactionFormatter.MethodBayar = transaction.MethodBayar
		transactionFormatter.Address.ID = transaction.Address.ID
		transactionFormatter.Address.JudulAlamat = transaction.Address.JudulAlamat
		transactionFormatter.Address.NamaPenerima = transaction.Address.NamaPenerima
		transactionFormatter.Address.NoTelp = transaction.Address.NoTelp
		transactionFormatter.Address.DetailAlamat = transaction.Address.DetailAlamat

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

		transactionFormatter.TransactionDetails = detailsFormatter

		transactionsFormatter = append(transactionsFormatter, transactionFormatter)
	}

	pagination.Data = transactionsFormatter

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

package repositories

import (
	"mini-project-evermos/models"
	"mini-project-evermos/models/entities"
	"mini-project-evermos/models/responder"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

// Contract
type ProductRepository interface {
	FindAllPagination(pagination responder.Pagination) (responder.Pagination, error)
	FindById(id uint) (entities.Product, error)
	Insert(product models.ProductRequest) (bool, error)
	Update(product models.ProductRequest, id uint) (bool, error)
	Destroy(id uint) (bool, error)
}

type productRepositoryImpl struct {
	database *gorm.DB
}

func NewProductRepository(database *gorm.DB) ProductRepository {
	return &productRepositoryImpl{database}
}

func (repository *productRepositoryImpl) FindAllPagination(pagination responder.Pagination) (responder.Pagination, error) {
	var products []entities.Product

	keyword := "%" + pagination.Keyword + "%"
	// where_value := func(keyword string) *gorm.DB {
	// 	return repository.database.Where("nama_toko LIKE ?", keyword)
	// }

	// err := where_value(keyword).
	// 	Scopes(responder.PaginationFormat(keyword, stores, &pagination, where_value(keyword))).
	// 	Find(&stores).Error

	err := repository.database.
		Preload("Store").
		Preload("Category").
		Preload("ProductPicture").
		Scopes(responder.PaginationFormat(keyword, products, &pagination, repository.database)).
		Find(&products).Error

	if err != nil {
		return pagination, err
	}

	productsFormatter := []models.ProductResponse{}

	for _, product := range products {
		productFormatter := models.ProductResponse{}
		productFormatter.ID = product.ID
		productFormatter.NamaProduk = product.NamaProduk
		productFormatter.Slug = product.Slug
		productFormatter.HargaReseller = product.HargaReseller
		productFormatter.HargaKonsumen = product.HargaKonsumen
		productFormatter.Stok = product.Stok
		productFormatter.Deskripsi = product.Deskripsi
		productFormatter.Store.ID = product.Store.ID
		productFormatter.Store.NamaToko = product.Store.NamaToko
		productFormatter.Store.UrlFoto = product.Store.UrlFoto
		productFormatter.Category.ID = product.Category.ID
		productFormatter.Category.NamaCategory = product.Category.NamaCategory

		picturesFormatter := []models.ProductPictureResponse{}

		for _, picture := range product.ProductPicture {
			pictureFormatter := models.ProductPictureResponse{}
			pictureFormatter.ID = picture.ID
			pictureFormatter.IDProduk = picture.IDProduk
			pictureFormatter.Url = picture.Url

			picturesFormatter = append(picturesFormatter, pictureFormatter)
		}
		productFormatter.Photos = picturesFormatter
		productsFormatter = append(productsFormatter, productFormatter)
	}

	pagination.Data = productsFormatter

	return pagination, nil
}

func (repository *productRepositoryImpl) FindById(id uint) (entities.Product, error) {
	var product entities.Product

	err := repository.database.
		Preload("Store").
		Preload("Category").
		Preload("ProductPicture").
		Where("id = ?", id).First(&product).Error

	if err != nil {
		return product, err
	}

	return product, nil
}

func (repository *productRepositoryImpl) Insert(product models.ProductRequest) (bool, error) {
	tx := repository.database.Begin()

	create_product := &entities.Product{
		NamaProduk:    product.NamaProduk,
		Slug:          slug.Make(product.NamaProduk),
		HargaReseller: product.HargaReseller,
		HargaKonsumen: product.HargaKonsumen,
		Stok:          product.Stok,
		Deskripsi:     &product.Deskripsi,
		IDCategory:    product.CategoryID,
		IDToko:        product.StoreID,
	}

	if err := tx.Create(create_product).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	for _, v := range product.Photos {
		if err := tx.Create(&entities.ProductPicture{
			IDProduk: create_product.ID,
			Url:      v,
		}).Error; err != nil {
			tx.Rollback()
			return false, err
		}
	}

	tx.Commit()

	return true, nil
}

func (repository *productRepositoryImpl) Update(product models.ProductRequest, id uint) (bool, error) {
	tx := repository.database.Begin()
	update_product := &entities.Product{
		NamaProduk:    product.NamaProduk,
		Slug:          slug.Make(product.NamaProduk),
		HargaReseller: product.HargaReseller,
		HargaKonsumen: product.HargaKonsumen,
		Stok:          product.Stok,
		Deskripsi:     &product.Deskripsi,
		IDCategory:    product.CategoryID,
		IDToko:        product.StoreID,
	}

	if err := tx.Where("id = ?", id).Updates(update_product).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	if err := tx.Where("id_produk = ?", id).Delete(&entities.ProductPicture{}).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	for _, v := range product.Photos {
		if err := tx.Create(&entities.ProductPicture{
			IDProduk: id,
			Url:      v,
		}).Error; err != nil {
			tx.Rollback()
			return false, err
		}
	}

	tx.Commit()

	return true, nil
}

func (repository *productRepositoryImpl) Destroy(id uint) (bool, error) {
	var product entities.Product
	err := repository.database.Where("id = ?", id).Delete(&product).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

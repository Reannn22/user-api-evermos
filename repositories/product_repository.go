package repositories

import (
	"mini-project-evermos/models"
	"mini-project-evermos/models/entities"
	"mini-project-evermos/models/responder"
	"time"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

// Contract
type ProductRepository interface {
	FindAllPagination(pagination responder.Pagination) (responder.Pagination, error)
	FindById(id uint) (entities.Product, error)
	Insert(product models.ProductRequest) (entities.Product, error)
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
		productFormatter.CreatedAt = product.CreatedAt
		productFormatter.UpdatedAt = product.UpdatedAt

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

func (repository *productRepositoryImpl) Insert(input models.ProductRequest) (entities.Product, error) {
	now := time.Now()
	product := entities.Product{
		NamaProduk:    input.NamaProduk,
		IDToko:        input.StoreID,
		IDCategory:    input.CategoryID,
		HargaReseller: input.HargaReseller,
		HargaKonsumen: input.HargaKonsumen,
		Stok:          input.Stok,
		Deskripsi:     &input.Deskripsi,
		Slug:          slug.Make(input.NamaProduk),
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}

	err := repository.database.Create(&product).Error
	if err != nil {
		return entities.Product{}, err
	}

	// Insert product pictures
	for _, url := range input.PhotoURLs {
		picture := entities.ProductPicture{
			IDProduk: product.ID,
			Url:      url,
		}
		err = repository.database.Create(&picture).Error
		if err != nil {
			return entities.Product{}, err
		}
	}

	return product, nil
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

	for _, photo := range product.PhotoURLs {
		productPicture := entities.ProductPicture{}
		productPicture.IDProduk = id
		productPicture.Url = photo
		repository.database.Create(&productPicture)
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

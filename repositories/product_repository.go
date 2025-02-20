package repositories

import (
	"fmt"
	"math"
	"mini-project-evermos/models"
	"mini-project-evermos/models/entities"
	"mini-project-evermos/models/responder"
	"time"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

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

func (repository *productRepositoryImpl) FindAllPagination(request responder.Pagination) (responder.Pagination, error) {
	var products []entities.Product
	var totalRows int64

	// Count total rows
	query := repository.database.Model(&entities.Product{})
	if request.Keyword != "" {
		query = query.Where("nama_produk LIKE ?", "%"+request.Keyword+"%")
	}
	query.Count(&totalRows)

	// Update preload to ensure complete data loading
	query = repository.database.Model(&entities.Product{}).
		Preload("Store").
		Preload("Category").
		Preload("ProductPicture").
		Order("id desc")

	if request.Keyword != "" {
		query = query.Where("nama_produk LIKE ?", "%"+request.Keyword+"%")
	}

	err := query.
		Limit(request.Limit).
		Offset(request.GetOffset()).
		Find(&products).Error

	if err != nil {
		return responder.Pagination{}, err
	}

	var responses []models.ProductResponse
	for _, product := range products {
		response := models.ProductResponse{
			ID:            product.ID,
			NamaProduk:    product.NamaProduk,
			Slug:          product.Slug,
			HargaReseller: product.HargaReseller,
			HargaKonsumen: product.HargaKonsumen,
			Stok:          product.Stok,
			Deskripsi:     product.Deskripsi,
			CreatedAt:     product.CreatedAt,
			UpdatedAt:     product.UpdatedAt,
			Store: models.StoreResponse{
				ID:        product.Store.ID,
				NamaToko:  product.Store.NamaToko,
				UrlFoto:   product.Store.UrlFoto,
				CreatedAt: product.Store.CreatedAt,
				UpdatedAt: product.Store.UpdatedAt,
			},
			Category: models.CategoryResponse{
				ID:           product.Category.ID,
				NamaCategory: product.Category.NamaCategory,
				CreatedAt:    product.Category.CreatedAt,
				UpdatedAt:    product.Category.UpdatedAt,
			},
		}

		var pictures []models.ProductPictureResponse
		for _, pic := range product.ProductPicture {
			pictures = append(pictures, models.ProductPictureResponse{
				ID:        pic.ID,
				IDProduk:  pic.IDProduk,
				Url:       pic.Url,
				CreatedAt: pic.CreatedAt,
				UpdatedAt: pic.UpdatedAt,
			})
		}
		response.Photos = pictures
		responses = append(responses, response)
	}

	request.Rows = responses
	request.TotalRows = totalRows
	request.TotalPages = int(math.Ceil(float64(totalRows) / float64(request.Limit)))

	return request, nil
}

func (repository *productRepositoryImpl) FindById(id uint) (entities.Product, error) {
	var product entities.Product

	fmt.Printf("Looking for product with ID: %d in table: %s\n", id, entities.Product{}.TableName())

	result := repository.database.Debug().
		Table(entities.Product{}.TableName()).
		Preload("Store", func(db *gorm.DB) *gorm.DB {
			return db.Table("toko") // explicitly use toko table
		}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Table("category") // explicitly use category table
		}).
		Preload("ProductPicture", func(db *gorm.DB) *gorm.DB {
			return db.Table("foto_produk") // explicitly use foto_produk table
		}).
		Where("id = ?", id).
		First(&product)

	if result.Error != nil {
		fmt.Printf("Database error: %v\n", result.Error)
		return entities.Product{}, fmt.Errorf("product with ID %d not found: %v", id, result.Error)
	}

	fmt.Printf("Found product: %+v\n", product)
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
		productPicture := entities.ProductPicture{
			IDProduk: id,
			Url:      photo,
		}
		if err := tx.Create(&productPicture).Error; err != nil {
			tx.Rollback()
			return false, err
		}
	}

	tx.Commit()
	return true, nil
}

func (repository *productRepositoryImpl) Destroy(id uint) (bool, error) {
	err := repository.database.Where("id = ?", id).Delete(&entities.Product{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

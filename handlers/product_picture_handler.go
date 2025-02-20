package handlers

import (
	"mini-project-evermos/models"
	"mini-project-evermos/models/entities"
	"mini-project-evermos/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	productPictureRepo repositories.ProductPictureRepository
}

func NewHandler(productPictureRepo repositories.ProductPictureRepository) *Handler {
	return &Handler{
		productPictureRepo: productPictureRepo,
	}
}

func (h *Handler) CreateProductPicture(c *gin.Context) {
	var request models.ProductPictureRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "BAD REQUEST",
			"errors":  err.Error(),
			"data":    nil,
		})
		return
	}

	// Verify if product exists
	if !h.productPictureRepo.VerifyProductExists(request.ProductID) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Product not found",
			"errors":  "Product with given ID does not exist",
			"data":    nil,
		})
		return
	}

	productPicture := &entities.ProductPicture{
		IDProduk: request.ProductID, // Changed from ProductID to IDProduk
		Url:      request.URL,       // Changed from URL to Url
	}

	if err := h.productPictureRepo.Create(productPicture); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "INTERNAL SERVER ERROR",
			"errors":  err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Successfully created product picture",
		"errors":  nil,
		"data": models.ProductPictureResponse{
			ID:       productPicture.ID,
			IDProduk: productPicture.IDProduk,
			Url:      productPicture.Url,
		},
	})
}

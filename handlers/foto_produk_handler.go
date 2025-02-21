package handlers

import (
	"mini-project-evermos/exceptions"
	"mini-project-evermos/middleware"
	"mini-project-evermos/models"
	"mini-project-evermos/models/responder"
	"mini-project-evermos/services"
	"mini-project-evermos/utils/jwt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type FotoProdukHandler struct {
	service services.FotoProdukService
}

func NewFotoProdukHandler(service *services.FotoProdukService) FotoProdukHandler {
	return FotoProdukHandler{*service}
}

func (h *FotoProdukHandler) Route(app *fiber.App) {
	routes := app.Group("/api/v1")
	routes.Use(middleware.JWTProtected())
	routes.Get("/product-photos", h.GetAll)
	routes.Get("/product-photos/:id", h.GetById)
	routes.Get("/product/:id/photos", h.GetByProductId)
	routes.Post("/product-photos", h.Create)
	routes.Put("/product-photos/:id", h.Update)
	routes.Delete("/product-photos/:id", h.Delete)
}

func (h *FotoProdukHandler) GetAll(c *fiber.Ctx) error {
	responses, err := h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to get photos",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	return c.JSON(responder.ApiResponse{
		Status:  true,
		Message: "Success get photos",
		Error:   nil,
		Data:    responses,
	})
}

func (h *FotoProdukHandler) GetById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Invalid ID",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	response, err := h.service.GetById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to get photo",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	return c.JSON(responder.ApiResponse{
		Status:  true,
		Message: "Success get photo",
		Error:   nil,
		Data:    response,
	})
}

func (h *FotoProdukHandler) GetByProductId(c *fiber.Ctx) error {
	productId, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Invalid product ID",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	responses, err := h.service.GetByProductId(uint(productId))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to get photos",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	return c.JSON(responder.ApiResponse{
		Status:  true,
		Message: "Success get photos",
		Error:   nil,
		Data:    responses,
	})
}

func (h *FotoProdukHandler) Create(c *fiber.Ctx) error {
	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Unauthorized",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	// Parse JSON input directly
	var input models.FotoProdukRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Invalid input",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	// Create the photo record
	response, err := h.service.Create(input, uint(claims.UserId))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to create photo",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	return c.Status(http.StatusCreated).JSON(responder.ApiResponse{
		Status:  true,
		Message: "Success create photo",
		Error:   nil,
		Data:    response,
	})
}

func (h *FotoProdukHandler) Update(c *fiber.Ctx) error {
	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Unauthorized",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Invalid ID",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	var input models.FotoProdukRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Invalid input",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	response, err := h.service.Update(uint(id), input, uint(claims.UserId))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to update photo",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	return c.JSON(responder.ApiResponse{
		Status:  true,
		Message: "Success update photo",
		Error:   nil,
		Data:    response,
	})
}

func (h *FotoProdukHandler) Delete(c *fiber.Ctx) error {
	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Unauthorized",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Invalid ID",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	// Get the deleted photo data
	deletedPhoto, err := h.service.Delete(uint(id), uint(claims.UserId))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to delete photo",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	return c.JSON(responder.ApiResponse{
		Status:  true,
		Message: "Success delete photo",
		Error:   nil,
		Data:    deletedPhoto,
	})
}

package handlers

import (
	"fmt"
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

type TransactionHandler struct {
	TransactionService services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) TransactionHandler {
	return TransactionHandler{*transactionService}
}

func (handler *TransactionHandler) Route(app *fiber.App) {
	routes := app.Group("/api/v1/trx")
	routes.Get("/", middleware.JWTProtected(), handler.GetAllTransaction)
	routes.Get("/:id", middleware.JWTProtected(), handler.DetailTransaction)
	routes.Post("/", middleware.JWTProtected(), handler.CreateTransaction)
	routes.Put("/:id", middleware.JWTProtected(), handler.UpdateTransaction)
	routes.Delete("/:id", middleware.JWTProtected(), handler.DeleteTransaction)
}

func (handler *TransactionHandler) GetAllTransaction(c *fiber.Ctx) error {
	// Default values
	defaultLimit := 10
	defaultPage := 1

	// Try to get limit from query param, use default if not provided
	limit, err := strconv.Atoi(c.FormValue("limit", strconv.Itoa(defaultLimit)))
	if err != nil {
		limit = defaultLimit
	}

	// Try to get page from query param, use default if not provided
	page, err := strconv.Atoi(c.FormValue("page", strconv.Itoa(defaultPage)))
	if err != nil {
		page = defaultPage
	}

	keyword := c.FormValue("search", "")

	responses, err := handler.TransactionService.GetAll(limit, page, keyword)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to GET data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	return c.Status(http.StatusOK).JSON(responder.ApiResponse{
		Status:  true,
		Message: "Succeed to GET data",
		Error:   nil,
		Data:    responses,
	})
}

func (handler *TransactionHandler) DetailTransaction(c *fiber.Ctx) error {
	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Unauthorized",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to GET data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	response, err := handler.TransactionService.GetById(uint(id), uint(claims.UserId))
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(http.StatusNotFound).JSON(responder.ApiResponse{
				Status:  false,
				Message: "Transaction not found",
				Error:   exceptions.NewString(err.Error()),
				Data:    nil,
			})
		}
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to GET data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	return c.Status(http.StatusOK).JSON(responder.ApiResponse{
		Status:  true,
		Message: "Succeed to GET data",
		Error:   nil,
		Data:    response,
	})
}

func (handler *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	// Debug incoming request
	fmt.Println("=== Transaction Create Request ===")

	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		fmt.Printf("JWT Error: %v\n", err)
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to POST data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}
	fmt.Printf("User ID from token: %d\n", claims.UserId)

	var input models.TransactionRequest
	err = c.BodyParser(&input)
	if err != nil {
		fmt.Printf("Body Parse Error: %v\n", err)
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to parse request body",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	fmt.Printf("Parsed Input: %+v\n", input)

	response, err := handler.TransactionService.Create(input, uint(claims.UserId))
	if err != nil {
		fmt.Printf("Service Error: %v\n", err)
		return c.Status(http.StatusNotFound).JSON(responder.ApiResponse{
			Status:  false,
			Message: "NOT FOUND",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	return c.Status(http.StatusOK).JSON(responder.ApiResponse{
		Status:  true,
		Message: "Succeed to POST data",
		Error:   nil,
		Data:    response,
	})
}

func (handler *TransactionHandler) UpdateTransaction(c *fiber.Ctx) error {
	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Unauthorized",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Invalid ID parameter",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	var input models.TransactionUpdateRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to parse request body",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	response, err := handler.TransactionService.Update(uint(id), uint(claims.UserId), input)
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(http.StatusNotFound).JSON(responder.ApiResponse{
				Status:  false,
				Message: "Transaction not found",
				Error:   exceptions.NewString(err.Error()),
				Data:    nil,
			})
		}
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to update transaction",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	return c.Status(http.StatusOK).JSON(responder.ApiResponse{
		Status:  true,
		Message: "Transaction updated successfully",
		Error:   nil,
		Data:    response,
	})
}

func (handler *TransactionHandler) DeleteTransaction(c *fiber.Ctx) error {
	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Unauthorized",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Invalid ID parameter",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	// Get transaction data before deletion
	response, err := handler.TransactionService.GetById(uint(id), uint(claims.UserId))
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(http.StatusNotFound).JSON(responder.ApiResponse{
				Status:  false,
				Message: "Transaction not found",
				Error:   exceptions.NewString(err.Error()),
				Data:    nil,
			})
		}
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to get transaction",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	// Delete the transaction
	err = handler.TransactionService.Delete(uint(id), uint(claims.UserId))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to delete transaction",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	return c.Status(http.StatusOK).JSON(responder.ApiResponse{
		Status:  true,
		Message: "Transaction deleted successfully",
		Error:   nil,
		Data:    response, // Return the transaction data that was deleted
	})
}

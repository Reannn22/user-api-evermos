package handlers

import (
	"mini-project-evermos/exceptions"
	"mini-project-evermos/models"
	"mini-project-evermos/models/responder"
	"mini-project-evermos/services"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	AuthService services.AuthService
}

func NewAuthHandler(authService *services.AuthService) AuthHandler {
	return AuthHandler{*authService}
}

// Route
func (handler *AuthHandler) Route(app *fiber.App) {
	routes := app.Group("/api/v1/auth")
	routes.Post("/register", handler.Register)
	routes.Post("/login", handler.Login)
}

func (handler *AuthHandler) Register(c *fiber.Ctx) error {
	var input models.RegisterRequest
	err := c.BodyParser(&input)

	// exception.ValidationForm(err)

	err = handler.AuthService.Register(input)

	if err != nil {
		//error
		if err.Error() == "unique" {
			return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
				Status:  false,
				Message: "Failed to POST data",
				Error:   exceptions.NewString("No Telp Registered, Please Login"),
				Data:    nil,
			})
		}
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to POST data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	return c.Status(http.StatusCreated).JSON(responder.ApiResponse{
		Status:  true,
		Message: "Succeed to POST data",
		Error:   nil,
		Data:    "Register Succeed",
	})
}

func (handler *AuthHandler) Login(c *fiber.Ctx) error {
	var input models.LoginRequest

	err := c.BodyParser(&input)

	// exception.ValidationForm(err)

	responses, err := handler.AuthService.Login(input)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to POST data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	return c.Status(http.StatusOK).JSON(responder.ApiResponse{
		Status:  true,
		Message: "Succeed to POST data",
		Error:   nil,
		Data:    responses,
	})
}

package handlers

import (
	"mini-project-evermos/exceptions"
	"mini-project-evermos/middleware"
	"mini-project-evermos/models"
	"mini-project-evermos/models/responder"
	"mini-project-evermos/services"
	"mini-project-evermos/utils/jwt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserService services.UserService
}

func NewUserHandler(userService *services.UserService) UserHandler {
	return UserHandler{*userService}
}

func (handler *UserHandler) Route(app *fiber.App) {
	routes := app.Group("/api/v1/user")
	routes.Get("/", middleware.JWTProtected(), handler.UserDetail)
	routes.Put("/", middleware.JWTProtected(), handler.UserUpdate)
	routes.Delete("/", middleware.JWTProtected(), handler.UserDelete)
}

func (handler *UserHandler) UserDetail(c *fiber.Ctx) error {
	//claim
	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to GET data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	user_id := claims.UserId

	responses, err := handler.UserService.GetById(uint(user_id))
	if err != nil {
		//error
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

func (handler *UserHandler) UserUpdate(c *fiber.Ctx) error {
	//claim
	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to PUT data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	user_id := claims.UserId

	var input models.UserRequest
	err = c.BodyParser(&input)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to PUT data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	responses, err := handler.UserService.Edit(uint(user_id), input)
	if err != nil {
		//error
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to PUT data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}
	return c.Status(http.StatusOK).JSON(responder.ApiResponse{
		Status:  true,
		Message: "Succeed to PUT data",
		Error:   nil,
		Data:    responses, // Now returns the updated user data
	})
}

func (handler *UserHandler) GetProfile(c *fiber.Ctx) error {
	// Get user_id from token claims
	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(responder.ApiResponse{
			Status:  false,
			Message: "UNAUTHORIZED",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	// Get user data using the ID
	user, err := handler.UserService.GetById(uint(claims.UserId))
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(responder.ApiResponse{
			Status:  false,
			Message: "NOT FOUND",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	return c.Status(http.StatusOK).JSON(responder.ApiResponse{
		Status:  true,
		Message: "SUCCESS",
		Error:   nil,
		Data:    user,
	})
}

func (handler *UserHandler) UserDelete(c *fiber.Ctx) error {
	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to DELETE data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	user_id := claims.UserId

	// Get user data before deletion
	response, err := handler.UserService.Delete(uint(user_id))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to DELETE data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	return c.Status(http.StatusOK).JSON(responder.ApiResponse{
		Status:  true,
		Message: "Succeed to DELETE data",
		Error:   nil,
		Data:    response, // Return the deleted user data
	})
}

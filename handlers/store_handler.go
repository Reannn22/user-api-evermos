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

type StoreHandler struct {
	StoreService services.StoreService
}

func NewStoreHandler(storeService *services.StoreService) StoreHandler {
	return StoreHandler{*storeService}
}

func (handler *StoreHandler) Route(app *fiber.App) {
	fmt.Println("Registering store routes") // Add debug log
	routes := app.Group("/api/v1/toko")
	routes.Post("/", middleware.JWTProtected(), handler.StoreCreate)
	routes.Get("/my", middleware.JWTProtected(), handler.MyStore)
	routes.Get("/", middleware.JWTProtected(), handler.GetAllStore)
	routes.Get("/:id_toko", middleware.JWTProtected(), handler.StoreDetail)
	routes.Put("/:id_toko", middleware.JWTProtected(), handler.EditStore)
	routes.Delete("/:id_toko", middleware.JWTProtected(), handler.DeleteStore)
}

func (handler *StoreHandler) MyStore(c *fiber.Ctx) error {
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

	responses, err := handler.StoreService.GetByUserId(uint(user_id))
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

func (handler *StoreHandler) GetAllStore(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.FormValue("limit"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to GET data",
			Error:   exceptions.NewString("limit required."),
			Data:    nil,
		})
	}

	page, err := strconv.Atoi(c.FormValue("page"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to GET data",
			Error:   exceptions.NewString("page required."),
			Data:    nil,
		})
	}

	keyword := c.FormValue("nama")

	responses, err := handler.StoreService.GetAll(limit, page, keyword)

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

func (handler *StoreHandler) StoreDetail(c *fiber.Ctx) error {
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

	id, err := c.ParamsInt("id_toko")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to GET data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	response, err := handler.StoreService.GetById(uint(id), uint(user_id))
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
		Data:    response,
	})
}

func (handler *StoreHandler) EditStore(c *fiber.Ctx) error {
	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to PUT data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	user_id := claims.UserId
	id_toko, err := c.ParamsInt("id_toko")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to PUT data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	name_store := c.FormValue("nama_toko")
	url_foto := c.FormValue("url_foto")

	// Handle both file upload and URL cases
	formHeader, err := c.FormFile("photo")
	if err == nil {
		// If file was uploaded, use the filename
		url_foto = formHeader.Filename
	} else if url_foto == "" {
		// If no file and no URL provided, keep existing photo
		store, err := handler.StoreService.GetById(uint(id_toko), uint(user_id))
		if err == nil && store.UrlFoto != nil {
			url_foto = *store.UrlFoto
		}
	}

	input := models.StoreProcess{
		ID:       uint(id_toko),
		UserID:   uint(user_id),
		NamaToko: &name_store,
		URL:      url_foto,
	}

	response, err := handler.StoreService.Edit(input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to PUT data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	// Save file only if one was uploaded
	if formHeader != nil {
		saveErr := c.SaveFile(formHeader, fmt.Sprintf("uploads/%s", *response.UrlFoto))
		if saveErr != nil {
			return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
				Status:  false,
				Message: "Failed to PUT data",
				Error:   exceptions.NewString(saveErr.Error()),
				Data:    nil,
			})
		}
	}

	return c.Status(http.StatusOK).JSON(responder.ApiResponse{
		Status:  true,
		Message: "Succeed to PUT data",
		Error:   nil,
		Data:    response,
	})
}

func (handler *StoreHandler) StoreCreate(c *fiber.Ctx) error {
	fmt.Println("StoreCreate handler called")
	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to POST data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	user_id := claims.UserId
	nama_toko := c.FormValue("nama_toko")
	url_foto := c.FormValue("url_foto")

	// Only use default if url_foto is empty string or wasn't provided at all
	if url_foto == "" {
		url_foto = "https://avatars.githubusercontent.com/u/174382151?s=48&v=4"
	}

	fmt.Printf("Using URL: %s\n", url_foto) // Debug log

	input := models.StoreProcess{
		UserID:   uint(user_id),
		NamaToko: &nama_toko,
		URL:      url_foto,
	}

	response, err := handler.StoreService.Create(input)
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
		Data:    response,
	})
}

func (handler *StoreHandler) DeleteStore(c *fiber.Ctx) error {
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
	id_toko, err := c.ParamsInt("id_toko")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to DELETE data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	response, err := handler.StoreService.Delete(uint(id_toko), uint(user_id))
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
		Data:    response,
	})
}

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
	"time"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	ProductService services.ProductService
}

func NewProductHandler(productService *services.ProductService) ProductHandler {
	return ProductHandler{*productService}
}

func (handler *ProductHandler) Route(app *fiber.App) {
	routes := app.Group("/api/v1/product")
	routes.Get("/", middleware.JWTProtected(), handler.GetAllProduct)
	routes.Get("/:id", middleware.JWTProtected(), handler.ProductDetail)
	routes.Post("/", middleware.JWTProtected(), handler.ProductCreate)
	routes.Put("/:id", middleware.JWTProtected(), handler.ProductUpdate)
	routes.Delete("/:id", middleware.JWTProtected(), handler.ProductDelete)
}

func (handler *ProductHandler) GetAllProduct(c *fiber.Ctx) error {
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

	keyword := c.FormValue("keyword")

	responses, err := handler.ProductService.GetAll(limit, page, keyword)

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

func (handler *ProductHandler) ProductDetail(c *fiber.Ctx) error {
	//claim
	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to POST data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	user_id := claims.UserId

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to GET data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	response, err := handler.ProductService.GetById(uint(id), uint(user_id))
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

func (handler *ProductHandler) ProductCreate(c *fiber.Ctx) error {
	//claim
	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to POST data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	user_id := claims.UserId

	formFile, err := c.MultipartForm()
	var file_name []string
	for _, fileHeaders := range formFile.File {
		for _, fileHeader := range fileHeaders {
			date_now := time.Now()
			string_date := date_now.Format("2006_01_02_15_04_05")

			filename := string_date + "-" + fileHeader.Filename
			c.SaveFile(fileHeader, fmt.Sprintf("uploads/%s", filename))
			file_name = append(file_name, filename)
		}
	}

	category_id, err := strconv.Atoi(c.FormValue("category_id"))
	stok, err := strconv.Atoi(c.FormValue("stok"))

	input := models.ProductRequest{}
	input.NamaProduk = c.FormValue("nama_produk")
	input.CategoryID = uint(category_id)
	input.HargaReseller = c.FormValue("harga_reseller")
	input.HargaKonsumen = c.FormValue("harga_konsumen")
	input.Stok = stok
	input.Deskripsi = c.FormValue("deskripsi")
	input.Photos = file_name

	response, err := handler.ProductService.Create(input, uint(user_id))

	if err != nil {
		//error
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

func (handler *ProductHandler) ProductUpdate(c *fiber.Ctx) error {
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

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to PUT data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	formFile, err := c.MultipartForm()
	var file_name []string
	for _, fileHeaders := range formFile.File {
		for _, fileHeader := range fileHeaders {
			date_now := time.Now()
			string_date := date_now.Format("2006_01_02_15_04_05")

			filename := string_date + "-" + fileHeader.Filename
			c.SaveFile(fileHeader, fmt.Sprintf("uploads/%s", filename))
			file_name = append(file_name, filename)
		}
	}

	category_id, err := strconv.Atoi(c.FormValue("category_id"))
	stok, err := strconv.Atoi(c.FormValue("stok"))

	input := models.ProductRequest{}
	input.NamaProduk = c.FormValue("nama_produk")
	input.CategoryID = uint(category_id)
	input.HargaReseller = c.FormValue("harga_reseller")
	input.HargaKonsumen = c.FormValue("harga_konsumen")
	input.Stok = stok
	input.Deskripsi = c.FormValue("deskripsi")
	input.Photos = file_name

	response, err := handler.ProductService.Update(input, uint(id), uint(user_id))

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
		Data:    response,
	})
}

func (handler *ProductHandler) ProductDelete(c *fiber.Ctx) error {
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

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Failed to DELETE data",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	response, err := handler.ProductService.Delete(uint(id), uint(user_id))

	if err != nil {
		//error
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

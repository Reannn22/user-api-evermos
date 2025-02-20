package routers

import (
	"mini-project-evermos/services"

	"github.com/gofiber/fiber"
)

// ...existing code...

func RegisterTrxRoutes(router fiber.Router, service services.TransactionService) {
	router.Post("/trx", handlers.TransactionCreateHandler(service))
	router.Get("/trx", handlers.TransactionAllHandler(service))
	router.Get("/trx/:id", handlers.TransactionGetByIdHandler(service)) // Add this line
	// ...existing code...
}

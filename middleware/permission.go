package middleware

import (
	"mini-project-evermos/exceptions"
	"mini-project-evermos/models/responder"
	"mini-project-evermos/utils/jwt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Admin
func RolePermissionAdmin(c *fiber.Ctx) error {
	// Get claims from JWT.
	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(http.StatusInternalServerError).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Something Wrong",
			Error:   exceptions.NewString(err.Error()),
			Data:    nil,
		})
	}

	is_admin := claims.IsAdmin
	if is_admin != true {
		return c.Status(http.StatusForbidden).JSON(responder.ApiResponse{
			Status:  false,
			Message: "Only Admin Access",
			Error:   exceptions.NewString("Forbidden"),
			Data:    nil,
		})
	}

	return c.Next()
}

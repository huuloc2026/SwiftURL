package response

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

// Success returns a standard success response
func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// Error returns a standard error response, hiding sensitive details
func Error(c *fiber.Ctx, status int, err error, publicMsg string) error {
	// Log the real error for debugging
	if err != nil {
		log.Printf("error: %v", err)
	}
	env := os.Getenv("ENV")
	showDetail := env == "development" || env == "debug"
	msg := publicMsg
	if msg == "" {
		if showDetail && err != nil {
			msg = err.Error()
		} else {
			msg = "An error occurred"
		}
	}
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"error":   msg,
	})
}

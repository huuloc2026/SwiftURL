package routes

import (
	"github.com/gofiber/fiber/v2"
	handler "github.com/huuloc2026/SwiftURL/internal/shorturl/delivery/http"
)

func RegisterRoutes(app *fiber.App, h *handler.URLHandler) {

	app.Get("/:code", h.ResolveURL)
	api := app.Group("/api")
	api.Post("/shorten", h.CreateShortURL)
	api.Delete("/shorten/:code", h.DeleteShortURL)

	// auth := api.Group("/auth")
	// admin := api.Group("/admin")

}

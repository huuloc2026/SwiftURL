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

	// User routes (example)
	// import user handler and usecase at the top
	// userHandler := userhandler.NewUserHandler(userUsecase)
	// api := app.Group("/api")
	// api.Post("/users/register", userHandler.Register)
	// api.Get("/users/:id", userHandler.GetByID)

	// auth := api.Group("/auth")
	// admin := api.Group("/admin")

}

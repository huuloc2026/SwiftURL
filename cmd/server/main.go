package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	handler "github.com/huuloc2026/SwiftURL/internal/shorturl/delivery/http"
	"github.com/huuloc2026/SwiftURL/internal/shorturl/repository"
	"github.com/huuloc2026/SwiftURL/internal/shorturl/usecase"
	"github.com/huuloc2026/SwiftURL/pkg/database"
)

func main() {
	// Initialize Fiber app in debug mode
	app := fiber.New(fiber.Config{
		AppName: "SwiftURL",
		Prefork: false,
		// Enable debug mode
		DisableStartupMessage: false,
	})

	// Global middlewares
	app.Use(logger.New())
	app.Use(recover.New())

	// Initialize dependencies
	db := database.InitDB()
	urlRepo := repository.NewShortURLRepository(db)
	urlUC := usecase.NewShortURLUsecase(urlRepo)
	h := handler.NewURLHandler(urlUC)

	// 🔍 Health check route
	app.Get("/healthz", func(c *fiber.Ctx) error {
		if err := db.DB.Ping(); err != nil {
			return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "unhealthy",
				"message": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})
	// Group all API routes under /api
	api := app.Group("/api")

	api.Post("/shorten", h.CreateShortURL)
	api.Get("/:code", h.ResolveURL)

	// Start server
	log.Fatal(app.Listen(":8080"))
}

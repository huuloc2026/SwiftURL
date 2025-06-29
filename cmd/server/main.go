package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/huuloc2026/SwiftURL/config"
	"github.com/huuloc2026/SwiftURL/internal/dependencies"

	"github.com/huuloc2026/SwiftURL/pkg/database"
)

func main() {
	config.LoadEnv()
	app := fiber.New(fiber.Config{
		AppName:               "SwiftURL",
		Prefork:               false,
		DisableStartupMessage: false,
	})

	app.Use(logger.New())
	app.Use(recover.New())

	db := database.InitDB()

	deps := dependencies.InitDependencies(db)

	app.Get("/healthz", func(c *fiber.Ctx) error {
		if err := db.DB.Ping(); err != nil {
			return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	dependencies.RegisterRoutes(app, deps)

	log.Fatal(app.Listen(":8080"))
}

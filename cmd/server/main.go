package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/huuloc2026/SwiftURL/config"
	"github.com/huuloc2026/SwiftURL/internal/routes"
	urlhandler "github.com/huuloc2026/SwiftURL/internal/shorturl/delivery/http"
	urlrepo "github.com/huuloc2026/SwiftURL/internal/shorturl/repository"
	urlusecase "github.com/huuloc2026/SwiftURL/internal/shorturl/usecase"
	"github.com/huuloc2026/SwiftURL/pkg/database"

	authhandler "github.com/huuloc2026/SwiftURL/internal/auth/handler"
	authusecase "github.com/huuloc2026/SwiftURL/internal/auth/usecase"
	authrepo "github.com/huuloc2026/SwiftURL/internal/user/repository"
	// import your OTPService and JWTService implementations as needed
)

func main() {
	config.LoadEnv()
	// Initialize Fiber app in debug mode
	app := fiber.New(fiber.Config{
		AppName:               "SwiftURL",
		Prefork:               false,
		DisableStartupMessage: false,
	})

	// Global middlewares
	app.Use(logger.New())
	app.Use(recover.New())

	// Initialize DB connection
	db := database.InitDB() // or your actual DB initialization

	urlRepo := urlrepo.NewShortURLRepository(db)
	urlUC := urlusecase.NewShortURLUsecase(urlRepo)
	h := urlhandler.NewURLHandler(urlUC)

	// Initialize AuthHandler dependencies
	userRepo := authrepo.NewUserRepository(db)
	// otpService := /* initialize your OTPService implementation */
	// jwtService := /* initialize your JWTService implementation */
	authUC := authusecase.NewAuthUsecase(userRepo, otpService, jwtService)
	authH := authhandler.NewAuthHandler(authUC)

	// Health check route
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

	routes.RegisterRoutes(app, h, authH)

	// Start server
	log.Fatal(app.Listen(":8080"))
}

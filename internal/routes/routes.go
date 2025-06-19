package routes

import (
	"github.com/gofiber/fiber/v2"
	authhandler "github.com/huuloc2026/SwiftURL/internal/auth/delivery/http"
	urlhandler "github.com/huuloc2026/SwiftURL/internal/shorturl/delivery/http"
)

func RegisterRoutes(app *fiber.App, h *urlhandler.URLHandler, authHandler *authhandler.AuthHandler) {

	app.Get("/:code", h.ResolveURL)
	api := app.Group("/api")
	api.Post("/shorten", h.CreateShortURL)
	api.Delete("/shorten/:code", h.DeleteShortURL)

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/forget-password", authHandler.ForgetPassword)
	auth.Post("/verify-otp", authHandler.VerifyOTP)
	auth.Post("/change-password", authHandler.ChangePassword)

	// User routes (example)
	// import user handler and usecase at the top
	// userHandler := userhandler.NewUserHandler(userUsecase)
	// api := app.Group("/api")
	// api.Post("/users/register", userHandler.Register)
	// api.Get("/users/:id", userHandler.GetByID)

	// auth := api.Group("/auth")
	// admin := api.Group("/admin")

}

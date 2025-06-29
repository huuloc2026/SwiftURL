package dependencies

import (
	"github.com/gofiber/fiber/v2"
	authhandler "github.com/huuloc2026/SwiftURL/internal/auth/delivery/http"
	authusecase "github.com/huuloc2026/SwiftURL/internal/auth/usecase"
	urlhandler "github.com/huuloc2026/SwiftURL/internal/shorturl/delivery/http"
	urlrepo "github.com/huuloc2026/SwiftURL/internal/shorturl/repository"
	urlusecase "github.com/huuloc2026/SwiftURL/internal/shorturl/usecase"
	userrepo "github.com/huuloc2026/SwiftURL/internal/user/repository"
	jwtService "github.com/huuloc2026/SwiftURL/pkg/jwt"
	"github.com/huuloc2026/SwiftURL/pkg/otp"
	"github.com/jmoiron/sqlx"
)

type Dependencies struct {
	URLHandler  *urlhandler.URLHandler
	AuthHandler *authhandler.AuthHandler
}

func InitDependencies(db *sqlx.DB) *Dependencies {
	// ShortURL module
	urlRepo := urlrepo.NewShortURLRepository(db)
	urlUC := urlusecase.NewShortURLUsecase(urlRepo)
	urlHandler := urlhandler.NewURLHandler(urlUC)

	// Auth module
	userRepo := userrepo.NewUserRepository(db)
	otpService := otp.NewInMemoryOTPService()
	jwt := jwtService.NewSimpleJWTService("your-secret-key")
	authUC := authusecase.NewAuthUsecase(userRepo, otpService, jwt)
	authHandler := authhandler.NewAuthHandler(authUC)

	return &Dependencies{
		URLHandler:  urlHandler,
		AuthHandler: authHandler,
	}
}

func RegisterRoutes(app *fiber.App, deps *Dependencies) {
	// Short URL routes
	app.Get("/:code", deps.URLHandler.ResolveURL)
	api := app.Group("/api")
	api.Post("/shorten", deps.URLHandler.CreateShortURL)
	api.Delete("/shorten/:code", deps.URLHandler.DeleteShortURL)

	// Auth routes under /api/auth/
	auth := api.Group("/auth")
	auth.Post("/login", deps.AuthHandler.Login)
	auth.Post("/register", deps.AuthHandler.Register)
	auth.Post("/verify-otp", deps.AuthHandler.VerifyOTP)
	auth.Post("/change-password", deps.AuthHandler.ChangePassword)
	auth.Post("/forget-password", deps.AuthHandler.ForgetPassword)
}

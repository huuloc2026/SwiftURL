package dependencies

import (
	authHandler "github.com/huuloc2026/SwiftURL/internal/auth/delivery/http"
	"github.com/huuloc2026/SwiftURL/internal/auth/repository"
	"github.com/huuloc2026/SwiftURL/internal/auth/usecase"

	"github.com/huuloc2026/SwiftURL/pkg/database"
)

type Dependencies struct {
	URLHandler  *userHandler.URLHandler
	AuthHandler *authHandler.AuthHandler
}

func InitDependencies(db *database.DB) *Dependencies {
	// Initialize URL dependencies
	urlRepo := repository.NewShortURLRepository(db)
	urlUC := usecase.NewShortURLUsecase(urlRepo)
	urlHandler := http.NewURLHandler(urlUC)

	// Initialize Auth dependencies
	userRepo := repository.NewUserRepository(db)
	otpService := NewInMemoryOTPService()
	jwtService := jwt.NewSimpleJWTService("your-secret-key")
	authUC := usecase.NewAuthUsecase(userRepo, otpService, jwtService)
	authHandler := http.NewAuthHandler(authUC)

	return &Dependencies{
		URLHandler:  urlHandler,
		AuthHandler: authHandler,
	}
}

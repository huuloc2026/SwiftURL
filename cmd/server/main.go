package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/huuloc2026/SwiftURL/config"
	"github.com/huuloc2026/SwiftURL/internal/routes"
	urlhandler "github.com/huuloc2026/SwiftURL/internal/shorturl/delivery/http"
	urlrepo "github.com/huuloc2026/SwiftURL/internal/shorturl/repository"
	urlusecase "github.com/huuloc2026/SwiftURL/internal/shorturl/usecase"
	"github.com/huuloc2026/SwiftURL/pkg/database"

	authhandler "github.com/huuloc2026/SwiftURL/internal/auth/delivery/http"
	authusecase "github.com/huuloc2026/SwiftURL/internal/auth/usecase"
	authrepo "github.com/huuloc2026/SwiftURL/internal/user/repository"
)

// --- Simple In-Memory OTPService Implementation ---
type InMemoryOTPService struct {
	store map[string]string
	mu    sync.Mutex
}

func NewInMemoryOTPService() *InMemoryOTPService {
	return &InMemoryOTPService{
		store: make(map[string]string),
	}
}

func (s *InMemoryOTPService) GenerateOTP(ctx context.Context, email string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	s.store[email] = otp
	return otp, nil
}

func (s *InMemoryOTPService) VerifyOTP(ctx context.Context, email, otp string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.store[email] == otp {
		delete(s.store, email)
		return true, nil
	}
	return false, nil
}

// --- Simple JWTService Implementation ---
type SimpleJWTService struct {
	secret string
}

func NewSimpleJWTService(secret string) *SimpleJWTService {
	return &SimpleJWTService{secret: secret}
}

func (s *SimpleJWTService) GenerateToken(userID int64, username string, exp time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(exp).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

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
	otpService := NewInMemoryOTPService()
	jwtService := NewSimpleJWTService("your-secret-key")
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

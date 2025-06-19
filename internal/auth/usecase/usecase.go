package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/huuloc2026/SwiftURL/internal/entity"
	"github.com/huuloc2026/SwiftURL/internal/user/repository"
	"golang.org/x/crypto/bcrypt"
)

// OTPService abstracts OTP generation/verification
type OTPService interface {
	GenerateOTP(ctx context.Context, email string) (string, error)
	VerifyOTP(ctx context.Context, email, otp string) (bool, error)
}

// JWTService abstracts JWT token generation
type JWTService interface {
	GenerateToken(userID int64, username string, exp time.Duration) (string, error)
}

// AuthUsecase defines authentication business logic
type AuthUsecase interface {
	Register(ctx context.Context, username, email, password string) (int64, error)
	Login(ctx context.Context, email, password string) (string, error)
	ForgetPassword(ctx context.Context, email string) error
	VerifyOTP(ctx context.Context, email, otp string) (bool, error)
	ChangePassword(ctx context.Context, email, otp, newPassword string) error
}

type authUsecase struct {
	userRepo   repository.UserRepository
	otpService OTPService
	jwtService JWTService
}

// NewAuthUsecase is the constructor with dependency injection
func NewAuthUsecase(userRepo repository.UserRepository, otpService OTPService, jwtService JWTService) AuthUsecase {
	return &authUsecase{
		userRepo:   userRepo,
		otpService: otpService,
		jwtService: jwtService,
	}
}

func (a *authUsecase) Register(ctx context.Context, username, email, password string) (int64, error) {
	hashed := hashPassword(password)
	user := &entity.User{
		Username: username,
		Email:    email,
		Password: hashed,
	}
	return a.userRepo.Create(ctx, user)
}

func (a *authUsecase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}
	if !checkPassword(password, user.Password) {
		return "", errors.New("invalid credentials")
	}
	return a.jwtService.GenerateToken(user.ID, user.Username, time.Hour*24)
}

func (a *authUsecase) ForgetPassword(ctx context.Context, email string) error {
	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	_, err = a.otpService.GenerateOTP(ctx, email)
	return err
}

func (a *authUsecase) VerifyOTP(ctx context.Context, email, otp string) (bool, error) {
	return a.otpService.VerifyOTP(ctx, email, otp)
}

func (a *authUsecase) ChangePassword(ctx context.Context, email, otp, newPassword string) error {
	ok, err := a.otpService.VerifyOTP(ctx, email, otp)
	if err != nil || !ok {
		return errors.New("invalid OTP")
	}
	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	hashed := hashPassword(newPassword)
	_, err = a.userRepo.UpdateByID(ctx, user.ID, nil, nil, &hashed)
	return err
}

// --- Helper functions ---
func hashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// In production, handle error properly
		return ""
	}
	return string(hashed)
}

func checkPassword(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

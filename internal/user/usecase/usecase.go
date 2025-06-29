package usecase

import (
	"context"
	"errors"

	"github.com/huuloc2026/SwiftURL/internal/entity"
	"github.com/huuloc2026/SwiftURL/internal/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type ctxKey string

const usernameCtxKey ctxKey = "username"

type UserUsecase interface {
	Register(ctx context.Context, username, email, password string) (int64, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	Update(ctx context.Context, id int64, username, email, password *string) (*entity.User, error)
	Delete(ctx context.Context, id int64) (*entity.User, error)
	List(ctx context.Context, limit, offset int) ([]*entity.User, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo}
}

func (u *userUsecase) Register(ctx context.Context, username, email, password string) (int64, error) {
	if username == "" || email == "" || password == "" {
		return 0, errors.New("username, email, and password cannot be empty")
	}
	if len(password) < 6 {
		return 0, errors.New("password must be at least 6 characters long")
	}
	if len(username) < 3 || len(username) > 20 {
		return 0, errors.New("username must be between 3 and 20 characters long")
	}
	if len(email) < 5 || len(email) > 50 {
		return 0, errors.New("email must be between 5 and 50 characters long")
	}
	if len(password) < 6 || len(password) > 100 {
		return 0, errors.New("password must be between 6 and 100 characters long")
	}
	// Check if username already exists
	ctx = context.WithValue(ctx, usernameCtxKey, username)
	existing, _ := u.repo.FindByUsername(ctx, username)
	if existing != nil {
		return 0, errors.New("username already exists")
	}
	// Hash password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.New("failed to hash password")
	}
	user := &entity.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}
	return u.repo.Create(ctx, user)
}

func (u *userUsecase) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *userUsecase) Update(ctx context.Context, id int64, username, email, password *string) (*entity.User, error) {
	return u.repo.UpdateByID(ctx, id, username, email, password)
}

func (u *userUsecase) Delete(ctx context.Context, id int64) (*entity.User, error) {
	return u.repo.DeleteByID(ctx, id)
}

func (u *userUsecase) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	return u.repo.List(ctx, limit, offset)
}

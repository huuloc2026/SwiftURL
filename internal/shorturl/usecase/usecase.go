package usecase

import (
	"context"
	"time"

	"github.com/huuloc2026/SwiftURL/internal/entity"
	"github.com/huuloc2026/SwiftURL/internal/shorturl/repository"
	"github.com/huuloc2026/SwiftURL/pkg/utils"
)

type ShortURLUsecase interface {
	Shorten(ctx context.Context, longURL string) (*entity.ShortURL, error)
	Resolve(ctx context.Context, shortCode string) (*entity.ShortURL, error)
}

type shortURLUsecase struct {
	repo repository.ShortURLRepository
}

func NewShortURLUsecase(repo repository.ShortURLRepository) ShortURLUsecase {
	return &shortURLUsecase{repo}
}

func (u *shortURLUsecase) Shorten(ctx context.Context, longURL string) (*entity.ShortURL, error) {
	code := utils.GenerateShortCode(6)
	now := time.Now()

	url := &entity.ShortURL{
		ShortCode: code,
		LongURL:   longURL,
		CreatedAt: now,
		Clicks:    0,
	}

	err := u.repo.Store(ctx, url)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (u *shortURLUsecase) Resolve(ctx context.Context, shortCode string) (*entity.ShortURL, error) {
	url, err := u.repo.FindByCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	_ = u.repo.IncrementClick(ctx, shortCode) // Không chặn user nếu lỗi log
	return url, nil
}

package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/huuloc2026/SwiftURL/internal/entity"
	"github.com/huuloc2026/SwiftURL/internal/shorturl/repository"
	"github.com/huuloc2026/SwiftURL/pkg/utils"
)

type ShortURLUsecase interface {
	Generate(ctx context.Context, longURL string, customCode *string, expireAt *time.Time) (string, error)
	Shorten(ctx context.Context, longURL string) (*entity.ShortURL, error)
	Resolve(ctx context.Context, shortCode string) (*entity.ShortURL, error)
	UpdateCode(ctx context.Context, code string, custom_code string, longURL string, expireAt *time.Time) (*entity.ShortURL, error)
	Delete(ctx context.Context, code string) error
	//TrackClick
	TrackClick(ctx context.Context, log *entity.ClickLog) error
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
	//TODO: DISABLE FOR TESTING ALGORITHM - DONT INSERT TO REPO
	err := u.repo.Store(ctx, url)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (u *shortURLUsecase) Generate(ctx context.Context, longURL string, customCode *string, expireAt *time.Time) (string, error) {
	var code string

	if customCode != nil {
		existing, err := u.repo.FindByCode(ctx, *customCode)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("error checking custom code: %w", err)
		}
		if existing != nil {
			return "", errors.New("short code already exists")
		}
		code = *customCode
	} else {
		const maxAttempts = 100
		for i := 0; i < maxAttempts; i++ {
			candidate := utils.GenerateShortCode(8)
			existing, err := u.repo.FindByCode(ctx, candidate)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return "", fmt.Errorf("error checking generated code: %w", err)
			}
			if existing != nil {
				code = candidate
				break
			}

		}
	}

	short := &entity.ShortURL{
		ShortCode: code,
		LongURL:   longURL,
		CreatedAt: time.Now(),
		ExpireAt:  expireAt,
		Clicks:    0,
	}

	if err := u.repo.Store(ctx, short); err != nil {
		return "", fmt.Errorf("failed to store short url: %w", err)
	}

	return code, nil
}

func (u *shortURLUsecase) Delete(ctx context.Context, code string) error {
	return u.repo.DeleteByCode(ctx, code)
}

func (u *shortURLUsecase) Resolve(ctx context.Context, code string) (*entity.ShortURL, error) {
	url, err := u.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if url.ExpireAt != nil && url.ExpireAt.Before(time.Now()) {
		return nil, errors.New("short url expired")
	}
	_ = u.repo.IncrementClick(ctx, code)

	return url, nil
}

func (u *shortURLUsecase) UpdateCode(ctx context.Context, code string, custom_code string, longURL string, expireAt *time.Time) (*entity.ShortURL, error) {
	existing, err := u.repo.FindByCode(ctx, code)

	if err != nil {
		return nil, fmt.Errorf("error checking existing code: %w", err)
	}
	if existing == nil {
		return nil, errors.New("short code does not exist")
	}
	fmt.Println("Existing URL:", code)
	fmt.Print("Custom Code:", custom_code)

	var expireAtStr *string
	if expireAt != nil {
		str := expireAt.Format(time.RFC3339)
		expireAtStr = &str
	}

	updatedURL, err := u.repo.UpdateCode(ctx, code, longURL, expireAtStr)

	if err != nil {
		return nil, fmt.Errorf("failed to update short url: %w", err)
	}

	return updatedURL, nil
}

func (u *shortURLUsecase) TrackClick(ctx context.Context, log *entity.ClickLog) error {
	geo, err := utils.LookupIP(log.IPAddress)
	if err == nil {
		log.Country = geo.Country
		log.City = geo.City
	}
	return u.repo.InsertClickLog(ctx, log)
}

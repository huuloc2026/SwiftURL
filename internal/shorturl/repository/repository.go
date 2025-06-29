package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/huuloc2026/SwiftURL/internal/entity"
	"github.com/jmoiron/sqlx"
)

func isSensitiveError() bool {
	env := os.Getenv("ENV")
	return env == "development" || env == "debug"
}

type ShortURLRepository interface {
	Store(ctx context.Context, url *entity.ShortURL) error
	FindByCode(ctx context.Context, code string) (*entity.ShortURL, error)
	IncrementClick(ctx context.Context, code string) error
	DeleteByCode(ctx context.Context, code string) error
	ExistsByCode(ctx context.Context, code string) (bool, error)
	UpdateCode(ctx context.Context, code string, longURL string, expireAt *string) (*entity.ShortURL, error)

	//CLICK LOG
	InsertClickLog(ctx context.Context, log *entity.ClickLog) error
}

type shortURLRepo struct {
	db *sqlx.DB
}

func NewShortURLRepository(db *sqlx.DB) ShortURLRepository {
	return &shortURLRepo{db}
}

func (r *shortURLRepo) Store(ctx context.Context, url *entity.ShortURL) error {
	_, err := r.db.ExecContext(ctx, sqlInsertShortURL,
		url.ShortCode,
		url.LongURL,
		url.CreatedAt,
		url.ExpireAt,
		url.Clicks,
	)
	if err != nil && !isSensitiveError() {
		return fmt.Errorf("failed to store short url")
	}
	return err
}

func (r *shortURLRepo) FindByCode(ctx context.Context, code string) (*entity.ShortURL, error) {
	var url entity.ShortURL

	err := r.db.GetContext(ctx, &url, sqlFindByCode, code)
	if err != nil && !isSensitiveError() {
		return nil, fmt.Errorf("short url not found")
	}
	return &url, err
}

func (r *shortURLRepo) IncrementClick(ctx context.Context, code string) error {
	_, err := r.db.ExecContext(ctx, sqlIncrementClick, code)
	if err != nil && !isSensitiveError() {
		return fmt.Errorf("failed to increment click")
	}
	return err
}

func (r *shortURLRepo) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var exists bool

	err := r.db.GetContext(ctx, &exists, sqlExistsByCode, code)
	if err != nil && !isSensitiveError() {
		return false, fmt.Errorf("failed to check existence")
	}
	return exists, err
}

func (r *shortURLRepo) UpdateCode(ctx context.Context, code string, longURL string, expireAt *string) (*entity.ShortURL, error) {
	var url entity.ShortURL

	err := r.db.GetContext(ctx, &url, sqlUpdateCode, longURL, expireAt, code)
	if err != nil && !isSensitiveError() {
		return nil, fmt.Errorf("failed to update short url")
	}
	return &url, err
}

func (r *shortURLRepo) DeleteByCode(ctx context.Context, code string) error {
	_, err := r.db.ExecContext(ctx, sqlDeleteByCode, code)
	if err != nil && !isSensitiveError() {
		return fmt.Errorf("failed to delete short url")
	}
	return err
}

func (r *shortURLRepo) FindValidByCode(ctx context.Context, code string) (*entity.ShortURL, error) {
	var url entity.ShortURL

	err := r.db.GetContext(ctx, &url, sqlFindValidByCode, code)
	if err != nil && !isSensitiveError() {
		return nil, fmt.Errorf("short url not found or expired")
	}
	return &url, err
}

func (r *shortURLRepo) InsertClickLog(ctx context.Context, log *entity.ClickLog) error {
	var clickLog entity.ClickLog
	fmt.Println(clickLog)
	err := r.db.GetContext(ctx, &clickLog, queryInsertClickLog, log.ShortCode,
		log.ClickedAt,
		log.Referrer,
		log.UserAgent,
		log.DeviceType,
		log.OS,
		log.Browser,
		log.Country,
		log.City,
		log.IPAddress)
	if err != nil && !isSensitiveError() {
		return fmt.Errorf("failed to insert click log")
	}
	return err
}

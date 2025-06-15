package repository

import (
	"context"

	"github.com/huuloc2026/SwiftURL/internal/entity"
	"github.com/jmoiron/sqlx"
)

type ShortURLRepository interface {
	Store(ctx context.Context, url *entity.ShortURL) error
	FindByCode(ctx context.Context, code string) (*entity.ShortURL, error)
	IncrementClick(ctx context.Context, code string) error
}

type shortURLRepo struct {
	db *sqlx.DB
}

func NewShortURLRepository(db *sqlx.DB) ShortURLRepository {
	return &shortURLRepo{db}
}

func (r *shortURLRepo) Store(ctx context.Context, url *entity.ShortURL) error {
	query := `INSERT INTO short_urls (short_code, long_url, created_at, expire_at, clicks)
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query,
		url.ShortCode,
		url.LongURL,
		url.CreatedAt,
		url.ExpireAt,
		url.Clicks,
	)
	return err
}

func (r *shortURLRepo) FindByCode(ctx context.Context, code string) (*entity.ShortURL, error) {
	var url entity.ShortURL
	query := `SELECT * FROM short_urls WHERE short_code = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &url, query, code)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *shortURLRepo) IncrementClick(ctx context.Context, code string) error {
	query := `UPDATE short_urls SET clicks = clicks + 1 WHERE short_code = $1`
	_, err := r.db.ExecContext(ctx, query, code)
	return err
}

package repository

import (
	"context"
	"fmt"

	"github.com/huuloc2026/SwiftURL/internal/entity"
	"github.com/jmoiron/sqlx"
)

type ShortURLRepository interface {
	Store(ctx context.Context, url *entity.ShortURL) error
	FindByCode(ctx context.Context, code string) (*entity.ShortURL, error)
	IncrementClick(ctx context.Context, code string) error
	DeleteByCode(ctx context.Context, code string) error

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
	return err
}

func (r *shortURLRepo) FindByCode(ctx context.Context, code string) (*entity.ShortURL, error) {
	var url entity.ShortURL

	err := r.db.GetContext(ctx, &url, sqlFindByCode, code)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *shortURLRepo) IncrementClick(ctx context.Context, code string) error {

	_, err := r.db.ExecContext(ctx, sqlIncrementClick, code)
	return err
}

func (r *shortURLRepo) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var exists bool

	err := r.db.GetContext(ctx, &exists, sqlExistsByCode, code)
	return exists, err
}

func (r *shortURLRepo) DeleteByCode(ctx context.Context, code string) error {
	_, err := r.db.ExecContext(ctx, sqlDeleteByCode, code)
	return err
}

func (r *shortURLRepo) FindValidByCode(ctx context.Context, code string) (*entity.ShortURL, error) {
	var url entity.ShortURL

	err := r.db.GetContext(ctx, &url, sqlFindValidByCode, code)
	if err != nil {
		return nil, err
	}
	return &url, nil
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

	return err
}

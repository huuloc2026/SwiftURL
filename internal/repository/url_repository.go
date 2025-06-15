package repository

import (
	"github.com/huuloc2026/SwiftURL.git/model"
	"github.com/jmoiron/sqlx"
)

type URLRepository interface {
	Save(url *model.URL) error
	FindByCode(code string) (*model.URL, error)
}

type urlRepo struct {
	db *sqlx.DB
}

func NewURLRepository(db *sqlx.DB) URLRepository {
	return &urlRepo{db: db}
}

func (r *urlRepo) Save(url *model.URL) error {
	return r.db.QueryRowx(
		"INSERT INTO urls (code, long_url, clicks) VALUES ($1, $2, $3) RETURNING id",
		url.Code, url.LongURL, url.Clicks).Scan(&url.ID)
}

func (r *urlRepo) FindByCode(code string) (*model.URL, error) {
	var url model.URL
	err := r.db.Get(&url, "SELECT * FROM urls WHERE code=$1", code)
	return &url, err
}

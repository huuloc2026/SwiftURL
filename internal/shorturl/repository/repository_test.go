package repository

import (
	"context"
	"regexp"
	"testing"

	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/huuloc2026/SwiftURL/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func newMockShortURLRepo(t *testing.T) (*shortURLRepo, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return &shortURLRepo{db: sqlxDB}, mock, func() { db.Close() }
}

func TestShortURLRepo_Store(t *testing.T) {
	repo, mock, close := newMockShortURLRepo(t)
	defer close()
	createdAt, _ := time.Parse(time.RFC3339, "2024-06-01T00:00:00Z")
	expireAt, _ := time.Parse(time.RFC3339, "2024-07-01T00:00:00Z")
	url := &entity.ShortURL{
		ShortCode: "abc123",
		LongURL:   "https://example.com",
		CreatedAt: createdAt,
		ExpireAt:  &expireAt,
		Clicks:    0,
	}

	mock.ExpectExec(regexp.QuoteMeta(sqlInsertShortURL)).
		WithArgs(url.ShortCode, url.LongURL, url.CreatedAt, url.ExpireAt, url.Clicks).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Store(context.Background(), url)
	assert.NoError(t, err)
}

func TestShortURLRepo_FindByCode(t *testing.T) {
	repo, mock, close := newMockShortURLRepo(t)
	defer close()
	createdAt, _ := time.Parse(time.RFC3339, "2024-06-01T00:00:00Z")
	expireAt, _ := time.Parse(time.RFC3339, "2024-07-01T00:00:00Z")
	url := &entity.ShortURL{
		ShortCode: "abc123",
		LongURL:   "https://example.com",
		CreatedAt: createdAt,
		ExpireAt:  &expireAt,
		Clicks:    0,
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlFindByCode)).
		WithArgs(url.ShortCode).
		WillReturnRows(sqlmock.NewRows([]string{"short_code", "long_url", "created_at", "expire_at", "clicks"}).
			AddRow(url.ShortCode, url.LongURL, url.CreatedAt, url.ExpireAt, url.Clicks))

	got, err := repo.FindByCode(context.Background(), url.ShortCode)
	assert.NoError(t, err)
	assert.Equal(t, url, got)
}

func TestShortURLRepo_IncrementClick(t *testing.T) {
	repo, mock, close := newMockShortURLRepo(t)
	defer close()

	code := "abc123"
	mock.ExpectExec(regexp.QuoteMeta(sqlIncrementClick)).
		WithArgs(code).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.IncrementClick(context.Background(), code)
	assert.NoError(t, err)
}

func TestShortURLRepo_DeleteByCode(t *testing.T) {
	repo, mock, close := newMockShortURLRepo(t)
	defer close()

	code := "abc123"
	mock.ExpectExec(regexp.QuoteMeta(sqlDeleteByCode)).
		WithArgs(code).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteByCode(context.Background(), code)
	assert.NoError(t, err)
}

func TestShortURLRepo_ExistsByCode(t *testing.T) {
	repo, mock, close := newMockShortURLRepo(t)
	defer close()

	code := "abc123"
	mock.ExpectQuery(regexp.QuoteMeta(sqlExistsByCode)).
		WithArgs(code).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.ExistsByCode(context.Background(), code)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestShortURLRepo_FindValidByCode(t *testing.T) {
	repo, mock, close := newMockShortURLRepo(t)
	defer close()
	createdAt, _ := time.Parse(time.RFC3339, "2024-06-01T00:00:00Z")
	url := &entity.ShortURL{
		ShortCode: "abc123",
		LongURL:   "https://example.com",
		CreatedAt: createdAt,
		Clicks:    0,
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlFindValidByCode)).
		WithArgs(url.ShortCode).
		WillReturnRows(sqlmock.NewRows([]string{"short_code", "long_url", "created_at", "expire_at", "clicks"}).
			AddRow(url.ShortCode, url.LongURL, url.CreatedAt, url.ExpireAt, url.Clicks))

	got, err := repo.FindValidByCode(context.Background(), url.ShortCode)
	assert.NoError(t, err)
	assert.Equal(t, url, got)
}

func TestShortURLRepo_InsertClickLog(t *testing.T) {
	repo, mock, close := newMockShortURLRepo(t)
	defer close()
	createdAt, _ := time.Parse(time.RFC3339, "2024-06-01T00:00:00Z")
	log := &entity.ClickLog{
		ShortCode:  "abc123",
		ClickedAt:  createdAt,
		Referrer:   "https://referrer.com",
		UserAgent:  "Mozilla/5.0",
		DeviceType: "desktop",
		OS:         "linux",
		Browser:    "firefox",
		Country:    "US",
		City:       "New York",
		IPAddress:  "127.0.0.1",
	}

	mock.ExpectQuery(regexp.QuoteMeta(queryInsertClickLog)).
		WithArgs(
			log.ShortCode,
			log.ClickedAt,
			log.Referrer,
			log.UserAgent,
			log.DeviceType,
			log.OS,
			log.Browser,
			log.Country,
			log.City,
			log.IPAddress,
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"short_code", "clicked_at", "referrer", "user_agent", "device_type", "os", "browser", "country", "city", "ip_address",
		}).AddRow(
			log.ShortCode,
			log.ClickedAt,
			log.Referrer,
			log.UserAgent,
			log.DeviceType,
			log.OS,
			log.Browser,
			log.Country,
			log.City,
			log.IPAddress,
		))

	err := repo.InsertClickLog(context.Background(), log)
	assert.NoError(t, err)
}

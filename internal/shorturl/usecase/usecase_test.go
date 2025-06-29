package usecase

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/huuloc2026/SwiftURL/internal/entity"
// 	"github.com/huuloc2026/SwiftURL/pkg/utils"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// // --- Mocks ---

// type mockShortURLRepo struct {
// 	mock.Mock
// }

// func (m *mockShortURLRepo) Store(ctx context.Context, url *entity.ShortURL) error {
// 	args := m.Called(ctx, url)
// 	return args.Error(0)
// }
// func (m *mockShortURLRepo) FindByCode(ctx context.Context, code string) (*entity.ShortURL, error) {
// 	args := m.Called(ctx, code)
// 	url, _ := args.Get(0).(*entity.ShortURL)
// 	return url, args.Error(1)
// }
// func (m *mockShortURLRepo) DeleteByCode(ctx context.Context, code string) error {
// 	args := m.Called(ctx, code)
// 	return args.Error(0)
// }
// func (m *mockShortURLRepo) IncrementClick(ctx context.Context, code string) error {
// 	args := m.Called(ctx, code)
// 	return args.Error(0)
// }
// func (m *mockShortURLRepo) UpdateCode(ctx context.Context, code string, longURL string, expireAt *string) (*entity.ShortURL, error) {
// 	args := m.Called(ctx, code, longURL, expireAt)
// 	url, _ := args.Get(0).(*entity.ShortURL)
// 	return url, args.Error(1)
// }
// func (m *mockShortURLRepo) InsertClickLog(ctx context.Context, log *entity.ClickLog) error {
// 	args := m.Called(ctx, log)
// 	return args.Error(0)
// }

// // --- Utils patching ---

// func patchGenerateShortCode(f func(n int) string) func() {
// 	orig := utils.GenerateShortCode
// 	utils.GenerateShortCode = f
// 	return func() { utils.GenerateShortCode = orig }
// }

// func patchLookupIP(f func(ip string) (*utils.GeoInfo, error)) func() {
// 	orig := utils.LookupIP
// 	utils.LookupIP = f
// 	return func() { utils.LookupIP = orig }
// }

// // --- Tests ---

// func TestShorten_Success(t *testing.T) {
// 	repo := new(mockShortURLRepo)
// 	defer patchGenerateShortCode(func(n int) string { return "abc123" })()

// 	uc := NewShortURLUsecase(repo)
// 	ctx := context.Background()
// 	longURL := "https://example.com"

// 	repo.On("Store", ctx, mock.AnythingOfType("*entity.ShortURL")).Return(nil)

// 	short, err := uc.Shorten(ctx, longURL)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "abc123", short.ShortCode)
// 	assert.Equal(t, longURL, short.LongURL)
// 	assert.Equal(t, 0, short.Clicks)
// 	repo.AssertExpectations(t)
// }

// func TestShorten_StoreError(t *testing.T) {
// 	repo := new(mockShortURLRepo)
// 	defer patchGenerateShortCode(func(n int) string { return "abc123" })()

// 	uc := NewShortURLUsecase(repo)
// 	ctx := context.Background()
// 	longURL := "https://example.com"

// 	repo.On("Store", ctx, mock.AnythingOfType("*entity.ShortURL")).Return(errors.New("db error"))

// 	short, err := uc.Shorten(ctx, longURL)
// 	assert.Nil(t, short)
// 	assert.Error(t, err)
// 	repo.AssertExpectations(t)
// }

// func TestGenerate_CustomCode_Success(t *testing.T) {
// 	repo := new(mockShortURLRepo)
// 	uc := NewShortURLUsecase(repo)
// 	ctx := context.Background()
// 	longURL := "https://example.com"
// 	customCode := "custom1"

// 	repo.On("FindByCode", ctx, customCode).Return(nil, sql.ErrNoRows)
// 	repo.On("Store", ctx, mock.AnythingOfType("*entity.ShortURL")).Return(nil)

// 	code, err := uc.Generate(ctx, longURL, &customCode, nil)
// 	assert.NoError(t, err)
// 	assert.Equal(t, customCode, code)
// 	repo.AssertExpectations(t)
// }

// func TestGenerate_CustomCode_AlreadyExists(t *testing.T) {
// 	repo := new(mockShortURLRepo)
// 	uc := NewShortURLUsecase(repo)
// 	ctx := context.Background()
// 	longURL := "https://example.com"
// 	customCode := "custom1"
// 	existing := &entity.ShortURL{ShortCode: customCode}

// 	repo.On("FindByCode", ctx, customCode).Return(existing, nil)

// 	code, err := uc.Generate(ctx, longURL, &customCode, nil)
// 	assert.Error(t, err)
// 	assert.Empty(t, code)
// 	repo.AssertExpectations(t)
// }

// func TestGenerate_GeneratedCode_Success(t *testing.T) {
// 	repo := new(mockShortURLRepo)
// 	defer patchGenerateShortCode(func(n int) string { return "gen12345" })()

// 	uc := NewShortURLUsecase(repo)
// 	ctx := context.Background()
// 	longURL := "https://example.com"

// 	repo.On("FindByCode", ctx, "gen12345").Return(nil, sql.ErrNoRows)
// 	repo.On("Store", ctx, mock.AnythingOfType("*entity.ShortURL")).Return(nil)

// 	code, err := uc.Generate(ctx, longURL, nil, nil)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "gen12345", code)
// 	repo.AssertExpectations(t)
// }

// func TestGenerate_StoreError(t *testing.T) {
// 	repo := new(mockShortURLRepo)
// 	defer patchGenerateShortCode(func(n int) string { return "gen12345" })()

// 	uc := NewShortURLUsecase(repo)
// 	ctx := context.Background()
// 	longURL := "https://example.com"

// 	repo.On("FindByCode", ctx, "gen12345").Return(nil, sql.ErrNoRows)
// 	repo.On("Store", ctx, mock.AnythingOfType("*entity.ShortURL")).Return(errors.New("db error"))

// 	code, err := uc.Generate(ctx, longURL, nil, nil)
// 	assert.Error(t, err)
// 	assert.Empty(t, code)
// 	repo.AssertExpectations(t)
// }

// func TestDelete_Success(t *testing.T) {
// 	repo := new(mockShortURLRepo)
// 	uc := NewShortURLUsecase(repo)
// 	ctx := context.Background()
// 	code := "abc123"

// 	repo.On("DeleteByCode", ctx, code).Return(nil)

// 	err := uc.Delete(ctx, code)
// 	assert.NoError(t, err)
// 	repo.AssertExpectations(t)
// }

// func TestDelete_Error(t *testing.T) {
// 	repo := new(mockShortURLRepo)
// 	uc := NewShortURLUsecase(repo)
// 	ctx := context.Background()
// 	code := "abc123"

// 	repo.On("DeleteByCode", ctx, code).Return(errors.New("not found"))

// 	err := uc.Delete(ctx, code)
// 	assert.Error(t, err)
// 	repo.AssertExpectations(t)
// }

// func TestResolve_Success(t *testing.T) {
// 	repo := new(mockShortURLRepo)
// 	uc := NewShortURLUsecase(repo)
// 	ctx := context.Background()
// 	code := "abc123"
// 	url := &entity.ShortURL{ShortCode: code, LongURL: "https://example.com", Clicks: 0}

// 	repo.On("FindByCode", ctx, code).Return(url, nil)
// 	repo.On("IncrementClick", ctx, code).Return(nil)

// 	result, err := uc.Resolve(ctx, code)
// 	assert.NoError(t, err)
// 	assert.Equal(t, url, result)
// 	repo.AssertExpectations(t)
// }

// func TestResolve_Expired(t *testing.T) {
// 	repo := new(mockShortURLRepo)
// 	uc := NewShortURLUsecase(repo)
// 	ctx := context.Background()
// 	code := "abc123"
// 	expired := time.Now().Add(-time.Hour)
// 	url := &entity.ShortURL{ShortCode: code, LongURL: "https://example.com", ExpireAt: &expired}

// 	repo.On("FindByCode", ctx, code).Return(url, nil)

// 	result, err := uc.Resolve(ctx, code)
// 	assert.Error(t, err)
// 	assert.Nil(t, result)
// 	repo.AssertExpectations(t)
// }

// func TestResolve_FindByCodeError(t *testing.T) {
// 	repo := new(mockShortURLRepo)
// 	uc := NewShortURLUsecase(repo)
// 	ctx := context.Background()
// 	code := "abc123"

// 	repo.On("FindByCode", ctx, code).Return(nil, errors.New("db error"))

// 	result, err := uc.Resolve(ctx, code)
// 	assert.Error(t, err)
// 	assert.Nil(t, result)
// 	repo.AssertExpectations(t)
// }

// func TestUpdateCode_Success(t *testing.T) {
// 	repo := new(mockShortURLRepo)
// 	uc := NewShortURLUsecase(repo)
// 	ctx := context.Background()
// 	code := "abc123"
// 	longURL := "https://new.com"
// 	now := time.Now()
// 	existing := &entity.ShortURL{ShortCode: code, LongURL: "old"}

// 	repo.On("FindByCode", ctx, code).Return(existing, nil)
// 	repo.On("UpdateCode", ctx, code, longURL, mock.AnythingOfType("*string")).Return(&entity.ShortURL{ShortCode: code, LongURL: longURL, ExpireAt: &now}, nil)

// 	updated, err := uc.UpdateCode(ctx, code, longURL, &now)
// 	assert.NoError(t, err)
// 	assert.Equal(t, longURL, updated.LongURL)
// 	repo.AssertExpectations(t)
// }

// func TestUpdateCode_NotFound(t *testing.T) {
// 	repo := new(mockShortURLRepo)
// 	uc := NewShortURLUsecase(repo)
// 	ctx := context.Background()
// 	code := "abc123"
// 	longURL := "https://new.com"

// 	repo.On("FindByCode", ctx, code).Return(nil, nil)

// 	updated, err := uc.UpdateCode(ctx, code, longURL, nil)
// 	assert.Error(t, err)
// 	assert.Nil(t, updated)
// 	repo.AssertExpectations(t)
// }

// func TestUpdateCode_UpdateError(t *testing.T) {
// 	repo := new(mockShortURLRepo)
// 	uc := NewShortURLUsecase(repo)
// 	ctx := context.Background()
// 	code := "abc123"
// 	longURL := "https://new.com"
// 	existing := &entity.ShortURL{ShortCode: code, LongURL: "old"}

// 	repo.On("FindByCode", ctx, code).Return(existing, nil)
// 	repo.On("UpdateCode", ctx, code, longURL, (*string)(nil)).Return(nil, errors.New("update error"))

// 	updated, err := uc.UpdateCode(ctx, code, longURL, nil)
// 	assert.Error(t, err)
// 	assert.Nil(t, updated)
// 	repo.AssertExpectations(t)
// }

// func TestTrackClick_Success(t *testing.T) {
// 	repo := new(mockShortURLRepo)
// 	defer patchLookupIP(func(ip string) (*utils.GeoInfo, error) {
// 		return &utils.GeoInfo{Country: "VN", City: "Hanoi"}, nil
// 	})()

// 	uc := NewShortURLUsecase(repo)
// 	ctx := context.Background()
// 	log := &entity.ClickLog{IPAddress: "1.2.3.4"}

// 	repo.On("InsertClickLog", ctx, log).Return(nil)

// 	err := uc.TrackClick(ctx, log)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "VN", log.Country)
// 	assert.Equal(t, "Hanoi", log.City)
// 	repo.AssertExpectations(t)
// }

// func TestTrackClick_LookupIPErr(t *testing.T) {
// 	repo := new(mockShortURLRepo)
// 	defer patchLookupIP(func(ip string) (*utils.GeoInfo, error) {
// 		return nil, errors.New("lookup error")
// 	})()

// 	uc := NewShortURLUsecase(repo)
// 	ctx := context.Background()
// 	log := &entity.ClickLog{IPAddress: "1.2.3.4"}

// 	repo.On("InsertClickLog", ctx, log).Return(nil)

// 	err := uc.TrackClick(ctx, log)
// 	assert.NoError(t, err)
// 	assert.Empty(t, log.Country)
// 	assert.Empty(t, log.City)
// 	repo.AssertExpectations(t)
// }

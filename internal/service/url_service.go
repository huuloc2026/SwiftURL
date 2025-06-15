package service

import (
	"math/rand"

	"github.com/huuloc2026/SwiftURL.git/internal/repository"
	"github.com/huuloc2026/SwiftURL.git/model"
)

type URLService interface {
	Shorten(url string) (*model.URL, error)
	Resolve(code string) (*model.URL, error)
}

type urlSvc struct {
	repo repository.URLRepository
}

func NewURLService(r repository.URLRepository) URLService {
	return &urlSvc{repo: r}
}

func generateCode(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *urlSvc) Shorten(longURL string) (*model.URL, error) {
	url := &model.URL{
		LongURL: longURL,
		Code:    generateCode(6),
		Clicks:  0,
	}
	err := s.repo.Save(url)
	return url, err
}

func (s *urlSvc) Resolve(code string) (*model.URL, error) {
	return s.repo.FindByCode(code)
}

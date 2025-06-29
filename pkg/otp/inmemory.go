package otp

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
)

type InMemoryOTPService struct {
	store map[string]string
	mu    sync.Mutex
}

func NewInMemoryOTPService() *InMemoryOTPService {
	return &InMemoryOTPService{
		store: make(map[string]string),
	}
}

func (s *InMemoryOTPService) GenerateOTP(ctx context.Context, email string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	s.store[email] = otp
	return otp, nil
}

func (s *InMemoryOTPService) VerifyOTP(ctx context.Context, email, otp string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.store[email] == otp {
		delete(s.store, email)
		return true, nil
	}
	return false, nil
}

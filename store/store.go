package store

import (
	"sync"

	"github.com/danieldzansi/amalitech/models"
)

type IdempotencyStore struct {
	mu    sync.RWMutex
	cache map[string]*models.PaymentResponse
}

func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{
		cache: make(map[string]*models.PaymentResponse),
	}
}

func (s *IdempotencyStore) Get(key string) (*models.PaymentResponse, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	resp, found := s.cache[key]
	return resp, found
}

func (s *IdempotencyStore) Set(key string, resp *models.PaymentResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache[key] = resp
}
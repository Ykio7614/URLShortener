package repository

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("not found")

type MemoryRepo struct {
	mu    sync.RWMutex
	store map[string]string // map[shortURL]originalURL
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		store: make(map[string]string),
	}
}

func (r *MemoryRepo) Save(shortURL, originalURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[shortURL] = originalURL
	return nil
}

func (r *MemoryRepo) Get(shortURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	originalURL, exists := r.store[shortURL]
	if !exists {
		return "", ErrNotFound
	}
	return originalURL, nil
}

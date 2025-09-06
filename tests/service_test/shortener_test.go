package service_test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/Ykio7614/URLShortener/internal/repository"
	. "github.com/Ykio7614/URLShortener/internal/service"
)

type mockRepository struct {
	data map[string]string
}

func (m *mockRepository) Save(shortURL, originalURL string) error {
	m.data[shortURL] = originalURL
	return nil
}

func (m *mockRepository) Get(shortURL string, logger *slog.Logger) (string, error) {
	originalURL, exists := m.data[shortURL]
	if !exists {
		return "", repository.ErrNotFound
	}
	return originalURL, nil
}

func TestShorterAndResolv(t *testing.T) {
	repo := &mockRepository{make(map[string]string)}
	svc := NewShortenerService(repo)

	_, err := svc.GetOriginalURL("missing", nil)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

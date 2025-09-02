package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/Ykio7614/URLShortener/internal/repository"
)

type ShortenerService struct {
	repo *repository.MemoryRepo
}

func NewShortenerService(repo *repository.MemoryRepo) *ShortenerService {
	return &ShortenerService{
		repo: repo,
	}
}

func (s *ShortenerService) ShortenURL(originalURL string, logger *slog.Logger) (string, error) {
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		logger.Error("invalid URL", slog.String("url", originalURL), slog.Any("error", err))
		return "", errors.New("invalid URL")
	}

	shortURL := generateShortURL(originalURL)

	s.repo.Save(shortURL, originalURL)

	return shortURL, nil
}

func generateShortURL(originalURL string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(originalURL)))[:6]
}

func (s *ShortenerService) GetOriginalURL(shortURL string, logger *slog.Logger) (string, error) {
	return s.repo.Get(shortURL, logger)
}

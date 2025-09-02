package service

import (
	"crypto/md5"
	"errors"
	"fmt"
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

func (s *ShortenerService) ShortenURL(originalURL string) (string, error) {
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return "", errors.New("invalid URL")
	}

	shortURL := generateShortURL(originalURL)

	s.repo.Save(shortURL, originalURL)

	return shortURL, nil
}

func generateShortURL(originalURL string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(originalURL)))[:6]
}

func (s *ShortenerService) GetOriginalURL(shortURL string) (string, error) {
	return s.repo.Get(shortURL)
}

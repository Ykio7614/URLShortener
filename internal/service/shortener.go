package service

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"log/slog"
	"net/url"

	"github.com/Ykio7614/URLShortener/internal/repository"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Shortener interface {
	ShortenURL(originalURL string, logger *slog.Logger) (string, error)
	GetOriginalURL(shortURL string, logger *slog.Logger) (string, error)
}

type ShortenerService struct {
	repo repository.Repo
}

func NewShortenerService(r repository.Repo) *ShortenerService {
	return &ShortenerService{
		repo: r,
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
	md5hash := md5.Sum([]byte(originalURL))

	num := binary.BigEndian.Uint64(md5hash[:8])

	shortURL := Base62Encode(int64(num))

	if len(shortURL) > 6 {
		shortURL = shortURL[:6]
	}

	return shortURL
}

func Base62Encode(num int64) string {
	if num == 0 {
		return string(alphabet[0])
	}

	result := ""
	base := int64(len(alphabet))

	for num > 0 {
		remainder := num % base
		result = string(alphabet[remainder]) + result
		num = num / base
	}

	return result
}

func (s *ShortenerService) GetOriginalURL(shortURL string, logger *slog.Logger) (string, error) {
	return s.repo.Get(shortURL, logger)
}

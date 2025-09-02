package repository

import (
	"errors"
	"log/slog"
)

var ErrNotFound = errors.New("not found")

type Repo interface {
	Save(shortURL, originalURL string) error
	Get(shortURL string, logger *slog.Logger) (string, error)
}

package repository

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (r *PostgresRepo) Save(shortURL, originalURL string) error {
	_, err := r.db.Exec("INSERT INTO links (short, original, created_at) VALUES ($1, $2, $3) ON CONFLICT (short) DO NOTHING", shortURL, originalURL, time.Now())
	return err
}

func (r *PostgresRepo) Get(shortURL string, logger *slog.Logger) (string, error) {
	var originalURL string
	err := r.db.QueryRow("SELECT original FROM links WHERE short = $1", shortURL).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	return originalURL, nil
}

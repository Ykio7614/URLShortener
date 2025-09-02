package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ykio7614/URLShortener/internal/httpapi"
	"github.com/Ykio7614/URLShortener/internal/platform"
	"github.com/Ykio7614/URLShortener/internal/repository"
	"github.com/Ykio7614/URLShortener/internal/service"
	_ "github.com/lib/pq"
)

func main() {
	logger := platform.NewLogger()
	port := getenv("PORT", "9090")

	dsn := "postgres://postgres:xxXX1234@localhost:5432/shortener?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Error("failed to ping database", slog.Any("error", err))
		return
	}

	repo := repository.NewMemoryRepo()
	svc := service.NewShortenerService(repo)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           httpapi.Router(logger, svc),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("starting server", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", slog.Any("error", err.Error()))
		}
	}()

	<-ctx.Done()
	logger.Info("shutting signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", slog.Any("error", err.Error()))
	} else {
		logger.Info("server stopped")
	}

}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ykio7614/URLShortener/internal/config"
	"github.com/Ykio7614/URLShortener/internal/httpapi"
	"github.com/Ykio7614/URLShortener/internal/platform"
	"github.com/Ykio7614/URLShortener/internal/repository"
	"github.com/Ykio7614/URLShortener/internal/service"
	"github.com/golang-migrate/migrate/v4"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {

	conf := config.LoadConfig()
	logger := platform.NewLogger()
	port := conf.Port
	var repo repository.Repo
	logger.Info("Config", slog.String("ENV", conf.Env))

	if conf.Env == "test" {
		repo = repository.NewMemoryRepo()
	} else {
		dsn := config.DSNbuilder(conf)
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

		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			logger.Error("failed to create migrate db instance", slog.Any("error", err))
			return
		}

		m, err := migrate.NewWithDatabaseInstance(
			"file://db/migrations",
			"postgres", driver)
		if err != nil {
			logger.Error("failed to create migrate instance", slog.Any("error", err))
			return
		}

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			logger.Error("failed to run migrations", slog.Any("error", err))
			return
		}

		repo = repository.NewPostgresRepo(db)
	}

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

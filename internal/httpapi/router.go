package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/Ykio7614/URLShortener/internal/service"
)

func Router(logger *slog.Logger, svc *service.ShortenerService) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"status":    "ok",
			"ts":        time.Now().Format(time.RFC3339),
			"userAgent": r.UserAgent(),
		})
	})

	mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"ready": "true",
		})
	})

	mux.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var payload struct {
			URL string `json:"url"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		short, err := svc.ShortenURL(payload.URL, logger)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error": err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"short_url": short,
		})
	})

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		shortURL := r.URL.Path[1:]
		if shortURL == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		originalURL, err := svc.GetOriginalURL(shortURL, logger)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			writeJSON(w, http.StatusNotFound, map[string]any{
				"error": "not found",
			})
			return
		}

		http.Redirect(w, r, originalURL, http.StatusFound)
	}))

	return withLogging(logger, mux)
}

func withLogging(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Info("request",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.String("remoteAddr", r.RemoteAddr),
			slog.String("userAgent", r.UserAgent()),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

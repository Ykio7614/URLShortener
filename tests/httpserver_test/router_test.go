package httpserver_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Ykio7614/URLShortener/internal/httpapi"
	"github.com/Ykio7614/URLShortener/internal/repository"
	"github.com/Ykio7614/URLShortener/internal/service"
	"github.com/stretchr/testify/require"
)

func TestShortenerAndResolv(t *testing.T) {
	repo := repository.NewMemoryRepo()
	svc := service.NewShortenerService(repo)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	srv := httptest.NewServer(httpapi.Router(logger, svc))
	defer srv.Close()

	payload := map[string]string{
		"url": "https://golang.org",
	}
	originalURL, err := json.Marshal(payload)
	require.NoError(t, err)
	resp, err := http.Post(srv.URL+"/shorten", "aplicateion/json", bytes.NewReader(originalURL))
	require.NoError(t, err)
	defer resp.Body.Close()

	var res map[string]string
	err = json.NewDecoder(resp.Body).Decode(&res)
	require.NoError(t, err)
	short := res["short_url"]
	require.NotEmpty(t, short)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp2, err := client.Get(srv.URL + "/" + short)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, resp2.StatusCode)
	require.Equal(t, "https://golang.org", resp2.Header.Get("Location"))
}

package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/presswintox/practicum-shortener/internal/config"
	"github.com/presswintox/practicum-shortener/internal/repository"
	"github.com/presswintox/practicum-shortener/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_DoShortUrlHandler(t *testing.T) {

	type want struct {
		code        int
		contentType string
		response    string
	}
	tests := []struct {
		name    string
		url     string
		want    want
		request string
	}{
		{
			name: "success",
			url:  "http://google.com",
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain; charset=UTF-8",
				response:    "http://localhost:8080/x7kg9X5V",
			},
			request: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Server: &config.ServerConfig{
					Port: "8080",
				},
				ShorterService: &config.ShorterServiceConfig{
					ShortUrlAddr: "http://localhost:8080",
				},
			}
			shorterService := service.NewShorterService(repository.NewMemoryRepository(), cfg.ShorterService.ShortUrlAddr)
			api := NewShorterApi(shorterService)

			e := echo.New()
			request := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.url))
			w := httptest.NewRecorder()
			c := e.NewContext(request, w)
			assert.NoError(t, api.DoShortUrlHandler(c))
			result := w.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			resultBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			require.NoError(t, result.Body.Close())

			assert.NotEqual(t, tt.url, string(resultBody))

		})
	}
}

func TestServer_GetUrlHandler(t *testing.T) {

	type want struct {
		code     int
		location string
	}
	tests := []struct {
		name string
		id   string
		want want
	}{
		{
			name: "success",
			id:   "",
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: "https://google.com",
			},
		},
		{
			name: "not found",
			id:   "unknown",
			want: want{
				code:     http.StatusBadRequest,
				location: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Server: &config.ServerConfig{
					Port: "8080",
				},
				ShorterService: &config.ShorterServiceConfig{
					ShortUrlAddr: "http://localhost:8080",
				},
			}
			shorterService := service.NewShorterService(repository.NewMemoryRepository(), cfg.ShorterService.ShortUrlAddr)

			api := NewShorterApi(shorterService)
			originalUrl := "https://google.com"
			shortId, _, err := shorterService.DoShortUrl(originalUrl)
			require.NoError(t, err)

			if tt.id == "unknown" {
				shortId = tt.id
			}

			e := echo.New()

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			c := e.NewContext(request, w)
			c.SetPathValues(echo.PathValues{{Name: "id", Value: shortId}})
			assert.NoError(t, api.GetUrlHandler(c))

			result := w.Result()
			require.NoError(t, result.Body.Close())

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}

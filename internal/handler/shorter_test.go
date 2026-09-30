package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
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
		{
			name: "empty url",
			url:  "",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
			},
			request: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				ShorterService: &config.ShorterServiceConfig{
					ShortURLAddr: "http://localhost:8080",
				},
			}
			shorterService := service.NewShorterService(repository.NewMemoryRepository(), cfg.ShorterService.ShortURLAddr)
			api := NewShorterAPI(shorterService)

			w := echotest.ContextConfig{
				Request: httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.url)),
			}.ServeWithHandler(t, api.DoShortURLHandler)

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
				ShorterService: &config.ShorterServiceConfig{
					ShortURLAddr: "http://localhost:8080",
				},
			}
			shorterService := service.NewShorterService(repository.NewMemoryRepository(), cfg.ShorterService.ShortURLAddr)

			api := NewShorterAPI(shorterService)
			originalURL := "https://google.com"
			shortID, _, err := shorterService.DoShortURL(originalURL)
			require.NoError(t, err)

			if tt.id == "unknown" {
				shortID = tt.id
			}

			w := echotest.ContextConfig{
				Request:    httptest.NewRequest(http.MethodGet, "/", nil),
				PathValues: echo.PathValues{{Name: "id", Value: shortID}},
			}.ServeWithHandler(t, api.GetURLHandler)

			result := w.Result()
			require.NoError(t, result.Body.Close())

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}

func TestServer_ShortenHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
		response    ShortenResponse
	}
	tests := []struct {
		name        string
		requestBody string
		want        want
	}{
		{
			name:        "success",
			requestBody: `{"url":"http://google.com"}`,
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
				response:    ShortenResponse{Result: "http://localhost:8080/x7kg9X5V"},
			},
		},
		{
			name:        "empty url",
			requestBody: `{"url":""}`,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
			},
		},
		{
			name:        "invalid json",
			requestBody: `{"url":`,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shorterService := &shorterServiceStub{
				shortURL: "http://localhost:8080/x7kg9X5V",
			}
			api := NewShorterAPI(shorterService)

			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(tt.requestBody))
			request.Header.Set("Content-Type", "application/json")
			w := echotest.ContextConfig{Request: request}.ServeWithHandler(t, api.ShortenHandler)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if tt.want.code == http.StatusCreated {
				var response ShortenResponse
				require.NoError(t, json.NewDecoder(result.Body).Decode(&response))
				assert.Equal(t, tt.want.response, response)
			}
		})
	}
}

type shorterServiceStub struct {
	shortURL string
}

func (s *shorterServiceStub) DoShortURL(url string) (string, string, error) {
	return "x7kg9X5V", s.shortURL, nil
}

func (s *shorterServiceStub) GetURL(string) (string, error) {
	return "", nil
}

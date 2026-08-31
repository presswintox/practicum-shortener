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

var cfg = &config.Config{
	Port:         "8080",
	ShortUrlAddr: "http://localhost:8080",
}

func TestServer_DoShortUrlHandler(t *testing.T) {

	shorterService := service.NewShorterService(repository.NewMemoryRepository(), cfg)
	server := NewServer(shorterService)
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
			e := echo.New()
			request := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.url))
			w := httptest.NewRecorder()
			c := e.NewContext(request, w)
			assert.NoError(t, server.DoShortUrlHandler(c))
			result := w.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			resultBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			require.NoError(t, result.Body.Close())

			assert.Equal(t, tt.want.response, string(resultBody))

		})
	}
}

func TestServer_GetUrlHandler(t *testing.T) {
	shorterService := service.NewShorterService(repository.NewMemoryRepository(), cfg)
	server := NewServer(shorterService)

	originalUrl := "http://google.com"
	shortId, _, err := shorterService.DoShortUrl(originalUrl)
	require.NoError(t, err)

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
			id:   shortId,
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: originalUrl,
			},
		},
		{
			name: "not found",
			id:   "unknown1",
			want: want{
				code:     http.StatusBadRequest,
				location: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			c := e.NewContext(request, w)
			c.SetPathValues(echo.PathValues{{Name: "id", Value: tt.id}})
			assert.NoError(t, server.GetUrlHandler(c))

			result := w.Result()
			require.NoError(t, result.Body.Close())

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}

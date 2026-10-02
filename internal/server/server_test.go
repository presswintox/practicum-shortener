package server

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/presswintox/practicum-shortener/internal/handler"
	"github.com/presswintox/practicum-shortener/internal/repository"
	"github.com/presswintox/practicum-shortener/internal/service"
	"github.com/stretchr/testify/require"
)

func TestGzipCompression(t *testing.T) {
	const shortURLAddr = "http://localhost:8080"

	shorterService := service.NewShorterService(repository.NewMemoryRepository(), shortURLAddr)
	api := handler.NewShorterAPI(shorterService)
	s := NewServer(":0", api)
	requestBody := `{"url":"https://example.com"}`

	t.Run("accepts gzip request", func(t *testing.T) {
		var compressed bytes.Buffer
		zw := gzip.NewWriter(&compressed)
		_, err := zw.Write([]byte(requestBody))
		require.NoError(t, err)
		require.NoError(t, zw.Close())

		req := httptest.NewRequest(http.MethodPost, "/api/shorten", &compressed)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderContentEncoding, "gzip")
		res := httptest.NewRecorder()

		s.echo.ServeHTTP(res, req)

		require.Equal(t, http.StatusCreated, res.Code)

		var response handler.ShortenResponse
		require.NoError(t, json.Unmarshal(res.Body.Bytes(), &response))
		require.Contains(t, response.Result, shortURLAddr)
	})

	t.Run("sends gzip response", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAcceptEncoding, "gzip")
		res := httptest.NewRecorder()

		s.echo.ServeHTTP(res, req)

		require.Equal(t, http.StatusCreated, res.Code)
		require.Equal(t, "gzip", res.Header().Get(echo.HeaderContentEncoding))

		zr, err := gzip.NewReader(res.Body)
		require.NoError(t, err)
		defer zr.Close()

		body, err := io.ReadAll(zr)
		require.NoError(t, err)

		var response handler.ShortenResponse
		require.NoError(t, json.Unmarshal(body, &response))
		require.Contains(t, response.Result, shortURLAddr)
	})
}

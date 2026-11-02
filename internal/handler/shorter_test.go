package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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
			shortURL, err := shorterService.Shorten(originalURL)
			require.NoError(t, err)
			parsedURL, err := url.Parse(shortURL)
			require.NoError(t, err)
			shortID := strings.TrimPrefix(parsedURL.Path, "/")

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
	shortURL     string
	shortenBatch func(map[string]string) (map[string]string, error)
}

func (s *shorterServiceStub) Shorten(string) (string, error) {
	return s.shortURL, nil
}

func (s *shorterServiceStub) GetURL(string) (string, error) {
	return "", nil
}

func (s *shorterServiceStub) ShortenBatch(urls map[string]string) (map[string]string, error) {
	return s.shortenBatch(urls)
}

func TestServer_ShortenBatchHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		serviceErr error
		wantCode   int
		wantCalls  int
	}{
		{name: "success", body: `[{"correlation_id":"first","original_url":"https://example.com/1"},{"correlation_id":"second","original_url":"https://example.com/2"}]`, wantCode: http.StatusCreated, wantCalls: 1},
		{name: "service error", body: `[{"correlation_id":"first","original_url":"https://example.com/1"},{"correlation_id":"second","original_url":"https://example.com/2"}]`, serviceErr: errors.New("storage unavailable"), wantCode: http.StatusInternalServerError, wantCalls: 1},
		{name: "invalid JSON", body: `[{"correlation_id":`, wantCode: http.StatusBadRequest},
		{name: "wrong JSON shape", body: `{"correlation_id":"first"}`, wantCode: http.StatusBadRequest},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			stub := &shorterServiceStub{shortenBatch: func(urls map[string]string) (map[string]string, error) {
				calls++
				assert.Equal(t, map[string]string{"first": "https://example.com/1", "second": "https://example.com/2"}, urls)
				return map[string]string{"first": "http://localhost:8080/first123", "second": "http://localhost:8080/second12"}, test.serviceErr
			}}
			api := NewShorterAPI(stub)
			request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(test.body))
			request.Header.Set("Content-Type", "application/json")
			w := echotest.ContextConfig{Request: request}.ServeWithHandler(t, api.ShortenBatchHandler)

			assert.Equal(t, test.wantCode, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, test.wantCalls, calls)
			if test.wantCode == http.StatusCreated {
				var response []ShortenBatchItemResponse
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
				assert.ElementsMatch(t, []ShortenBatchItemResponse{
					{CorrelationID: "first", ShortURL: "http://localhost:8080/first123"},
					{CorrelationID: "second", ShortURL: "http://localhost:8080/second12"},
				}, response)
			}
		})
	}
}

func TestServer_ShortenBatchHandlerStoresURLs(t *testing.T) {
	repo := repository.NewMemoryRepository()
	shorter := service.NewShorterService(repo, "http://localhost:8080")
	api := NewShorterAPI(shorter)
	request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(`[{"correlation_id":"first","original_url":"https://example.com/1"},{"correlation_id":"second","original_url":"https://example.com/2"}]`))
	request.Header.Set("Content-Type", "application/json")
	w := echotest.ContextConfig{Request: request}.ServeWithHandler(t, api.ShortenBatchHandler)
	require.Equal(t, http.StatusCreated, w.Code)
	var response []ShortenBatchItemResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	require.Len(t, response, 2)
	want := map[string]string{"first": "https://example.com/1", "second": "https://example.com/2"}
	for _, item := range response {
		require.Contains(t, want, item.CorrelationID)
		parsed, err := url.Parse(item.ShortURL)
		require.NoError(t, err)
		original, err := shorter.GetURL(strings.TrimPrefix(parsed.Path, "/"))
		require.NoError(t, err)
		assert.Equal(t, want[item.CorrelationID], original)
		delete(want, item.CorrelationID)
	}
	assert.Empty(t, want)
}

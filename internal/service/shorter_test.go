package service

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/presswintox/practicum-shortener/internal/config"
	"github.com/presswintox/practicum-shortener/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type batchRepositoryStub struct {
	ShorterRepository
	batchSave func(map[string]string) error
}

func (r *batchRepositoryStub) BatchSave(urls map[string]string) error {
	return r.batchSave(urls)
}

func TestShorterService_ShortenBatch(t *testing.T) {
	storageErr := errors.New("storage unavailable")
	tests := []struct {
		name      string
		failures  int
		saveErr   error
		wantErr   error
		wantCalls int
	}{
		{name: "success", wantCalls: 1},
		{name: "retry after collision", failures: 1, saveErr: fmt.Errorf("duplicate: %w", repository.ErrAlreadyExists), wantCalls: 2},
		{name: "collision limit", failures: maxHashAttempts, saveErr: repository.ErrAlreadyExists, wantErr: ErrHashCollision, wantCalls: maxHashAttempts},
		{name: "storage error", failures: 1, saveErr: storageErr, wantErr: storageErr, wantCalls: 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			urls := map[string]string{"first": "https://example.com/1", "second": "https://example.com/2", "same-url": "https://example.com/1"}
			memory := repository.NewMemoryRepository()
			calls := 0
			repo := &batchRepositoryStub{ShorterRepository: memory, batchSave: func(batch map[string]string) error {
				calls++
				require.Len(t, batch, 2)
				values := make([]string, 0, len(batch))
				for hash, original := range batch {
					assert.Len(t, hash, length)
					values = append(values, original)
				}
				assert.ElementsMatch(t, []string{urls["first"], urls["second"]}, values)
				if calls <= test.failures {
					return test.saveErr
				}
				return memory.BatchSave(batch)
			}}
			s := NewShorterService(repo, "http://localhost:8080")
			result, err := s.ShortenBatch(urls)
			assert.Equal(t, test.wantCalls, calls)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				assert.Nil(t, result)
				return
			}
			require.NoError(t, err)
			require.Len(t, result, len(urls))
			assert.Equal(t, result["first"], result["same-url"])
			for correlationID, original := range urls {
				shortURL, err := url.Parse(result[correlationID])
				require.NoError(t, err)
				assert.Equal(t, "http", shortURL.Scheme)
				assert.Equal(t, "localhost:8080", shortURL.Host)
				saved, err := memory.Get(strings.TrimPrefix(shortURL.Path, "/"))
				require.NoError(t, err)
				assert.Equal(t, original, saved)
			}
		})
	}
}

func TestShorterService_ShortenBatchInvalidBaseURL(t *testing.T) {
	calls := 0
	repo := &batchRepositoryStub{ShorterRepository: repository.NewMemoryRepository(), batchSave: func(map[string]string) error {
		calls++
		return nil
	}}
	s := NewShorterService(repo, "://invalid")
	result, err := s.ShortenBatch(map[string]string{"first": "https://example.com"})
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Zero(t, calls)
}

func TestShorterService_ShortenBatchEmpty(t *testing.T) {
	for _, urls := range []map[string]string{nil, {}} {
		s := NewShorterService(repository.NewMemoryRepository(), "http://localhost:8080")
		result, err := s.ShortenBatch(urls)
		require.NoError(t, err)
		assert.Empty(t, result)
	}
}

func TestShorterService_DoShortUrl(t *testing.T) {

	tests := []struct {
		name string
		url  string
	}{
		{
			name: "simple test 1",
			url:  "http://google.com",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := &config.Config{
				ShorterService: &config.ShorterServiceConfig{ShortURLAddr: "http://localhost:8080"},
			}
			s := NewShorterService(repository.NewMemoryRepository(), cfg.ShorterService.ShortURLAddr)

			shortURL, err := s.Shorten(test.url)
			require.NoError(t, err)
			assert.NotEqual(t, test.url, shortURL)
		})
	}
}

func TestShortenURL(t *testing.T) {

	tests := []struct {
		name string
		url  string
	}{
		{
			name: "simple test 1",
			url:  "http://google.com",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.NotEqual(t, test.url, urlHash(test.url))
		})
	}
}

func TestShortenURL_NotEqualSameUrl(t *testing.T) {

	tests := []struct {
		name string
		url  string
	}{
		{
			name: "simple test 1",
			url:  "http://google.com",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.NotEqual(t, urlHash(test.url), urlHash(test.url))
		})
	}
}

func TestShorterService_GetUrl(t *testing.T) {

	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "simple test 1",
			url:  "https://google.com",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := &config.Config{
				ShorterService: &config.ShorterServiceConfig{ShortURLAddr: "http://localhost:8080"},
			}
			s := NewShorterService(repository.NewMemoryRepository(), cfg.ShorterService.ShortURLAddr)

			shortURL, err := s.Shorten(test.url)
			require.NoError(t, err)

			parsedURL, err := url.Parse(shortURL)
			require.NoError(t, err)
			hash := strings.TrimPrefix(parsedURL.Path, "/")

			url, err := s.GetURL(hash)
			require.NoError(t, err)
			assert.Equal(t, test.url, url)
		})
	}
}

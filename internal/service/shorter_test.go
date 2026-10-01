package service

import (
	"net/url"
	"strings"
	"testing"

	"github.com/presswintox/practicum-shortener/internal/config"
	"github.com/presswintox/practicum-shortener/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

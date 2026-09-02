package service

import (
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
				Server:         &config.ServerConfig{Port: ":8080"},
				ShorterService: &config.ShorterServiceConfig{ShortUrlAddr: "http://localhost:8080"},
			}
			s := NewShorterService(repository.NewMemoryRepository(), cfg.ShorterService.ShortUrlAddr)

			_, shortUrl, err := s.DoShortUrl(test.url)
			require.NoError(t, err)
			assert.NotEqual(t, test.url, shortUrl)
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
				Server:         &config.ServerConfig{Port: ":8080"},
				ShorterService: &config.ShorterServiceConfig{ShortUrlAddr: "http://localhost:8080"},
			}
			s := NewShorterService(repository.NewMemoryRepository(), cfg.ShorterService.ShortUrlAddr)

			hash, _, err := s.DoShortUrl(test.url)
			require.NoError(t, err)
			url, err := s.GetUrl(hash)
			require.NoError(t, err)
			assert.Equal(t, test.url, url)
		})
	}
}

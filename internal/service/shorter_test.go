package service

import (
	"testing"

	"github.com/presswintox/practicum-shortener/internal/config"
	"github.com/presswintox/practicum-shortener/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cfg = &config.Config{
	Port:         "8080",
	ShortUrlAddr: "http://localhost:8080",
}

func TestShorterService_DoShortUrl(t *testing.T) {

	s := NewShorterService(repository.NewMemoryRepository(), cfg)

	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "simple test 1",
			url:  "http://google.com",
			want: "http://localhost:8080/x7kg9X5V",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, shortUrl, err := s.DoShortUrl(test.url)
			require.NoError(t, err)
			assert.Equal(t, test.want, shortUrl)
		})
	}
}

func TestShortenURL(t *testing.T) {

	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "simple test 1",
			url:  "http://google.com",
			want: "x7kg9X5V",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Run(test.name, func(t *testing.T) {
				assert.Equal(t, test.want, URLHash(test.url))
			})
		})
	}
}

func TestShorterService_GetUrl(t *testing.T) {
	s := NewShorterService(repository.NewMemoryRepository(), cfg)

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
			hash, _, err := s.DoShortUrl(test.url)
			require.NoError(t, err)
			url, err := s.GetUrl(hash)
			require.NoError(t, err)
			assert.Equal(t, test.url, url)
		})
	}
}

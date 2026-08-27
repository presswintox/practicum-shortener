package service

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
)

const length int = 8

type ShorterService struct {
	mu sync.Mutex
	db map[string]string
}

func NewShorterService() *ShorterService {
	return &ShorterService{db: make(map[string]string)}
}

func (s *ShorterService) DoShortUrl(url string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	shortUrl := ShortenURL(url)
	s.db[shortUrl] = url
	return fmt.Sprintf("http://localhost:8080/%s", shortUrl)
}

func (s *ShorterService) GetUrl(shortUrl string) (string, error) {
	if url, ok := s.db[shortUrl]; ok {
		return url, nil
	}
	return "", errors.New("shorter url not found")
}

func ShortenURL(longURL string) string {
	urlLength := length
	hash := md5.Sum([]byte(longURL))
	// base64 URL-safe без паддинга, чтобы не было / + =
	encoded := base64.RawURLEncoding.EncodeToString(hash[:])
	if urlLength > len(encoded) {
		urlLength = len(encoded)
	}
	return encoded[:urlLength]
}

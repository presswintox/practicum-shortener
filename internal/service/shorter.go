package service

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"

	"github.com/presswintox/practicum-shortener/internal/repository"
)

const length int = 8

type ShorterService struct {
	db repository.ShorterRepository
}

func NewShorterService(db repository.ShorterRepository) *ShorterService {
	return &ShorterService{db: db}
}

func (s *ShorterService) DoShortUrl(url string) (string, string, error) {
	urlHash := URLHash(url)
	if err := s.db.Save(urlHash, url); err != nil {
		return "", "", err
	}
	return urlHash, fmt.Sprintf("http://localhost:8080/%s", urlHash), nil
}

func (s *ShorterService) GetUrl(shortUrl string) (string, error) {
	return s.db.Get(shortUrl)
}

func URLHash(longURL string) string {
	urlLength := length
	hash := md5.Sum([]byte(longURL))
	// base64 URL-safe без паддинга, чтобы не было / + =
	encoded := base64.RawURLEncoding.EncodeToString(hash[:])
	if urlLength > len(encoded) {
		urlLength = len(encoded)
	}
	return encoded[:urlLength]
}

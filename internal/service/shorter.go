package service

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/url"

	"github.com/presswintox/practicum-shortener/internal/repository"
)

const length int = 8

const maxHashAttempts int = 10

const saltLength int = 16

var ErrHashCollision = errors.New("failed to generate unique short url")

type ShorterRepository interface {
	Save(shortUrl, url string) error
	Get(shortUrl string) (string, error)
}
type ShorterService struct {
	db           ShorterRepository
	shortUrlAddr string
}

func NewShorterService(db ShorterRepository, shortUrlAddr string) *ShorterService {
	return &ShorterService{db: db, shortUrlAddr: shortUrlAddr}
}

func (s *ShorterService) DoShortUrl(url string) (string, string, error) {
	for range maxHashAttempts {
		hash := urlHash(url)

		err := s.db.Save(hash, url)
		switch {
		case err == nil:
			shortUrl, err2 := s.shortUrl(hash)
			if err2 != nil {
				return "", "", fmt.Errorf("failed to generate short url: %w", err2)
			}
			return hash, shortUrl, nil
		case errors.Is(err, repository.ErrAlreadyExists):
			continue
		default:
			return "", "", fmt.Errorf("failed to save short url: %w", err)
		}
	}
	return "", "", fmt.Errorf("failed to save short url: %w", ErrHashCollision)
}

func (s *ShorterService) GetUrl(shortUrl string) (string, error) {
	return s.db.Get(shortUrl)
}

func (s *ShorterService) shortUrl(hash string) (string, error) {
	return url.JoinPath(s.shortUrlAddr, hash)
}

func generateSalt(length int) string {
	const letterRunes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, length)
	for i := range b {
		b[i] = letterRunes[rand.IntN(len(letterRunes))]
	}
	return string(b)
}

func urlHash(longURL string) string {
	urlLength := length
	payload := longURL + generateSalt(saltLength)

	hash := md5.Sum([]byte(payload))
	// base64 URL-safe без паддинга, чтобы не было / + =
	encoded := base64.RawURLEncoding.EncodeToString(hash[:])
	if urlLength > len(encoded) {
		urlLength = len(encoded)
	}
	return encoded[:urlLength]
}

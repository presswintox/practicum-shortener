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
	Save(shortURL, url string) error
	Get(shortURL string) (string, error)
	GetByOriginalURL(originalURL string) (string, error)
	BatchSave(shortURLs map[string]string) error
}
type ShorterService struct {
	db           ShorterRepository
	shortURLAddr string
}

func NewShorterService(db ShorterRepository, shortURLAddr string) *ShorterService {
	return &ShorterService{db: db, shortURLAddr: shortURLAddr}
}

func (s *ShorterService) Shorten(url string) (string, error) {
	for range maxHashAttempts {
		hash := urlHash(url)

		err := s.db.Save(hash, url)
		switch {
		case err == nil:
			shortURL, err2 := s.shortURL(hash)
			if err2 != nil {
				return "", fmt.Errorf("failed to generate short url: %w", err2)
			}
			return shortURL, nil
		case errors.Is(err, repository.ErrOriginalURLExists):
			existing, lookupErr := s.db.GetByOriginalURL(url)
			if lookupErr != nil {
				return "", fmt.Errorf("failed to find existing short url: %w", lookupErr)
			}
			shortURL, buildErr := s.shortURL(existing)
			if buildErr != nil {
				return "", fmt.Errorf("failed to generate short url: %w", buildErr)
			}
			return shortURL, repository.ErrOriginalURLExists
		case errors.Is(err, repository.ErrAlreadyExists):
			continue
		default:
			return "", fmt.Errorf("failed to save short url: %w", err)
		}
	}
	return "", fmt.Errorf("failed to save short url: %w", ErrHashCollision)
}

func (s *ShorterService) GetURL(shortURL string) (string, error) {
	return s.db.Get(shortURL)
}

func (s *ShorterService) shortURL(hash string) (string, error) {
	return url.JoinPath(s.shortURLAddr, hash)
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

type Url struct {
	corrID string
	hash   string
	url    string
}

func (s *ShorterService) ShortenBatch(urls map[string]string) (map[string]string, error) {
	for range maxHashAttempts {
		urlsToSave := make(map[string]string, len(urls))
		result := make(map[string]string, len(urls))
		for corrID, url := range urls {
			hash := urlHash(url)
			urlsToSave[hash] = url
			shortURL, err := s.shortURL(hash)
			if err != nil {
				return nil, fmt.Errorf("failed to generate short url: %w", err)
			}
			result[corrID] = shortURL
		}
		if len(urlsToSave) != len(urls) {
			continue
		}
		err := s.db.BatchSave(urlsToSave)
		switch {
		case err == nil:
			return result, nil
		case errors.Is(err, repository.ErrAlreadyExists):
			continue
		default:
			return nil, fmt.Errorf("failed to save short url: %w", err)
		}
	}
	return nil, fmt.Errorf("failed to save short url: %w", ErrHashCollision)
}

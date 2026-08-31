package repository

import "errors"

var ErrNotFound = errors.New("shorter url not found")

type ShorterRepository interface {
	Save(shortUrl, url string) error
	Get(shortUrl string) (string, error)
}

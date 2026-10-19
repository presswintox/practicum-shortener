package repository

import "errors"

// Общие ошибки для repository package

var ErrNotFound = errors.New("shorter url not found")
var ErrAlreadyExists = errors.New("shorter url already exists")

type URLRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

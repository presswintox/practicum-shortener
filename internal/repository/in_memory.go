package repository

import (
	"fmt"
	"sync"
)

type MemoryRepository struct {
	mu              sync.Mutex
	db              map[string]string
	shortByOriginal map[string]string
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		db:              make(map[string]string),
		shortByOriginal: make(map[string]string),
	}
}

func (r *MemoryRepository) Save(id, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.shortByOriginal[value]; exists {
		return ErrOriginalURLExists
	}
	if _, exist := r.db[id]; exist {
		return fmt.Errorf("%w: %q", ErrAlreadyExists, id)
	}
	r.db[id] = value
	r.shortByOriginal[value] = id
	return nil
}

func (r *MemoryRepository) Get(id string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	url, ok := r.db[id]
	if !ok {
		return "", fmt.Errorf("%w: %q", ErrNotFound, id)
	}
	return url, nil
}

func (r *MemoryRepository) BatchSave(urls map[string]string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	addedIDs := make([]string, 0, len(urls))
	for id, url := range urls {
		if _, exists := r.db[id]; exists {
			for _, addedID := range addedIDs {
				delete(r.db, addedID)
			}
			return fmt.Errorf("%w: %q", ErrAlreadyExists, id)
		}
		r.db[id] = url
		addedIDs = append(addedIDs, id)
	}
	for id, url := range urls {
		r.shortByOriginal[url] = id
	}
	return nil
}

func (r *MemoryRepository) GetByOriginalURL(originalURL string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if id, ok := r.shortByOriginal[originalURL]; ok {
		return id, nil
	}
	return "", ErrNotFound
}

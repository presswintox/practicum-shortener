package repository

import (
	"sync"
)

type MemoryRepository struct {
	mu sync.RWMutex
	db map[string]string
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{db: make(map[string]string)}
}

func (r *MemoryRepository) Save(id, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.db[id] = value
	return nil
}

func (r *MemoryRepository) Get(id string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, ok := r.db[id]
	if !ok {
		return "", ErrNotFound
	}
	return url, nil
}

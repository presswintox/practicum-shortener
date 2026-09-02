package repository

import (
	"sync"
)

type MemoryRepository struct {
	mu sync.Mutex
	db map[string]string
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{db: make(map[string]string)}
}

func (r *MemoryRepository) Save(id, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exist := r.db[id]; exist {
		return ErrAlreadyExists
	}
	r.db[id] = value
	return nil
}

func (r *MemoryRepository) Get(id string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	url, ok := r.db[id]
	if !ok {
		return "", ErrNotFound
	}
	return url, nil
}

package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/google/uuid"
)

type URLRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type FileRepository struct {
	mu      sync.Mutex
	records map[string]URLRecord
	file    *os.File
}

func NewFileRepository(path string) (*FileRepository, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", path, err)
	}
	db := &FileRepository{file: file, records: make(map[string]URLRecord)}
	if err = db.Load(); err != nil {
		return nil, fmt.Errorf("error loading file %s: %w", path, err)
	}
	return db, nil
}

func (r *FileRepository) Save(id, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.records[id]; ok {
		return fmt.Errorf("%w: %q", ErrAlreadyExists, id)
	}

	record := URLRecord{
		UUID:        uuid.NewString(),
		ShortURL:    id,
		OriginalURL: value,
	}
	if err := r.saveRecord(record); err != nil {
		return fmt.Errorf("failed to save file storage: %w", err)
	}
	r.records[id] = record

	return nil
}

func (r *FileRepository) Get(id string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if record, ok := r.records[id]; ok {
		return record.OriginalURL, nil
	}
	return "", fmt.Errorf("%w: %q", ErrNotFound, id)
}

func (r *FileRepository) Load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, err := r.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	scanner := bufio.NewScanner(r.file)
	for scanner.Scan() {
		var record URLRecord
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return err
		}
		r.records[record.ShortURL] = record
	}
	return scanner.Err()
}

func (r *FileRepository) saveRecord(record URLRecord) error {
	if _, err := r.file.Seek(0, io.SeekEnd); err != nil {
		return err
	}
	if err := json.NewEncoder(r.file).Encode(record); err != nil {
		return err
	}
	return r.file.Sync()
}

func (r *FileRepository) Close() error {
	return r.file.Close()
}

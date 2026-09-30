package repository

import (
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
	records []URLRecord
	file    *os.File
}

func NewFileRepository(file *os.File) *FileRepository {
	return &FileRepository{file: file}
}

func (r *FileRepository) Save(id, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, record := range r.records {
		if record.ShortURL == id {
			return fmt.Errorf("%w: %q", ErrAlreadyExists, id)
		}
	}

	record := URLRecord{
		UUID:        uuid.NewString(),
		ShortURL:    id,
		OriginalURL: value,
	}
	r.records = append(r.records, record)
	if err := r.saveRecords(r.records); err != nil {
		return fmt.Errorf("failed to save file storage: %w", err)
	}

	return nil
}

func (r *FileRepository) Get(id string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, record := range r.records {
		if record.ShortURL == id {
			return record.OriginalURL, nil
		}
	}
	return "", fmt.Errorf("%w: %q", ErrNotFound, id)
}

func (r *FileRepository) Load() error {
	if _, err := r.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	decoder := json.NewDecoder(r.file)
	if err := decoder.Decode(&r.records); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	return nil
}

func (r *FileRepository) saveRecords(records []URLRecord) error {
	if _, err := r.file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if err := r.file.Truncate(0); err != nil {
		return err
	}
	if err := json.NewEncoder(r.file).Encode(records); err != nil {
		return err
	}
	return r.file.Sync()
}

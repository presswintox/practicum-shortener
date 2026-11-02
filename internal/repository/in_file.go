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

func (r *FileRepository) BatchSave(shortURLs map[string]string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	records := make([]URLRecord, 0, len(shortURLs))
	for id, url := range shortURLs {
		if _, ok := r.records[id]; ok {
			return fmt.Errorf("%w: %q", ErrAlreadyExists, id)
		}
		record := URLRecord{
			UUID:        uuid.NewString(),
			ShortURL:    id,
			OriginalURL: url,
		}
		records = append(records, record)
	}
	if err := r.saveRecords(records); err != nil {
		return fmt.Errorf("failed to save file storage: %w", err)
	}
	for _, record := range records {
		r.records[record.ShortURL] = record
	}
	return nil
}

func (r *FileRepository) saveRecord(record URLRecord) error {
	return r.saveRecords([]URLRecord{record})
}

func (r *FileRepository) saveRecords(records []URLRecord) error {
	if len(records) == 0 {
		return nil
	}
	if _, err := r.file.Seek(0, io.SeekEnd); err != nil {
		return err
	}
	encoder := json.NewEncoder(r.file)
	for _, record := range records {
		if err := encoder.Encode(record); err != nil {
			return err
		}
	}
	return r.file.Sync()
}

func (r *FileRepository) Close() error {
	return r.file.Close()
}

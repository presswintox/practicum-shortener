package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/presswintox/practicum-shortener/internal/config/db"
	"github.com/presswintox/practicum-shortener/migrations"
)

type PostgresRepository struct {
	mu  sync.Mutex
	db  *db.Database
	dsn string
}

func NewPostgresRepository(dsn string) (*PostgresRepository, error) {
	postgres, err := db.NewPostgresDatabase(dsn)
	if err != nil {
		return nil, fmt.Errorf("error init postgres: %w", err)
	}

	repo := &PostgresRepository{db: postgres, dsn: dsn}
	return repo, nil
}

func (r *PostgresRepository) Save(id, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.Conn.Exec("INSERT INTO short_links (original_url, short_url) VALUES ($1, $2)", value, id)
	if err != nil {
		return fmt.Errorf("failed to save postgres: %w", err)
	}
	return nil
}

func (r *PostgresRepository) Get(id string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	row := r.db.Conn.QueryRow("SELECT id, original_url, short_url FROM short_links WHERE short_url = $1", id)
	record := URLRecord{}
	err := row.Scan(&record.UUID, &record.OriginalURL, &record.ShortURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%w: %q", ErrNotFound, id)
	}
	if err != nil {
		return "", fmt.Errorf("failed to query row: %w", err)
	}
	return record.OriginalURL, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

func (r *PostgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *PostgresRepository) Migrate() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	source, err := iofs.New(migrations.Files, ".")
	if err != nil {
		return fmt.Errorf("failed to create migrations source: %w", err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", source, r.dsn)
	if err != nil {
		_ = source.Close()
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	migrateErr := migrator.Up()
	migrator.Close()

	if migrateErr != nil && !errors.Is(migrateErr, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", migrateErr)
	}

	return nil
}

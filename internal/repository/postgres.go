package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/presswintox/practicum-shortener/migrations"
)

type PostgresRepository struct {
	db  *sql.DB
	dsn string
}

func NewPostgresRepository(dsn string) (*PostgresRepository, error) {
	postgres, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error init postgres: %w", err)
	}

	repo := &PostgresRepository{db: postgres, dsn: dsn}
	return repo, nil
}

func (r *PostgresRepository) Save(id, value string) error {
	_, err := r.db.Exec("INSERT INTO short_links (original_url, short_url) VALUES ($1, $2)", value, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			if pgErr.ConstraintName == "short_links_original_url_key" {
				return ErrOriginalURLExists
			}
			return fmt.Errorf("%w: %q: %w", ErrAlreadyExists, id, err)
		}
		return fmt.Errorf("failed to save postgres: %w", err)
	}
	return nil
}

func (r *PostgresRepository) BatchSave(urls map[string]string) error {
	if len(urls) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin batch save: %w", err)
	}
	defer tx.Rollback()

	for id, url := range urls {
		_, err := tx.Exec("INSERT INTO short_links (original_url, short_url) VALUES ($1, $2)", url, id)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				return fmt.Errorf("%w: %q: %w", ErrAlreadyExists, id, err)
			}
			return fmt.Errorf("failed to batch save postgres: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit batch save: %w", err)
	}
	return nil
}

func (r *PostgresRepository) Get(id string) (string, error) {
	row := r.db.QueryRow("SELECT id, original_url, short_url FROM short_links WHERE short_url = $1", id)
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

func (r *PostgresRepository) GetByOriginalURL(originalURL string) (string, error) {
	var shortURL string
	err := r.db.QueryRow("SELECT short_url FROM short_links WHERE original_url = $1", originalURL).Scan(&shortURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("failed to find existing short url: %w", err)
	}
	return shortURL, nil
}

func (r *PostgresRepository) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return r.db.PingContext(ctx)
}

func (r *PostgresRepository) Migrate() error {
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

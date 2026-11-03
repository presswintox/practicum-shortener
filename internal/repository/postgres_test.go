package repository

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresTestSuite struct {
	suite.Suite

	ctx       context.Context
	container *postgres.PostgresContainer
}

func (s *PostgresTestSuite) SetupSuite() {
	s.ctx = context.Background()

	postgresContainer, err := postgres.Run(
		s.ctx,
		"postgres:17",
		postgres.BasicWaitStrategies(),
	)
	s.container = postgresContainer
	testcontainers.CleanupContainer(s.T(), postgresContainer)
	require.NoError(s.T(), err)
}

func (s *PostgresTestSuite) TestSaveAndGet() {
	tests := []struct {
		name  string
		id    string
		value string
	}{
		{
			name:  "simple test 1",
			id:    "1",
			value: "value",
		},
		{
			name:  "simple test 2",
			id:    "2",
			value: "abcd1234",
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			t := s.T()
			dsn, err := s.container.ConnectionString(s.ctx, "sslmode=disable")
			require.NoError(t, err)

			repository, err := NewPostgresRepository(dsn)
			require.NoError(t, err)
			t.Cleanup(func() {
				assert.NoError(t, repository.Close())
			})
			require.NoError(t, repository.Migrate())
			require.NoError(t, repository.Ping())
			require.NoError(t, repository.Save(test.id, test.value))

			value, err := repository.Get(test.id)
			require.NoError(t, err)
			assert.Equal(t, test.value, value)
		})
	}
}

func (s *PostgresTestSuite) TestBatchSave() {
	t := s.T()
	dsn, err := s.container.ConnectionString(s.ctx, "sslmode=disable")
	require.NoError(t, err)
	repo, err := NewPostgresRepository(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, repo.Close()) })
	require.NoError(t, repo.Migrate())

	tests := []struct {
		name    string
		urls    map[string]string
		wantErr bool
		err     error
	}{
		{name: "multiple records", urls: map[string]string{"batch-first": "https://example.com/first", "batch-second": "https://example.com/second"}},
		{name: "duplicate key rolls back batch", urls: map[string]string{"batch-existing": "replacement", "batch-new": "https://example.com/new"}, wantErr: true, err: ErrAlreadyExists},
		{name: "duplicate URL rolls back batch", urls: map[string]string{"batch-duplicate": "https://example.com/existing", "batch-new": "https://example.com/new"}, wantErr: true, err: ErrAlreadyExists},
		{name: "invalid value rolls back batch", urls: map[string]string{"batch-invalid": strings.Repeat("x", 1001), "batch-new": "https://example.com/new"}, wantErr: true},
		{name: "empty batch", urls: map[string]string{}},
		{name: "nil batch"},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			t := s.T()
			_, err := repo.db.Exec("TRUNCATE TABLE short_links")
			require.NoError(t, err)
			require.NoError(t, repo.Save("batch-existing", "https://example.com/existing"))

			err = repo.BatchSave(test.urls)
			if test.wantErr {
				require.Error(t, err)
				if test.err != nil {
					require.ErrorIs(t, err, test.err)
				}
			} else {
				require.NoError(t, err)
				for id, want := range test.urls {
					got, err := repo.Get(id)
					require.NoError(t, err)
					assert.Equal(t, want, got)
				}
			}

			var count int
			require.NoError(t, repo.db.QueryRow("SELECT COUNT(*) FROM short_links").Scan(&count))
			wantCount := 1
			if !test.wantErr {
				wantCount += len(test.urls)
			}
			assert.Equal(t, wantCount, count)
			value, err := repo.Get("batch-existing")
			require.NoError(t, err)
			assert.Equal(t, "https://example.com/existing", value)
		})
	}
}

func TestPostgresRepository(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}

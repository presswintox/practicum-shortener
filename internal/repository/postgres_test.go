package repository

import (
	"context"
	"os/exec"
	"path/filepath"
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

	dsn, err := s.container.ConnectionString(s.ctx, "sslmode=disable")
	require.NoError(s.T(), err)

	cmd := exec.CommandContext(s.ctx, "migrate", "-database", dsn, "-path", "./migrations", "up")
	cmd.Dir = filepath.Join("..", "..")
	output, err := cmd.CombinedOutput()
	require.NoErrorf(s.T(), err, "migrate up failed: %s", output)
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
			require.NoError(t, repository.db.Ping())
			require.NoError(t, repository.Save(test.id, test.value))

			value, err := repository.Get(test.id)
			require.NoError(t, err)
			assert.Equal(t, test.value, value)
		})
	}
}

func TestPostgresRepository(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}

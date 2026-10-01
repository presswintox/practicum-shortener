package repository

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileRepository_SaveAndGet(t *testing.T) {
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
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "storage.json")
			r, err := NewFileRepository(path)
			t.Cleanup(func() { _ = r.Close() })
			require.NoError(t, err)
			require.NoError(t, r.Save(test.id, test.value))

			value, err := r.Get(test.id)
			require.NoError(t, err)
			assert.Equal(t, test.value, value)
		})
	}
}

func TestFileRepository_Load(t *testing.T) {
	path := filepath.Join(t.TempDir(), "storage.json")
	r, err := NewFileRepository(path)
	require.NoError(t, err)
	require.NoError(t, r.Save("4rSPg8ap", "http://yandex.ru"))
	require.NoError(t, r.Close())

	restored, err := NewFileRepository(path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = restored.Close() })
	require.NoError(t, restored.Load())
	value, err := restored.Get("4rSPg8ap")
	require.NoError(t, err)
	assert.Equal(t, "http://yandex.ru", value)
}

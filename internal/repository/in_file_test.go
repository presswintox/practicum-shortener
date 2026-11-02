package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileRepository_BatchSave(t *testing.T) {
	path := filepath.Join(t.TempDir(), "storage.json")
	r, err := NewFileRepository(path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = r.Close() })
	require.NoError(t, r.Save("existing", "https://example.com"))
	before, err := os.ReadFile(path)
	require.NoError(t, err)

	urls := map[string]string{"first": "https://example.com/1", "second": "https://example.com/2"}
	require.NoError(t, r.BatchSave(urls))
	after, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Greater(t, len(after), len(before))
	assert.Equal(t, before, after[:len(before)])
	for id, want := range urls {
		got, err := r.Get(id)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	}

	restored, err := NewFileRepository(path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = restored.Close() })
	urls["existing"] = "https://example.com"
	for id, want := range urls {
		got, err := restored.Get(id)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	}
	assert.Len(t, restored.records, len(urls))
}

func TestFileRepository_BatchSaveDuplicate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "storage.json")
	r, err := NewFileRepository(path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = r.Close() })
	require.NoError(t, r.Save("existing", "original"))
	before, err := os.ReadFile(path)
	require.NoError(t, err)

	require.ErrorIs(t, r.BatchSave(map[string]string{"new": "new value", "existing": "replacement"}), ErrAlreadyExists)
	after, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, before, after)
	assert.Len(t, r.records, 1)
	value, err := r.Get("existing")
	require.NoError(t, err)
	assert.Equal(t, "original", value)
}

func TestFileRepository_BatchSaveWriteError(t *testing.T) {
	r, err := NewFileRepository(filepath.Join(t.TempDir(), "storage.json"))
	require.NoError(t, err)
	require.NoError(t, r.Close())

	require.Error(t, r.BatchSave(map[string]string{"new": "value"}))
	_, err = r.Get("new")
	require.ErrorIs(t, err, ErrNotFound)
}

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

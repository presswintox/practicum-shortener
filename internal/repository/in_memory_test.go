package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryRepository_SaveAndGet(t *testing.T) {
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
			value: "abcd1234",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := NewMemoryRepository()
			require.NoError(t, r.Save(test.id, test.value))

			value, err := r.Get(test.id)
			require.NoError(t, err)
			assert.Equal(t, test.value, value)
		})
	}
}

func TestMemoryRepository_BatchSave(t *testing.T) {
	tests := []struct {
		name string
		urls map[string]string
		want map[string]string
		err  error
	}{
		{
			name: "multiple records",
			urls: map[string]string{"first": "value 1", "second": "value 2"},
			want: map[string]string{"existing": "original", "first": "value 1", "second": "value 2"},
		},
		{
			name: "duplicate rejects entire batch",
			urls: map[string]string{"first": "value 1", "existing": "replacement", "second": "value 2"},
			want: map[string]string{"existing": "original"},
			err:  ErrAlreadyExists,
		},
		{
			name: "empty batch",
			urls: map[string]string{},
			want: map[string]string{"existing": "original"},
		},
		{
			name: "nil batch",
			want: map[string]string{"existing": "original"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := NewMemoryRepository()
			require.NoError(t, r.Save("existing", "original"))

			err := r.BatchSave(test.urls)
			if test.err != nil {
				require.ErrorIs(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, test.want, r.db)
			for id, want := range test.want {
				got, err := r.Get(id)
				require.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}

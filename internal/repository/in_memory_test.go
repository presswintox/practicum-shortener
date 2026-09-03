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

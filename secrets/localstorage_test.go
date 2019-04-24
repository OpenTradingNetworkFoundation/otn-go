package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalStorage(t *testing.T) {
	cfg := LocalStorageConfig{
		Values: map[string]interface{}{
			"str":   "string value",
			"arr":   []string{"one", "two"},
			"other": 100,
		},
	}

	stg := NewLocalStorage(&cfg)
	assert.NotNil(t, stg)

	str, err := stg.ReadStringValue("str")
	assert.NoError(t, err)
	assert.Equal(t, "string value", str)

	arr, err := stg.ReadStringArray("arr")
	assert.NoError(t, err)
	assert.Equal(t, []string{"one", "two"}, arr)

	other, err := stg.ReadStringValue("other")
	assert.Error(t, err)
	assert.Equal(t, "", other)
}

package util_test

import (
	"testing"

	"github.com/8thgencore/passfort/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCacheKey(t *testing.T) {
	t.Run("Successfully generates cache key", func(t *testing.T) {
		prefix := "test"
		params := "123"
		expectedKey := "test:123"

		key := util.GenerateCacheKey(prefix, params)
		assert.Equal(t, expectedKey, key)
	})
}

func TestGenerateCacheKeyParams(t *testing.T) {
	t.Run("Successfully generates cache key params", func(t *testing.T) {
		params := []any{"123", "456", "789"}
		expectedKeyParams := "123-456-789"

		keyParams := util.GenerateCacheKeyParams(params...)
		assert.Equal(t, expectedKeyParams, keyParams)
	})
}

func TestSerialize(t *testing.T) {
	t.Run("Successfully serializes data", func(t *testing.T) {
		data := map[string]any{"key": "value"}
		expectedSerializedData := []byte(`{"key":"value"}`)

		serializedData, err := util.Serialize(data)
		assert.NoError(t, err)
		assert.Equal(t, expectedSerializedData, serializedData)
	})

	t.Run("Returns error when serialization fails", func(t *testing.T) {
		data := make(chan int)

		_, err := util.Serialize(data)
		assert.Error(t, err)
	})
}

func TestDeserialize(t *testing.T) {
	t.Run("Successfully deserializes data", func(t *testing.T) {
		serializedData := []byte(`{"key":"value"}`)
		var data map[string]any
		expectedData := map[string]any{"key": "value"}

		err := util.Deserialize(serializedData, &data)
		assert.NoError(t, err)
		assert.Equal(t, expectedData, data)
	})

	t.Run("Returns error when deserialization fails", func(t *testing.T) {
		serializedData := []byte(`{"key":}`)
		var data map[string]any

		err := util.Deserialize(serializedData, &data)
		assert.Error(t, err)
	})
}

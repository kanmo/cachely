package cachely

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	c := New()
	assert.NotNil(t, c, "New() should return a Cache object")
}

type TestCacheRecord struct {
	Message string
}

func TestGet(t *testing.T) {
	t.Run("cache does not exist", func(t *testing.T) {
		c := New()
		var result TestCacheRecord
		key := "key"

		// Save data to cache because there is no data in the cache
		err := c.Get(key, &result, 0, func(dest interface{}) error {
			result.Message = "testrecord"
			return nil
		})
		assert.Nil(t, err, "Get() should not return an error")

		var result2 TestCacheRecord
		// Get saved data from the cache
		err = c.Get(key, &result2, 0, func(dest interface{}) error {
			t.Fatal("fn should not be called")
			return nil
		})
		assert.Nil(t, err, "Get() should not return an error")
		assert.Equal(t, "testrecord", result2.Message, "Get() should return the saved data")

	})

	t.Run("it expires cached data when cache has expiration", func(t *testing.T) {
		c := New()
		var result TestCacheRecord
		key := "key"

		// Save data to cache with expiration
		err := c.Get(key, &result, 20*time.Millisecond, func(dest interface{}) error {
			result.Message = "testrecord"
			return nil
		})
		assert.Nil(t, err, "Get() should not return an error")

		var result2 TestCacheRecord
		// Get saved data from the cache
		err = c.Get(key, &result2, 0, func(dest interface{}) error {
			return nil
		})
		assert.Nil(t, err, "Get() should not return an error")
		assert.Equal(t, "testrecord", result2.Message, "Get() should return the saved data")

		var result3 TestCacheRecord
		// Data that has expired cannot be retrieved
		<-time.After(25 * time.Millisecond)
		err = c.Get(key, &result3, 0, func(dest interface{}) error {
			return nil
		})
		assert.Nil(t, err, "Get() should not return an error")
		assert.Equal(t, "", result3.Message, "Get() should return an empty string")
	})
}

func TestSet(t *testing.T) {
	t.Run("set data to cache", func(t *testing.T) {
		c := New()
		var result TestCacheRecord
		result.Message = "testrecord"
		key := "key"

		c.Set(key, &result, 0)
		assert.Equal(t, 1, len(c.cache), "Set() should add data to the cache")
		assert.Equal(t, &result, c.cache[key].Records, "Set() should add the data to the cache")
	})

	t.Run("set data to cache with expiration", func(t *testing.T) {
		c := New()
		var result TestCacheRecord
		result.Message = "testrecord"
		key := "key"

		c.Set(key, &result, 20*time.Millisecond)
		assert.Equal(t, 1, len(c.cache), "Set() should add data to the cache")
		assert.Equal(t, &result, c.cache[key].Records, "Set() should add the data to the cache")

		<-time.After(25 * time.Millisecond)
		_, ok := c.get(key)
		assert.False(t, ok, "Set() should remove the data from the cache after expiration")
	})
}

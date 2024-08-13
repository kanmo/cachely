package cachely

import (
	"reflect"
	"sync"
	"time"
)

// CacheObject is a struct that holds the records and the expiration time.
// Records is the data that is stored in the cache.
type CacheObject struct {
	Records    interface{}
	Expiration int64
}

// Cache is a simple in-memory cache with expiration.
// cache field is a map of string keys to CacheObject values.
type Cache struct {
	mu    sync.RWMutex
	cache map[string]CacheObject
}

// New creates a new Cache object.
// It initializes the cache map.
func New() *Cache {
	c := &Cache{
		cache: make(map[string]CacheObject),
	}
	return c
}

// Get retrieves the data from the cache.
// If the data is not found in the cache, it calls the provided function to cache the data.
// If the data is found in the cache, it copies the data to the provided destination.
func (c *Cache) Get(key string, dest interface{}, d time.Duration, fn func(dest interface{}) error) error {
	records, found := c.get(key)

	if !found {
		err := fn(dest)
		if err != nil {
			return err
		}

		rv := reflect.ValueOf(dest).Elem()

		if (rv.Kind() == reflect.Slice && rv.Len() > 0) || rv.Kind() != reflect.Slice {
			c.Set(key, rv.Interface(), d)
		}
		return nil
	}

	rv := reflect.ValueOf(dest).Elem()
	rv2 := reflect.ValueOf(records)
	rv.Set(rv2)

	return nil
}

// Flush clears the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	c.cache = map[string]CacheObject{}
	c.mu.Unlock()
}

func (c *Cache) Set(k string, x interface{}, d time.Duration) {
	var e int64
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}

	c.mu.Lock()
	c.cache[k] = CacheObject{
		x,
		e,
	}
	c.mu.Unlock()
}

func (c *Cache) get(k string) (interface{}, bool) {
	c.mu.RLock()
	obj, ok := c.cache[k]
	if !ok {
		c.mu.RUnlock()
		return nil, false
	}

	if obj.Expiration > 0 {
		if time.Now().UnixNano() > obj.Expiration {
			c.mu.RUnlock()
			return nil, false
		}
	}

	c.mu.RUnlock()
	return obj.Records, true
}

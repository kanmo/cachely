# cachely

cachely is a lightweight, in-memory caching library for Go applications. It provides a simple and thread-safe way to cache data with expiration, making it ideal for improving the performance of your applications by reducing the need for repeated expensive computations or database queries.

## Features

- In-Memory Cache: Store data in memory for quick access.
- Expiration: Set expiration times for cached data to ensure that stale data is not served.
- Thread-Safe: Safe to use across multiple goroutines with efficient locking mechanisms.
- Flexible Data Retrieval: Automatically populate the cache if the data is not found, using a provided function.

## Usage

### Basic Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/kanmo/cachely"
)

type CacheRecord struct {
	Message string
}

func main() {
	cache := cachely.New()

	key := "user:123"
	var cacheData CacheRecord

	// Try to get the data from the cache
	err := cache.Get(key, &cacheData, 5*time.Minute, func(dest interface{}) error {
		cacheData.Message = "Data loaded from cache"
		return nil
	})
	if err != nil {
		fmt.Println("Error retrieving data:", err)
		return
	}

	fmt.Println("Cache Data:", cacheData.Message)

	// Flush the cache
	cache.Flush()
}
```

## API

### New() *Cache

Creates a new Cache instance.

### Get(key string, dest interface{}, d time.Duration, fn func(dest interface{}) error) error

Retrieves the data associated with the given key from the cache. If the data is not found, the provided function fn is called to fetch and cache the data. The cached data is stored for the duration specified by d.

- key: The key under which the data is cached.
- dest: A pointer to the variable where the data will be stored.
- d: The expiration duration for the cached data.
- fn: A function that fetches the data if it’s not found in the cache.

Returns an error if the data could not be fetched or cached.

### Set(k string, x interface{}, d time.Duration)

Manually sets a value in the cache with the specified expiration duration.

- k: The key under which the data is stored.
- x: The data to be cached.
- d: The expiration duration for the cached data.

### Flush()

Clears all data from the cache.

### License

This library is licensed under the MIT License. See the [LICENSE.txt](LICENSE.txt) file for details.

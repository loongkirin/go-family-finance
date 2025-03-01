package cache

import (
	"errors"
	"time"
)

const (
	DEFAULT = time.Duration(0)
	FOREVER = time.Duration(-1)
)

var (
	ErrNotStored  = errors.New("cache: item not stored")
	ErrNotSupport = errors.New("cache: not support")
	ErrKeyExists  = errors.New("cache: item already exists")
	ErrCacheMiss  = errors.New("cache: item not found")
)

// CacheStore is the interface of a cache backend
type CacheStore interface {
	// Get retrieves an item from the cache. Returns the item or nil
	Get(key string) (string, error)

	// Set sets an item to the cache, replacing any existing item.
	Set(key string, value interface{}, expire time.Duration) error

	// Add adds an item to the cache only if an item doesn't already exist for the given
	// key, or if the existing item has expired. Returns an error otherwise.
	Add(key string, value interface{}, expire time.Duration) error

	// Replace sets a new value for the cache key only if it already exists. Returns an
	// error if it does not.
	Replace(key string, value interface{}, expire time.Duration) error

	// Delete removes an item from the cache. Does nothing if the key is not in the cache.
	Delete(key string) error

	// Increment increments a real number, and returns error if the value is not real
	Increment(key string, value int64) (int64, error)

	// Decrement decrements a real number, and returns error if the value is not real
	Decrement(key string, value int64) (int64, error)

	// Flush deletes all items from the cache.
	Flush() error
}

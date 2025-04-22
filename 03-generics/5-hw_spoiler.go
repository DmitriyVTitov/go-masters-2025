package generics

import (
	"sync"
	"time"
)

// Cache - структура для обобщённого кэша
type Cache[K comparable, V any] struct {
	data     map[K]cacheItem[V]
	mutex    sync.RWMutex
	lifetime time.Duration
}

type cacheItem[V any] struct {
	value      V
	expiration time.Time
}

// NewCache - создание нового кэша с заданным временем жизни элементов
func NewCache[K comparable, V any](lifetime time.Duration) *Cache[K, V] {
	return &Cache[K, V]{
		data:     make(map[K]cacheItem[V]),
		lifetime: lifetime,
	}
}

// Set - добавление значения в кэш
func (c *Cache[K, V]) Set(key K, value V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = cacheItem[V]{
		value:      value,
		expiration: time.Now().Add(c.lifetime),
	}
}

// Get - получение значения из кэша
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists || time.Now().After(item.expiration) {
		var zeroValue V
		return zeroValue, false
	}
	return item.value, true
}

// Delete - удаление значения из кэша
func (c *Cache[K, V]) Delete(key K) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.data, key)
}

// Cleanup - удаление устаревших элементов из кэша
func (c *Cache[K, V]) Cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for key, item := range c.data {
		if now.After(item.expiration) {
			delete(c.data, key)
		}
	}
}

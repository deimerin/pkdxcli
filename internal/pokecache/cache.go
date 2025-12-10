package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cache    map[string]cacheEntry
	interval time.Duration
	mu       sync.Mutex
}

func NewCache(interval time.Duration) *Cache {

	c := &Cache{
		cache:    map[string]cacheEntry{},
		interval: interval,
		mu:       sync.Mutex{},
	}

	go c.reapLoop()

	return c

}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.cache[key]
	if ok {
		return entry.val, true
	}
	return nil, false
}

func (c *Cache) reapLoop() {

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()

		for key, entry := range c.cache {
			age := time.Since(entry.createdAt)
			if age > c.interval {
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}

}

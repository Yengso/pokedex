package pokecache

import (
	"time"
	"sync"
)

type Cache struct {
	entries		map[string]cacheEntry
	mu			sync.RWMutex
	interval 	time.Duration
}

type cacheEntry struct {
	val 		[]byte
	createdAt	time.Time
}

func NewCache(t time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]cacheEntry),
		interval: t,
	}
 
	go c.reapLoop()

	return c
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		val:		value,
		createdAt:  time.Now(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, ok
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.entries {
			age := time.Since(entry.createdAt)
			if age > c.interval {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}
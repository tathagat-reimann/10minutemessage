package cache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheMap sync.Map
}

type CacheItem struct {
	Value      string
	Expiration time.Time
}

func (c *Cache) Set(key string, value string, expiration time.Duration) {
	c.cacheMap.Store(key, CacheItem{value, time.Now().Add(expiration)})
}

func (c *Cache) Get(key string) (string, bool) {
	if val, ok := c.cacheMap.Load(key); ok {
		item := val.(CacheItem)
		if time.Now().Before(item.Expiration) {
			return item.Value, true
		} else {
			c.cacheMap.Delete(key)
		}
	}
	return "", false
}

func (c *Cache) Delete(key string) {
	c.cacheMap.Delete(key)
}

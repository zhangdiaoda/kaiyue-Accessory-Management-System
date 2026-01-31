package notification

import (
	"sync"
	"time"
)

// SimpleCache 简单的内存缓存实现(生产环境建议使用Redis)
type SimpleCache struct {
	data map[string]cacheItem
	mu   sync.RWMutex
}

type cacheItem struct {
	value      string
	expiration time.Time
}

// NewSimpleCache 创建简单缓存
func NewSimpleCache() *SimpleCache {
	cache := &SimpleCache{
		data: make(map[string]cacheItem),
	}

	// 启动过期清理goroutine
	go cache.cleanupExpired()

	return cache
}

// Get 获取缓存
func (c *SimpleCache) Get(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.data[key]
	if !ok {
		return "", nil
	}

	if time.Now().After(item.expiration) {
		return "", nil
	}

	return item.value, nil
}

// Set 设置缓存
func (c *SimpleCache) Set(key string, value string, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(expiration),
	}

	return nil
}

// Delete 删除缓存
func (c *SimpleCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
	return nil
}

// cleanupExpired 清理过期缓存
func (c *SimpleCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.data {
			if now.After(item.expiration) {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}

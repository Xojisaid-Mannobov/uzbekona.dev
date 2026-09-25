// Package cache public API javoblarini qisqa muddat saqlash uchun.
// REDIS_URL berilsa Redis, aks holda jarayon ichidagi xotira ishlatiladi.
package cache

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, bool)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration)
	DeletePrefix(ctx context.Context, prefix string)
}

// ─── Redis ──────────────────────────────────────────────────

type redisCache struct {
	client *redis.Client
}

func NewRedis(client *redis.Client) Cache {
	return &redisCache{client: client}
}

func (c *redisCache) Get(ctx context.Context, key string) ([]byte, bool) {
	v, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}
	return v, true
}

func (c *redisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) {
	_ = c.client.Set(ctx, key, value, ttl).Err()
}

func (c *redisCache) DeletePrefix(ctx context.Context, prefix string) {
	iter := c.client.Scan(ctx, 0, prefix+"*", 200).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if len(keys) > 0 {
		_ = c.client.Del(ctx, keys...).Err()
	}
}

// ─── Xotira ─────────────────────────────────────────────────

type entry struct {
	value   []byte
	expires time.Time
}

type memoryCache struct {
	mu    sync.RWMutex
	items map[string]entry
}

func NewMemory() Cache {
	c := &memoryCache{items: make(map[string]entry)}
	go c.janitor()
	return c
}

func (c *memoryCache) Get(_ context.Context, key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.items[key]
	if !ok || time.Now().After(e.expires) {
		return nil, false
	}
	return e.value, true
}

func (c *memoryCache) Set(_ context.Context, key string, value []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = entry{value: value, expires: time.Now().Add(ttl)}
}

func (c *memoryCache) DeletePrefix(_ context.Context, prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k := range c.items {
		if strings.HasPrefix(k, prefix) {
			delete(c.items, k)
		}
	}
}

// janitor — muddati o'tgan yozuvlarni vaqti-vaqti bilan tozalaydi.
func (c *memoryCache) janitor() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		now := time.Now()
		c.mu.Lock()
		for k, e := range c.items {
			if now.After(e.expires) {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}

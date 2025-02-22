package storage

import (
	"context"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Cache interface {
	Load(ctx context.Context, key string) (int, bool)
	Set(ctx context.Context, key string, value int, expiration time.Duration) error
}

type InMemoryCache struct {
	data sync.Map
	log  *zap.Logger
}

type cacheEntry struct {
	value      int
	expiration time.Time
}

func NewInMemoryCache(log *zap.Logger) *InMemoryCache {
	return &InMemoryCache{
		data: sync.Map{},
		log:  log,
	}
}

func (c *InMemoryCache) Load(ctx context.Context, key string) (int, bool) {
	select {
	case <-ctx.Done():
		return 0, false
	default:
		if entry, ok := c.data.Load(key); ok {
			cacheEntry, ok := entry.(cacheEntry)
			if !ok {
				return 0, false
			}

			if !cacheEntry.expiration.IsZero() && time.Now().After(cacheEntry.expiration) {
				c.data.Delete(key)
				return 0, false
			}

			return cacheEntry.value, true
		}
		return 0, false
	}
}

func (c *InMemoryCache) Set(ctx context.Context, key string, value int, expiration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		var expTime time.Time
		if expiration > 0 {
			expTime = time.Now().Add(expiration)
		}

		c.data.Store(key, cacheEntry{value: value, expiration: expTime})

		if expiration > 0 {
			keyCopy := key
			time.AfterFunc(expiration, func() {
				if _, ok := c.data.Load(keyCopy); ok {
					c.data.Delete(keyCopy)
					c.log.Info("Cache entry expired", zap.String("key", keyCopy))
				}
			})
		}
		c.log.Info("Cache set", zap.String("key", key), zap.Int("value", value))
		return nil
	}
}

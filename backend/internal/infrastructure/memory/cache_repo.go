package memory

import (
	"context"
	"time"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/repository"
)

type memoryCacheRepo struct {
	cache *MemoryCache
}

// NewCacheRepository creates a CacheRepository backed by MemoryCache.
func NewCacheRepository(ttl time.Duration, cleanupInterval time.Duration) (repository.CacheRepository, func()) {
	c := NewMemoryCache(ttl, cleanupInterval)
	repo := &memoryCacheRepo{cache: c}
	return repo, c.Close
}

func (r *memoryCacheRepo) Get(ctx context.Context, key string) (interface{}, bool) {
	return r.cache.Get(key)
}

func (r *memoryCacheRepo) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	r.cache.Set(key, value, ttl)
}

func (r *memoryCacheRepo) Delete(ctx context.Context, key string) {
	r.cache.Delete(key)
}

func (r *memoryCacheRepo) Stats() map[string]interface{} {
	return r.cache.Stats()
}

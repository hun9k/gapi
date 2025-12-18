package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

func NewLocal(DefaultExpiration, CleanupInterval time.Duration) Cacher {
	// parse dsn
	// create cache
	return &LocalCacher{
		cache.New(DefaultExpiration, CleanupInterval),
	}
}

type LocalCacher struct {
	*cache.Cache
}

func (c *LocalCacher) Add(key string, value any, expiration time.Duration) error {
	// parse options
	return c.Cache.Add(key, value, expiration)
}

func (c *LocalCacher) Set(key string, value any, expiration time.Duration) {
	// parse options
	c.Cache.Set(key, value, expiration)
}

func (c *LocalCacher) Get(key string) (any, bool) {
	// TODO: Implement cache get
	return c.Cache.Get(key)
}

func (c *LocalCacher) Delete(key string) {
	// TODO: Implement cache delete
	c.Cache.Delete(key)
}

func (c *LocalCacher) Clear() {
	// TODO: Implement cache clear
	c.Cache.Flush()
}

func (c *LocalCacher) Stats() map[string]any {
	return nil
}

package cacheman

import (
	"time"

	"github.com/allegro/bigcache"
)

// NewBigCache creates big cache
func NewBigCache(config *Config) (*bigcache.BigCache, error) {
	ttl, e := time.ParseDuration(config.TTL)
	if e != nil {
		ttl, _ = time.ParseDuration(defaultTTL)
	}
	return bigcache.NewBigCache(bigcache.DefaultConfig(ttl))
}

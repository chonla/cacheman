package cacheman

import (
	"fmt"
	"time"

	"github.com/allegro/bigcache"
)

type BigCacheClient struct {
	client *bigcache.BigCache
}

// NewBigCache creates big cache client
func NewBigCache(config *Config) (*BigCacheClient, error) {
	ttl, e := time.ParseDuration(config.TTL)
	if e != nil {
		ttl, _ = time.ParseDuration(defaultTTL)
	}
	client, e := bigcache.NewBigCache(bigcache.DefaultConfig(ttl))
	if e != nil {
		return nil, e
	}
	return &BigCacheClient{
		client: client,
	}, nil
}

func (c *BigCacheClient) Get(key string) ([]byte, error) {
	return c.client.Get(key)
}

func (c *BigCacheClient) Set(key string, value []byte) error {
	return c.client.Set(key, value)
}

func (c *BigCacheClient) Delete(key string) error {
	return c.client.Delete(key)
}

func (c *BigCacheClient) Reset() error {
	return c.client.Reset()
}

func (c *BigCacheClient) Type() string {
	return fmt.Sprintf("%T", c)
}

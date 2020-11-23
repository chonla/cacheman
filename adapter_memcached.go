package cacheman

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedClient struct {
	client *memcache.Client
	ttl    time.Duration
}

// NewMemcached creates big cache client
func NewMemcached(config *Config) (*MemcachedClient, error) {
	ttl, e := time.ParseDuration(config.TTL)
	if e != nil {
		ttl, _ = time.ParseDuration(defaultTTL)
	}
	client := memcache.New(config.Server)
	return &MemcachedClient{
		client: client,
		ttl:    ttl,
	}, nil
}

func (c *MemcachedClient) Get(key string) ([]byte, error) {
	result, e := c.client.Get(key)
	if e != nil {
		return nil, e
	}
	return result.Value, nil
}

func (c *MemcachedClient) Set(key string, value []byte) error {
	return c.client.Set(&memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: int32(c.ttl.Seconds()),
	})
}

func (c *MemcachedClient) Delete(key string) error {
	return c.client.Delete(key)
}

func (c *MemcachedClient) Reset() error {
	return c.client.DeleteAll()
}

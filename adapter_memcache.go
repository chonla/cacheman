package cacheman

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheClient struct {
	client *memcache.Client
	ttl    time.Duration
}

// NewMemcache creates big cache client
func NewMemcache(config *Config) (*MemcacheClient, error) {
	ttl, e := time.ParseDuration(config.TTL)
	if e != nil {
		ttl, _ = time.ParseDuration(defaultTTL)
	}
	client := memcache.New(config.Server)
	return &MemcacheClient{
		client: client,
		ttl:    ttl,
	}, nil
}

func (c *MemcacheClient) Get(key string) ([]byte, error) {
	result, e := c.client.Get(key)
	if e != nil {
		return nil, e
	}
	return result.Value, nil
}

func (c *MemcacheClient) Set(key string, value []byte) error {
	return c.client.Set(&memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: int32(c.ttl.Seconds()),
	})
}

func (c *MemcacheClient) Delete(key string) error {
	return c.client.Delete(key)
}

func (c *MemcacheClient) Reset() error {
	return c.client.DeleteAll()
}

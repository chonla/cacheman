package cacheman

// CacheInterface defines interface for cache
type CacheInterface interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}

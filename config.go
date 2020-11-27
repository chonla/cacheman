package cacheman

// Config for cacheman
type Config struct {
	// Enabled to enable/disable cacheman
	Enabled bool
	// Verbose allow activities of cacheman to be display on console
	Verbose bool
	// TTL is age of cache entry in duration format, e.g. 1d for one day
	TTL string
	// Paths that will be cached
	Paths []string
	// ExcludedPaths are paths to be excluded from cache
	ExcludedPaths []string
	// AdditionalHeaders are injected in return cache
	AdditionalHeaders map[string]string
	// Server is cache server in host:port format
	Server string
	// Password is credential for accessing cache service
	Password string
	// Database is database name or index
	Database interface{}
	// CacheInfoPath is URI to request cache information
	CacheInfoPath string
	// PurgePath is URI to purge all content in cache
	PurgePath string
	// Namespace to be automatically added into cache key
	Namespace string
}

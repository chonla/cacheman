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
	// AdditionalHeaders are injected in return cache
	AdditionalHeaders map[string]string
}

# CacheMan

CacheMan was designed to be middleware for Echo for caching response from `GET` request for a period of time.

## Usage (echo 3)

```go
	server := echo.New()

	store, e := cacheman.NewBigCache(&cfg.Cache)
	if e == nil {
		server.Use(cacheman.Middleware(&cfg.Cache, store))
	}
```

## Usage (echo 4)

```go
	server := echo.New()

	store, e := cacheman.NewBigCache(&cfg.Cache)
	if e == nil {
		server.Use(cacheman.MiddlewareV4(&cfg.Cache, store))
	}
```

## Working example

[cacheman-example](https://github.com/chonla/cacheman-example)

## Cache support

* BigCache - [allegro/bigcache](github.com/allegro/bigcache)
* Memcache - [bradfitz/gomemcache](github.com/bradfitz/gomemcache/memcache)
* Redis - [go-redis/redis](github.com/go-redis/redis/v8)

## Custom cache

Just implement this interface and pass it into `cacheman.Middleware`.

```go
type CacheInterface interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
	Reset() error
	Type() string
}
```

## Configuration

### Enabled

Set to true to enable cacheman or false to disable. Default is `false`,

### Verbose

Set to true to print out cacheman log or falst to make cacheman quiet. Default is `false`,

### TTL

Cache entry life span in duration format. For example, `5m` for 5 minutes, `1d` for 1 day. Default is `1d`.

### Paths

Paths to be cached. Path with embedded variables like `/user/:id` is supported. Regular expression string is also supported, `/.*` to cache every path. Default is `[]string{}`,

### ExcludedPaths

Paths to be excluded from cache. ExcludedPaths has higher priority than Paths. Default is `[]string{}`,

### AdditionalHeaders

Custom headers added into returned cache. Default is `map[string]string{}`,

### CacheInfoPath
URI to request cache information. Make it empty to disable it. Default is `<empty>`.

## License

[MIT](LICENSE)
# CacheMan

Cache Middleware for Echo

## Usage

```go
	server := echo.New()

	store, e := cacheman.NewBigCache(&cfg.Cache)
	if e == nil {
		server.Use(cacheman.Middleware(&cfg.Cache, store))
	}
```

## Cache support

* BigCache

## Custom cache

Just implement this interface and pass it into `cacheman.Middleware`.

```go
type CacheInterface interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
	Reset() error
}
```

## License

[MIT](LICENSE)
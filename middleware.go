package cacheman

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/labstack/echo"
	echo4 "github.com/labstack/echo/v4"
)

// Middleware creates a middleware to handle cache
func Middleware(config *Config, cache CacheInterface) echo.MiddlewareFunc {
	cm = NewCacheManager(config, cache)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			cm.Log(fmt.Sprintf("Test path: %s", ctx.Request().RequestURI))
			if cm.Enabled &&
				ctx.Request().Method == "GET" &&
				cm.TestPath(ctx.Request().URL.Path) {

				cm.Log(fmt.Sprintf("Path matches: %s", ctx.Request().RequestURI))

				interceptor := NewInterceptor(ctx.Response().Writer)
				ctx.Response().Writer = interceptor

				if !cm.TryWrite(ctx) {
					e := next(ctx)
					// Store into cache only if status is 200
					if e == nil && interceptor.Status() == 200 {
						content := Content{
							Status:  interceptor.Status(),
							Headers: interceptor.Header(),
							Content: base64.StdEncoding.EncodeToString(interceptor.Content()),
						}
						stringifiedCache, e := json.Marshal(content)
						if e == nil {
							cacheKey := ctx.Request().RequestURI
							cm.Set(cacheKey, stringifiedCache)
						}
					}
					return e
				}
				return nil
			}
			return next(ctx)
		}
	}
}

// MiddlewareV4 creates a middleware to handle cache for echo V4
func MiddlewareV4(config *Config, cache CacheInterface) echo4.MiddlewareFunc {
	cm = NewCacheManager(config, cache)
	return func(next echo4.HandlerFunc) echo4.HandlerFunc {
		return func(ctx echo4.Context) error {
			if cm.Enabled {
				cm.Log(fmt.Sprintf("Test path: %s", ctx.Request().RequestURI))
				if ctx.Request().Method == "GET" &&
					cm.TestPath(ctx.Request().URL.Path) {

					cm.Log(fmt.Sprintf("Path matches: %s", ctx.Request().RequestURI))

					interceptor := NewInterceptor(ctx.Response().Writer)
					ctx.Response().Writer = interceptor

					if !cm.TryWriteV4(ctx) {
						e := next(ctx)
						// Store into cache only if status is 200
						if e == nil && interceptor.Status() == 200 {
							content := Content{
								Status:  interceptor.Status(),
								Headers: interceptor.Header(),
								Content: base64.StdEncoding.EncodeToString(interceptor.Content()),
							}
							stringifiedCache, e := json.Marshal(content)
							if e == nil {
								cacheKey := ctx.Request().RequestURI
								cm.Set(cacheKey, stringifiedCache)
							}
						}
						return e
					}
					return nil
				} else {
					cm.Log(fmt.Sprintf("Path does not match: %s", ctx.Request().RequestURI))
				}
			}
			return next(ctx)
		}
	}
}

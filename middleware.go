package cacheman

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/labstack/echo"
)

var cm *Manager
var defaultTTL = "5m"

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

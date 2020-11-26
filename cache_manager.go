package cacheman

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/labstack/echo"
	echo4 "github.com/labstack/echo/v4"
)

// Manager is cache manager
type Manager struct {
	Enabled                  bool
	Verbose                  bool
	Cache                    CacheInterface
	Routes                   []string
	ExcludedRoutes           []string
	RouteCount               int
	ExcludedRouteCount       int
	ComparableRoutes         []*regexp.Regexp
	ComparableExcludedRoutes []*regexp.Regexp
	AdditionalHeaders        map[string]string
}

// Content is cached content
type Content struct {
	Status  int         `json:"status"`
	Headers http.Header `json:"headers"`
	Content string      `json:"content"`
}

const (
	// HealthCheckUntested tells that operation has not been tested
	HealthCheckUntested string = "untested"
	// HealthCheckFailed tells that operation has failed the test
	HealthCheckFailed string = "failed"
	// HealthCheckPassed tells that operation has passed the test
	HealthCheckPassed string = "passed"
)

type healthCheckResult struct {
	SetResult    string `json:"setResult"`
	GetResult    string `json:"getResult"`
	DeleteResult string `json:"deleteResult"`
}

var cm *Manager
var defaultTTL = "5m"

// NewCacheManager creates a cache manager
func NewCacheManager(conf *Config, cache CacheInterface) *Manager {
	comparableRoutes := convertToComparableRoutes(conf.Paths)
	comparableExcludedRoutes := convertToComparableRoutes(conf.ExcludedPaths)

	return &Manager{
		Enabled:                  conf.Enabled,
		Verbose:                  conf.Verbose,
		Cache:                    cache,
		Routes:                   conf.Paths,
		ComparableRoutes:         comparableRoutes,
		ExcludedRoutes:           conf.ExcludedPaths,
		ComparableExcludedRoutes: comparableExcludedRoutes,
		RouteCount:               len(conf.Paths),
		ExcludedRouteCount:       len(conf.ExcludedPaths),
		AdditionalHeaders:        conf.AdditionalHeaders,
	}
}

// convertToComparableRoutes converts routes array of string to array of regular expression
func convertToComparableRoutes(routes []string) []*regexp.Regexp {
	// Good route
	// /some/path
	// /some/other/path/with/:variable-inside
	output := []*regexp.Regexp{}
	for routeIndex, routeCount := 0, len(routes); routeIndex < routeCount; routeIndex++ {
		path := routes[routeIndex]
		if path == "" {
			path = "/"
		}
		if path[0] != '/' {
			path = "/" + path
		}
		fragments := strings.Split(path, "/")

		for fragmentIndex, fragmentCount := 0, len(fragments); fragmentIndex < fragmentCount; fragmentIndex++ {
			if len(fragments[fragmentIndex]) > 0 && fragments[fragmentIndex][0] == ':' {
				fragments[fragmentIndex] = ".+"
			}
		}
		regString := fmt.Sprintf("^%s$", strings.Join(fragments, "/"))
		output = append(output, regexp.MustCompile(regString))
	}
	return output
}

// TestPath return true if path matches a route, otherwise returns false
func (c *Manager) TestPath(path string) bool {
	for routeIndex := 0; routeIndex < c.ExcludedRouteCount; routeIndex++ {
		if c.ComparableExcludedRoutes[routeIndex].MatchString(path) {
			return false
		}
	}
	for routeIndex := 0; routeIndex < c.RouteCount; routeIndex++ {
		if c.ComparableRoutes[routeIndex].MatchString(path) {
			return true
		}
	}
	return false
}

// Get gets byte content from path key
func (c *Manager) Get(path string) ([]byte, bool) {
	content, e := c.Cache.Get(path)
	if e != nil {
		c.Log(fmt.Sprintf("Cache misses: %s", path))
		return []byte{}, false
	}
	c.Log(fmt.Sprintf("Cache hits: %s", path))
	return content, true
}

// Set sets byte content to path key
func (c *Manager) Set(path string, b []byte) error {
	c.Log(fmt.Sprintf("Cache sets: %s", path))
	return c.Cache.Set(path, b)
}

// TryWrite tries to write cached content if hit and return true, return false if miss
func (c *Manager) TryWrite(ctx echo.Context) bool {
	cacheKey := ctx.Request().RequestURI
	stringifiedCache, e := c.Get(cacheKey)
	if !e {
		return false
	}

	var content Content
	err := json.Unmarshal(stringifiedCache, &content)
	if err != nil {
		return false
	}

	writer := ctx.Response().Writer
	for headerKey, headerValues := range content.Headers {
		for _, headerValue := range headerValues {
			writer.Header().Set(headerKey, headerValue)
		}
	}
	for headerKey, headerValue := range c.AdditionalHeaders {
		writer.Header().Set(headerKey, headerValue)
	}

	writer.WriteHeader(content.Status)
	byteContent, err := base64.StdEncoding.DecodeString(content.Content)
	if err != nil {
		return false
	}
	writer.Write(byteContent)
	return true
}

// TryWriteV4 tries to write cached content if hit and return true, return false if miss
func (c *Manager) TryWriteV4(ctx echo4.Context) bool {
	cacheKey := ctx.Request().RequestURI
	stringifiedCache, e := c.Get(cacheKey)
	if !e {
		return false
	}

	var content Content
	err := json.Unmarshal(stringifiedCache, &content)
	if err != nil {
		return false
	}

	writer := ctx.Response().Writer
	for headerKey, headerValues := range content.Headers {
		for _, headerValue := range headerValues {
			writer.Header().Set(headerKey, headerValue)
		}
	}
	for headerKey, headerValue := range c.AdditionalHeaders {
		writer.Header().Set(headerKey, headerValue)
	}

	writer.WriteHeader(content.Status)
	byteContent, err := base64.StdEncoding.DecodeString(content.Content)
	if err != nil {
		return false
	}
	writer.Write(byteContent)
	return true
}

// Log prints log message
func (c *Manager) Log(msg string) {
	if c.Verbose {
		fmt.Println(msg)
	}
}

func (c *Manager) healthCheck() *healthCheckResult {
	cacheContent := "CacheMan"
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	cacheKey := fmt.Sprintf("cacheman-%d", r1.Intn(1000))
	result := &healthCheckResult{
		SetResult:    HealthCheckUntested,
		GetResult:    HealthCheckUntested,
		DeleteResult: HealthCheckUntested,
	}
	err := c.Set(cacheKey, []byte(cacheContent))
	if err != nil {
		result.SetResult = HealthCheckFailed
	} else {
		result.SetResult = HealthCheckPassed
	}
	if result.SetResult == HealthCheckPassed {
		content, found := c.Get(cacheKey)
		if !found || string(content) != cacheContent {
			result.GetResult = HealthCheckFailed
		} else {
			result.GetResult = HealthCheckPassed
		}
		if result.GetResult == HealthCheckPassed {
			err = c.Cache.Delete(cacheKey)
			if err != nil {
				result.DeleteResult = HealthCheckFailed
			} else {
				result.DeleteResult = HealthCheckPassed
			}
		}
	}
	return result
}

// WriteInfo print cacheman information out to client
func (c *Manager) WriteInfo(ctx echo.Context) {
	healthResult := c.healthCheck()

	response := map[string]interface{}{
		"type":            c.Cache.Type(),
		"operationHealth": healthResult,
	}

	ctx.JSON(200, response)
}

// WriteInfoV4 print cacheman information out to client
func (c *Manager) WriteInfoV4(ctx echo4.Context) {
	healthResult := c.healthCheck()

	response := map[string]interface{}{
		"type":            c.Cache.Type(),
		"operationHealth": healthResult,
	}
	ctx.JSON(200, response)
}

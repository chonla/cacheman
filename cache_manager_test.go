package cacheman

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func (o *MockCache) Get(k string) ([]byte, error) {
	args := o.Called(k)
	return args.Get(0).([]byte), args.Error(1)
}

func (o *MockCache) Set(k string, b []byte) error {
	args := o.Called(k, b)
	return args.Error(0)
}

func TestMatchPathWithExactMatch(t *testing.T) {
	conf := &Config{
		Enabled: true,
		Verbose: false,
		TTL:     "1m",
		Paths: []string{
			"/test1",
			"/test2",
		},
		AdditionalHeaders: map[string]string{},
	}
	cm := NewCacheManager(conf, nil)

	result := cm.TestPath("/test2")

	assert.True(t, result)
}

func TestMatchPathWithMissingPrecedingSlash(t *testing.T) {
	conf := &Config{
		Enabled: true,
		Verbose: false,
		TTL:     "1m",
		Paths: []string{
			"test1",
			"test2",
		},
		AdditionalHeaders: map[string]string{},
	}
	cm := NewCacheManager(conf, nil)

	result := cm.TestPath("/test2")

	assert.True(t, result)
}

func TestMatchPathWithEmptyPath(t *testing.T) {
	conf := &Config{
		Enabled: true,
		Verbose: false,
		TTL:     "1m",
		Paths: []string{
			"test1",
			"test2",
			"",
		},
		AdditionalHeaders: map[string]string{},
	}
	cm := NewCacheManager(conf, nil)

	result := cm.TestPath("/")

	assert.True(t, result)
}

func TestMatchPathWithUnmatch(t *testing.T) {
	conf := &Config{
		Enabled: true,
		Verbose: false,
		TTL:     "1m",
		Paths: []string{
			"/test1",
			"/test2",
		},
		AdditionalHeaders: map[string]string{},
	}
	cm := NewCacheManager(conf, nil)

	result := cm.TestPath("/test3")

	assert.False(t, result)
}

func TestMatchPathWithVariable(t *testing.T) {
	conf := &Config{
		Enabled: true,
		Verbose: false,
		TTL:     "1m",
		Paths: []string{
			"/test1",
			"/test2/:id",
		},
		AdditionalHeaders: map[string]string{},
	}
	cm := NewCacheManager(conf, nil)

	result := cm.TestPath("/test2/123")

	assert.True(t, result)
}

func TestMatchPathWithVariableUnmatchFromMissingVariable(t *testing.T) {
	conf := &Config{
		Enabled: true,
		Verbose: false,
		TTL:     "1m",
		Paths: []string{
			"/test1",
			"/test2/:id",
		},
		AdditionalHeaders: map[string]string{},
	}
	cm := NewCacheManager(conf, nil)

	result := cm.TestPath("/test2")

	assert.False(t, result)
}

func TestSetCacheShouldInvokeSetCache(t *testing.T) {
	mockCache := new(MockCache)
	expectedContent := []byte{1, 2, 3, 4}
	conf := &Config{
		Enabled: true,
		Verbose: false,
		TTL:     "1m",
		Paths: []string{
			"/test1",
			"/test2",
		},
		AdditionalHeaders: map[string]string{},
	}
	mockCache.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
	cm := NewCacheManager(conf, mockCache)

	cm.Set("/test2", expectedContent)

	mockCache.AssertCalled(t, "Set", "/test2", expectedContent)
}

func TestGetCacheShouldInvokeGetCache(t *testing.T) {
	mockCache := new(MockCache)
	expectedContent := []byte{1, 2, 3, 4}
	conf := &Config{
		Enabled: true,
		Verbose: false,
		TTL:     "1m",
		Paths: []string{
			"/test1",
			"/test2",
		},
		AdditionalHeaders: map[string]string{},
	}
	mockCache.On("Get", mock.AnythingOfType("string")).Return(expectedContent, nil)
	cm := NewCacheManager(conf, mockCache)

	result, ok := cm.Get("/test2")

	mockCache.AssertCalled(t, "Get", "/test2")
	assert.Equal(t, expectedContent, result)
	assert.True(t, ok)
}

func TestGetCacheShouldReturnFalseWhenCacheMisses(t *testing.T) {
	mockCache := new(MockCache)
	expectedContent := []byte{}
	conf := &Config{
		Enabled: true,
		Verbose: false,
		TTL:     "1m",
		Paths: []string{
			"/test1",
			"/test2",
		},
		AdditionalHeaders: map[string]string{},
	}
	mockCache.On("Get", mock.AnythingOfType("string")).Return(expectedContent, errors.New("Error"))
	cm := NewCacheManager(conf, mockCache)

	result, ok := cm.Get("/test2")

	mockCache.AssertCalled(t, "Get", "/test2")
	assert.Equal(t, expectedContent, result)
	assert.False(t, ok)
}

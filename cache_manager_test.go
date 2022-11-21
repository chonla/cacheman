package cacheman

import (
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

func (o *MockCache) Delete(k string) error {
	args := o.Called(k)
	return args.Error(0)
}

func (o *MockCache) Reset() error {
	args := o.Called()
	return args.Error(0)
}

func TestMatchPathWithWildcard(t *testing.T) {
	conf := &Config{
		Enabled: true,
		Verbose: false,
		TTL:     "1m",
		Paths: []string{
			"/.*",
		},
		AdditionalHeaders: map[string]string{},
	}
	cm := NewCacheManager(conf, nil)

	result := cm.TestPath("/test2")

	assert.True(t, result)
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

func TestMatchPathWithWildcardAndExcludsion(t *testing.T) {
	conf := &Config{
		Enabled: true,
		Verbose: false,
		TTL:     "1m",
		Paths: []string{
			"/.*",
		},
		ExcludedPaths: []string{
			"/test2",
		},
		AdditionalHeaders: map[string]string{},
	}
	cm := NewCacheManager(conf, nil)

	result := cm.TestPath("/test2")

	assert.False(t, result)
}

func TestMatchPathWithWildcardAndExcludsionWildcard(t *testing.T) {
	conf := &Config{
		Enabled: true,
		Verbose: false,
		TTL:     "1m",
		Paths: []string{
			"/.*",
		},
		ExcludedPaths: []string{
			"/.*",
		},
		AdditionalHeaders: map[string]string{},
	}
	cm := NewCacheManager(conf, nil)

	result := cm.TestPath("/")

	assert.False(t, result)
}

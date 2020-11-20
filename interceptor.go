package cacheman

import (
	"net/http"
)

// Interceptor is response interceptor
type Interceptor struct {
	writer    http.ResponseWriter
	committed bool

	status  int
	header  http.Header
	content []byte
}

// NewInterceptor creates a new response interceptor
func NewInterceptor(writer http.ResponseWriter) *Interceptor {
	return &Interceptor{
		writer: writer,
	}
}

// Header returns response header
func (c *Interceptor) Header() http.Header {
	c.header = c.writer.Header()
	return c.header
}

// Write writes out the content. Automatically writes out the header if it has not been written out.
func (c *Interceptor) Write(b []byte) (int, error) {
	if !c.committed {
		c.WriteHeader(http.StatusOK)
	}
	c.content = b
	return c.writer.Write(c.content)
}

// WriteHeader writes out the header with given status code
func (c *Interceptor) WriteHeader(statusCode int) {
	if c.committed {
		return
	}
	c.status = statusCode
	c.writer.WriteHeader(c.status)
	c.committed = true
}

// Status returns the captured status
func (c *Interceptor) Status() int {
	return c.status
}

// Content returns the captured content
func (c *Interceptor) Content() []byte {
	return c.content
}

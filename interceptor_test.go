package cacheman

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockResponseWriter struct {
	mock.Mock
}

func (o *MockResponseWriter) Header() http.Header {
	args := o.Called()
	return args.Get(0).(http.Header)
}

func (o *MockResponseWriter) Write(b []byte) (int, error) {
	args := o.Called(b)
	return args.Int(0), args.Error(1)
}

func (o *MockResponseWriter) WriteHeader(code int) {
	o.Called(code)
}

func TestGetHeaderShouldGetHeaderFromWriter(t *testing.T) {
	mockWriter := new(MockResponseWriter)
	var mockReturnHeader http.Header = map[string][]string{
		"Content-Type": {
			"application/octet-stream",
		},
		"Some-Other-Header": {
			"passed",
		},
	}

	mockWriter.On("Header").Return(mockReturnHeader)

	interceptor := NewInterceptor(mockWriter)

	header := interceptor.Header()

	mockWriter.AssertNumberOfCalls(t, "Header", 1)
	assert.Equal(t, "application/octet-stream", header.Get("Content-Type"))
	assert.Equal(t, "passed", header.Get("Some-Other-Header"))
}

func TestWriteHeaderShouldCacheStatusCode(t *testing.T) {
	mockWriter := new(MockResponseWriter)
	expectedCode := 204

	mockWriter.On("WriteHeader", mock.AnythingOfType("int"))

	interceptor := NewInterceptor(mockWriter)

	interceptor.WriteHeader(expectedCode)

	mockWriter.AssertNumberOfCalls(t, "WriteHeader", 1)
	mockWriter.AssertCalled(t, "WriteHeader", expectedCode)
	assert.Equal(t, expectedCode, interceptor.Status())
}

func TestWriteHeaderShouldNotReinvokeWriteHeader(t *testing.T) {
	mockWriter := new(MockResponseWriter)
	expectedCode := 204

	mockWriter.On("WriteHeader", mock.AnythingOfType("int"))

	interceptor := NewInterceptor(mockWriter)

	interceptor.WriteHeader(expectedCode)
	interceptor.WriteHeader(200)

	mockWriter.AssertNumberOfCalls(t, "WriteHeader", 1)
	mockWriter.AssertCalled(t, "WriteHeader", expectedCode)
	assert.Equal(t, expectedCode, interceptor.Status())
}

func TestWriteWithoutWritingHeaderShouldCacheContentWithHttp200OK(t *testing.T) {
	mockWriter := new(MockResponseWriter)
	expectedCode := 200
	expectedContent := []byte{1, 2, 3, 4}

	mockWriter.On("WriteHeader", mock.AnythingOfType("int"))
	mockWriter.On("Write", mock.AnythingOfType("[]uint8")).Return(4, nil)

	interceptor := NewInterceptor(mockWriter)

	interceptor.Write(expectedContent)

	mockWriter.AssertNumberOfCalls(t, "WriteHeader", 1)
	mockWriter.AssertNumberOfCalls(t, "Write", 1)
	mockWriter.AssertCalled(t, "WriteHeader", expectedCode)
	mockWriter.AssertCalled(t, "Write", expectedContent)
	assert.Equal(t, expectedContent, interceptor.Content())
}

func TestWriteWithWritingHeaderEarlierShouldCacheContentWithHttp200OK(t *testing.T) {
	mockWriter := new(MockResponseWriter)
	expectedCode := 201
	expectedContent := []byte{1, 2, 3, 4}

	mockWriter.On("WriteHeader", mock.AnythingOfType("int"))
	mockWriter.On("Write", mock.AnythingOfType("[]uint8")).Return(4, nil)

	interceptor := NewInterceptor(mockWriter)

	interceptor.WriteHeader(expectedCode)
	interceptor.Write(expectedContent)

	mockWriter.AssertNumberOfCalls(t, "WriteHeader", 1)
	mockWriter.AssertNumberOfCalls(t, "Write", 1)
	mockWriter.AssertCalled(t, "WriteHeader", expectedCode)
	mockWriter.AssertCalled(t, "Write", expectedContent)
	assert.Equal(t, expectedContent, interceptor.Content())
}

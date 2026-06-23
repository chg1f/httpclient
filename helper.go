package httpclient

import "net/http"

// RoundTripper adapts a function to the http.RoundTripper interface.
type RoundTripper func(*http.Request) (*http.Response, error)

// RoundTrip executes fn for req.
func (fn RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

// Closer adapts a function to the io.Closer interface.
type Closer func() error

// Close executes fn.
func (fn Closer) Close() error {
	return fn()
}

// Reader adapts a function to the io.Reader interface.
type Reader func(bs []byte) (int, error)

// Read executes fn with bs.
func (fn Reader) Read(bs []byte) (int, error) {
	return fn(bs)
}

// ReadCloser implements io.ReadCloser by combining Reader and Closer.
type ReadCloser struct {
	Reader
	Closer
}

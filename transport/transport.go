package transport

import "net/http"

// RoundTripper adapts a function to the http.RoundTripper interface.
type RoundTripper func(*http.Request) (*http.Response, error)

// RoundTrip executes an HTTP request.
func (fn RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func New(tr http.RoundTripper, opts ...Option) http.RoundTripper {
	for _, apply := range opts {
		tr = apply(tr)
	}
	return tr
}

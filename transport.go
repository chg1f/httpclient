package httpclient

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// Transport wraps an http.RoundTripper with additional behavior.
type Transport func(http.RoundTripper) http.RoundTripper

// Timeout applies a request timeout at the transport layer.
func Timeout(ttl time.Duration) Transport {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripper(func(req *http.Request) (*http.Response, error) {
			ctx, cancel := context.WithTimeout(req.Context(), ttl)
			defer cancel()

			return next.RoundTrip(req.Clone(ctx))
		})
	}
}

// DownstreamBandwidth limits response body read bandwidth with limiter.
func DownstreamBandwidth(limiter *rate.Limiter) Transport {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripper(func(req *http.Request) (*http.Response, error) {
			resp, err := next.RoundTrip(req)
			if resp != nil && resp.Body != nil {
				body := resp.Body
				resp.Body = ReadCloser{
					Reader: Reader(func(bs []byte) (int, error) {
						n, err := body.Read(bs)
						if n > 0 {
							err = limiter.WaitN(req.Context(), n)
							if err != nil {
								return n, err
							}
						}
						return n, err
					}),
					Closer: Closer(func() error {
						return body.Close()
					}),
				}
			}
			return resp, err
		})
	}
}

// UpstreamBandwidth limits request body read bandwidth with limiter.
func UpstreamBandwidth(limiter *rate.Limiter) Transport {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripper(func(req *http.Request) (*http.Response, error) {
			if req.Body != nil {
				body := req.Body
				req.Body = ReadCloser{
					Reader: Reader(func(bs []byte) (int, error) {
						n, err := body.Read(bs)
						if n > 0 {
							err = limiter.WaitN(req.Context(), n)
							if err != nil {
								return n, err
							}
						}
						return n, err
					}),
					Closer: Closer(func() error {
						return body.Close()
					}),
				}
			}
			return next.RoundTrip(req)
		})
	}
}

// Ratelimit limits request throughput with limiter.
func Ratelimit(limiter *rate.Limiter) Transport {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripper(func(req *http.Request) (*http.Response, error) {
			if err := limiter.Wait(req.Context()); err != nil {
				return nil, err
			}
			return next.RoundTrip(req)
		})
	}
}

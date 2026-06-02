package transport

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// Option wraps an http.RoundTripper with additional behavior.
type Option func(http.RoundTripper) http.RoundTripper

// Timeout applies a request timeout at the transport layer.
func Timeout(ttl time.Duration) Option {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripper(func(req *http.Request) (*http.Response, error) {
			ctx, cancel := context.WithTimeout(req.Context(), ttl)
			defer cancel()

			return next.RoundTrip(req.Clone(ctx))
		})
	}
}

// DownstreamBandwidth limits response body read bandwidth.
func DownstreamBandwidth(size uint64) Option {
	limiter := rate.NewLimiter(rate.Limit(size), int(size))
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripper(func(req *http.Request) (*http.Response, error) {
			resp, err := next.RoundTrip(req)
			if resp != nil && resp.Body != nil {
				resp.Body = ReadCloser{
					Context:    req.Context(),
					ReadCloser: resp.Body,
					Limiter:    limiter,
				}
			}
			return resp, err
		})
	}
}

// UpstreamBandwidth limits request body read bandwidth.
func UpstreamBandwidth(size uint64) Option {
	limiter := rate.NewLimiter(rate.Limit(size), int(size))
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripper(func(req *http.Request) (*http.Response, error) {
			if req.Body != nil {
				req.Body = ReadCloser{
					Context:    req.Context(),
					ReadCloser: req.Body,
					Limiter:    limiter,
				}
			}
			return next.RoundTrip(req)
		})
	}
}

// Ratelimit limits the number of requests allowed within a period.
func Ratelimit(period time.Duration, limit int) Option {
	limiter := rate.NewLimiter(rate.Every(period), limit)
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripper(func(req *http.Request) (*http.Response, error) {
			if err := limiter.Wait(req.Context()); err != nil {
				return nil, err
			}
			return next.RoundTrip(req)
		})
	}
}

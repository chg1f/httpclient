package transport

import (
	"context"
	"io"

	"golang.org/x/time/rate"
)

// ReadCloser applies a rate-limit bucket while reading from a stream.
type ReadCloser struct {
	context.Context
	io.ReadCloser
	*rate.Limiter
}

// Read reads bytes from the stream and waits for bucket capacity for bytes read.
func (rc ReadCloser) Read(bs []byte) (int, error) {
	n, err := rc.ReadCloser.Read(bs)
	if n > 0 {
		err = rc.Limiter.WaitN(rc.Context, n)
		if err != nil {
			return n, err
		}
	}
	return n, err
}

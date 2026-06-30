package httpclient

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDownstreamBandwidthReadsWrappedResponseBody(t *testing.T) {
	next := RoundTripper(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("response")),
		}, nil
	})

	resp, err := DownstreamBandwidth(1024)(next).RoundTrip(&http.Request{})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "response" {
		t.Fatalf("body = %q, want %q", bs, "response")
	}
}

func TestUpstreamBandwidthReadsWrappedRequestBody(t *testing.T) {
	next := RoundTripper(func(req *http.Request) (*http.Response, error) {
		bs, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		if string(bs) != "request" {
			t.Fatalf("body = %q, want %q", bs, "request")
		}
		return &http.Response{StatusCode: http.StatusOK}, nil
	})

	req, err := http.NewRequest(http.MethodPost, "http://example.test", strings.NewReader("request"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = UpstreamBandwidth(1024)(next).RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
}

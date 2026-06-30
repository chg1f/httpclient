# httpclient

`httpclient` is a configurable Go wrapper around `net/http.Client`.

The package is intended to keep the standard library HTTP client interface while
making common client and transport behavior easier to compose, configure, and
reuse.

## Installation

```sh
go get github.com/chg1f/httpclient
```

## Client

Use `Default` when the standard transport defaults are enough:

```go
resp, err := httpclient.Default.Get("https://example.com")
if err != nil {
	return err
}
defer resp.Body.Close()
```

Create a client when request timeout or transport settings need to be configured:

```go
cfg := httpclient.Config{}
cfg.Timeout = 10 * time.Second
cfg.Transport.Dial.Timeout = 3 * time.Second
cfg.Transport.Dial.KeepAlive = 30 * time.Second
cfg.Transport.ForceAttemptHttp2 = true
cfg.Transport.MaxIdleConns = 100
cfg.Transport.MaxIdleConnsPerHost = 10

client := httpclient.NewClient(cfg)
```

`Config` fields include `json`, `yaml`, and `mapstructure` tags so the same
structure can be populated by common configuration loaders.

## Transport Options

`httpclient` wraps `http.RoundTripper` implementations with small, composable
options:

```go
limiter := rate.NewLimiter(rate.Every(time.Second), 100)

client.Transport = httpclient.New(
	http.DefaultTransport,
	httpclient.Timeout(5*time.Second),
	httpclient.Ratelimit(limiter),
)
```

`New` applies options in the supplied order. Available options:

- `Timeout` applies a timeout at the transport layer.
- `DownstreamBandwidth` limits response body read bandwidth with a limiter.
- `UpstreamBandwidth` limits request body read bandwidth with a limiter.
- `Ratelimit` limits request throughput with a limiter.

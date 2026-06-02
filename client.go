package httpclient

import (
	"net"
	"net/http"
	"runtime"
	"time"
)

// Default is the package-level client configured with standard transport defaults.
var Default *http.Client

func init() {
	var cfg Config
	cfg.Transport.Dial.Timeout = time.Second * 30
	cfg.Transport.Dial.KeepAlive = time.Second * 30
	cfg.Transport.Dial.DualStack = true
	cfg.Transport.MaxIdleConns = 100
	cfg.Transport.IdleConnTimeout = time.Second * 90
	cfg.Transport.TlsHandshakeTimeout = time.Second * 10
	cfg.Transport.ExpectContinueTimeout = time.Second * 1
	cfg.Transport.ForceAttemptHttp2 = true
	cfg.Transport.MaxIdleConnsPerHost = runtime.GOMAXPROCS(0) + 1
	Default = NewClient(cfg)
}

// Config controls how a Client is created.
type Config struct {
	// Timeout limits the total duration of a request, including redirects and body reads.
	Timeout time.Duration `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	// Transport controls the underlying HTTP transport.
	Transport struct {
		// Dial controls TCP connection establishment.
		Dial struct {
			// Timeout is the maximum time allowed for a connection attempt.
			Timeout time.Duration `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
			// KeepAlive is the interval between keep-alive probes for an active connection.
			KeepAlive time.Duration `json:"keep_alive" yaml:"keep_alive" mapstructure:"keep_alive"`
			// DualStack enables RFC 6555 fast fallback between IPv4 and IPv6.
			DualStack bool `json:"dual_stack" yaml:"dual_stack" mapstructure:"dual_stack"`
			// FallbackDelay is the delay before falling back to another address family.
			FallbackDelay time.Duration `json:"fallback_delay" yaml:"fallback_delay" mapstructure:"fallback_delay"`
		} `json:"dial" yaml:"dial" mapstructure:"dial"`
		// DisableKeepAlives disables HTTP keep-alives and uses each connection for a single request.
		DisableKeepAlives bool `json:"disable_keep_alives" yaml:"disable_keep_alives" mapstructure:"disable_keep_alives"`
		// DisableCompression prevents transparent response decompression.
		DisableCompression bool `json:"disable_compression" yaml:"disable_compression" mapstructure:"disable_compression"`
		// ForceAttemptHttp2 controls whether HTTP/2 is attempted when supported.
		ForceAttemptHttp2 bool `json:"force_attempt_http2" yaml:"force_attempt_http2" mapstructure:"force_attempt_http2"`
		// TlsHandshakeTimeout is the maximum time allowed for a TLS handshake.
		TlsHandshakeTimeout time.Duration `json:"tls_handshake_timeout" yaml:"tls_handshake_timeout" mapstructure:"tls_handshake_timeout"`
		// ResponseHeaderTimeout limits how long to wait for response headers after writing a request.
		ResponseHeaderTimeout time.Duration `json:"response_header_timeout" yaml:"response_header_timeout" mapstructure:"response_header_timeout"`
		// ExpectContinueTimeout limits how long to wait for a server's first response headers after Expect: 100-continue.
		ExpectContinueTimeout time.Duration `json:"expect_continue_timeout" yaml:"expect_continue_timeout" mapstructure:"expect_continue_timeout"`
		// IdleConnTimeout is the maximum time an idle connection remains open.
		IdleConnTimeout time.Duration `json:"idle_conn_timeout" yaml:"idle_conn_timeout" mapstructure:"idle_conn_timeout"`
		// MaxIdleConns limits the number of idle connections across all hosts.
		MaxIdleConns int `json:"max_idle_conns" yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
		// MaxIdleConnsPerHost limits the number of idle connections per host.
		MaxIdleConnsPerHost int `json:"max_idle_conns_per_host" yaml:"max_idle_conns_per_host" mapstructure:"max_idle_conns_per_host"`
		// MaxConnsPerHost limits the total number of connections per host.
		MaxConnsPerHost int `json:"max_conns_per_host" yaml:"max_conns_per_host" mapstructure:"max_conns_per_host"`
		// WriteBufferSize sets the transport write buffer size.
		WriteBufferSize int `json:"write_buffer_size" yaml:"write_buffer_size" mapstructure:"write_buffer_size"`
		// ReadBufferSize sets the transport read buffer size.
		ReadBufferSize int `json:"read_buffer_size" yaml:"read_buffer_size" mapstructure:"read_buffer_size"`
	} `json:"transport" yaml:"transport" mapstructure:"transport"`
}

// NewClient creates a configured standard HTTP client.
func NewClient(cfg Config) *http.Client {
	return &http.Client{
		Timeout: cfg.Timeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:       cfg.Transport.Dial.Timeout,
				KeepAlive:     cfg.Transport.Dial.KeepAlive,
				DualStack:     cfg.Transport.Dial.DualStack,
				FallbackDelay: cfg.Transport.Dial.FallbackDelay,
			}).DialContext,
			DisableKeepAlives:     cfg.Transport.DisableKeepAlives,
			DisableCompression:    cfg.Transport.DisableCompression,
			ForceAttemptHTTP2:     cfg.Transport.ForceAttemptHttp2,
			TLSHandshakeTimeout:   cfg.Transport.TlsHandshakeTimeout,
			ExpectContinueTimeout: cfg.Transport.ExpectContinueTimeout,
			ResponseHeaderTimeout: cfg.Transport.ResponseHeaderTimeout,
			IdleConnTimeout:       cfg.Transport.IdleConnTimeout,
			MaxIdleConns:          cfg.Transport.MaxIdleConns,
			MaxIdleConnsPerHost:   cfg.Transport.MaxIdleConnsPerHost,
			MaxConnsPerHost:       cfg.Transport.MaxConnsPerHost,
			WriteBufferSize:       cfg.Transport.WriteBufferSize,
			ReadBufferSize:        cfg.Transport.ReadBufferSize,
		},
	}
}

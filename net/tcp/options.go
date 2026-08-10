package tcp

import (
	"crypto/tls"
	"time"
)

const (
	defaultPort = "2525"
)

// ServerOptions holds configuration settings for a TCP server.
type ServerOptions struct {
	Addresses       []string
	KeepAlivePeriod time.Duration
	RetryInitial    time.Duration
	RetryMax        time.Duration
	RetryAttempts   int
	TLSConfig       *tls.Config
}

// ServerOption defines a function signature for setting ServerOptions.
type ServerOption func(*ServerOptions)

func defaultServerOptions() *ServerOptions {
	return &ServerOptions{
		Addresses:       []string{":" + defaultPort},
		KeepAlivePeriod: time.Second * 30,
		RetryInitial:    time.Millisecond * 500,
		RetryMax:        5 * time.Minute,
		RetryAttempts:   5,
	}
}

// NewServerOptions creates ServerOptions initialized with defaults and updated by options.
func NewServerOptions(opts ...ServerOption) *ServerOptions {
	o := defaultServerOptions()
	for _, opt := range opts {
		opt(o)
	}

	return o
}

// WithAddresses returns a ServerOption that sets the listening addresses for the server.
func WithAddresses(addrs ...string) ServerOption {
	return func(o *ServerOptions) {
		o.Addresses = addrs
	}
}

// ConnOptions holds configuration settings for TCP connections.
type ConnOptions struct {
	DialRetries       int
	DialRetryDelay    time.Duration
	DialTimeout       time.Duration
	KeepAlive         time.Duration
	KeepAlivePeriod   time.Duration
	ReceiveTimeout    time.Duration
	SendTimeout       time.Duration
	ReceiveBufferWait time.Duration
}

// ConnOption defines a function signature for setting ConnOptions.
type ConnOption func(*ConnOptions)

// DefaultConnOptions returns the default ConnOptions.
func DefaultConnOptions() *ConnOptions {
	return &ConnOptions{
		DialRetries:       3,
		DialRetryDelay:    time.Second,
		DialTimeout:       5 * time.Second,
		ReceiveBufferWait: time.Millisecond,
	}
}

// NewConnOptions creates ConnOptions initialized with defaults and updated by options.
func NewConnOptions(opts ...ConnOption) *ConnOptions {
	o := DefaultConnOptions()
	for _, opt := range opts {
		opt(o)
	}

	return o
}

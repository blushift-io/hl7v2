package tcp

import (
	"crypto/tls"
	"time"
)

const (
	defaultPort = "2525"
)

type ServerOptions struct {
	Addresses       []string
	KeepAlivePeriod time.Duration
	RetryInitial    time.Duration
	RetryMax        time.Duration
	RetryAttempts   int
	TLSConfig       *tls.Config
}

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

func NewServerOptions(opts ...ServerOption) *ServerOptions {
	o := defaultServerOptions()
	for _, opt := range opts {
		opt(o)
	}

	return o
}

func WithAddresses(addrs ...string) ServerOption {
	return func(o *ServerOptions) {
		o.Addresses = addrs
	}
}

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

type ConnOption func(*ConnOptions)

func DefaultConnOptions() *ConnOptions {
	return &ConnOptions{
		DialRetries:       3,
		DialRetryDelay:    time.Second,
		DialTimeout:       5 * time.Second,
		ReceiveBufferWait: time.Millisecond,
	}
}

func NewConnOptions(opts ...ConnOption) *ConnOptions {
	o := DefaultConnOptions()
	for _, opt := range opts {
		opt(o)
	}

	return o
}

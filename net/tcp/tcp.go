package tcp

import (
	"context"
	"time"

	"github.com/blushift-io/hl7v2"
)

type SendOptions struct {
	WaitForAck bool
	Timeout    time.Duration
}

type SendOption func(*SendOptions)

func WaitForAck() SendOption {
	return func(o *SendOptions) {
		o.WaitForAck = true
	}
}

func WithTimeout(d time.Duration) SendOption {
	return func(o *SendOptions) {
		o.Timeout = d
	}
}

func Send(ctx context.Context, host string, msg *hl7v2.RawMessage, opts ...SendOption) (*hl7v2.RawMessage, error) {
	options := &SendOptions{
		WaitForAck: false,
		Timeout:    30 * time.Second,
	}

	for _, opt := range opts {
		opt(options)
	}

	conn, err := Dial(ctx, host)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err := conn.WriteMessage(msg); err != nil {
		return nil, err
	}

	if !options.WaitForAck {
		return nil, nil
	}

	ackCh := make(chan *hl7v2.RawMessage, 1)
	errCh := make(chan error, 1)

	go func() {
		ack, err := conn.ReadMessage()
		if err != nil {
			errCh <- err
			return
		}
		ackCh <- ack
	}()

	select {
	case ack := <-ackCh:
		return ack, nil
	case err := <-errCh:
		return nil, err
	case <-time.After(options.Timeout):
		return nil, context.DeadlineExceeded
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

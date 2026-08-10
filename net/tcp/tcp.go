// Package tcp implements TCP networking for HL7 v2 messages, including MLLP framing,
// connection management, message handling, and server lifecycle options.
package tcp

import (
	"context"
	"fmt"
	"time"

	"github.com/blushift-io/hl7v2"
)

// SendOptions holds options for sending an HL7 v2 message over TCP.
type SendOptions struct {
	WaitForAck bool
	Timeout    time.Duration
}

// SendOption defines a function signature for setting SendOptions.
type SendOption func(*SendOptions)

// WaitForAck returns a SendOption that specifies whether to wait for an acknowledgment response.
func WaitForAck() SendOption {
	return func(o *SendOptions) {
		o.WaitForAck = true
	}
}

// WithTimeout returns a SendOption that configures the operation timeout.
func WithTimeout(d time.Duration) SendOption {
	return func(o *SendOptions) {
		o.Timeout = d
	}
}

// Send connects to a target host, sends a raw HL7 v2 message over MLLP, and optionally waits for an acknowledgment.
func Send(ctx context.Context, host string, msg *hl7v2.RawMessage, opts ...SendOption) (*hl7v2.RawMessage, error) {
	if msg == nil {
		return nil, fmt.Errorf("tcp: message cannot be nil")
	}

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

	timer := time.NewTimer(options.Timeout)
	defer timer.Stop()

	select {
	case ack := <-ackCh:
		return ack, nil
	case err := <-errCh:
		return nil, err
	case <-timer.C:
		return nil, context.DeadlineExceeded
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

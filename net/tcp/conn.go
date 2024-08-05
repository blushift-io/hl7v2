package tcp

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/mllp"
)

type Conn interface {
	net.Conn
	WriteMessage(*hl7v2.RawMessage) error
	ReadMessage() (*hl7v2.RawMessage, error)
	AckMessage(*hl7v2.RawMessage) error
}

func Dial(ctx context.Context, addr string, opts ...ConnOption) (Conn, error) {
	options := NewConnOptions(opts...)

	conn, err := tryConnect(addr, options)
	if err != nil {
		return nil, err
	}

	return newConn(conn, options)
}

type tcpConn struct {
	net.Conn
	cancel context.CancelFunc
}

func newConn(conn net.Conn, opts *ConnOptions) (*tcpConn, error) {
	nc := conn.(*net.TCPConn)
	if err := nc.SetKeepAlive(true); err != nil {
		return nil, err
	}

	if err := nc.SetKeepAlivePeriod(opts.KeepAlivePeriod); err != nil {
		return nil, err
	}

	return &tcpConn{
		Conn: nc,
	}, nil
}

func (c *tcpConn) WriteMessage(m *hl7v2.RawMessage) error {
	b := m.Value().Bytes()

	if err := mllp.Write(c, b); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (c *tcpConn) ReadMessage() (*hl7v2.RawMessage, error) {
	b, err := mllp.Read(c)
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	m, err := hl7v2.NewRawMessageFromBytes(b, hl7v2.FixLineEndings())
	if err != nil {
		return nil, fmt.Errorf("failed to parse message: %w", err)
	}

	return m, nil
}

func (c *tcpConn) AckMessage(m *hl7v2.RawMessage) error {
	//TODO: better implement this silly thing
	v := m.Value().Bytes()
	ack, err := hl7v2.AckRawMessage(v)
	if err != nil {
		return fmt.Errorf("failed to ack message: %w", err)
	}

	if err := mllp.Write(c, ack); err != nil {
		return fmt.Errorf("failed to write ack: %w", err)
	}

	return nil
}

func (c *tcpConn) Close() error {
	if c.cancel != nil {
		c.cancel()
	}

	return c.Conn.Close()
}

func tryConnect(addr string, opts *ConnOptions) (net.Conn, error) {
	tries := 0
	retry := &expBackoffRetry{
		Count:        opts.DialRetries,
		InitialDelay: opts.DialRetryDelay,
		MaxDelay:     opts.DialRetryDelay * 10,
	}

	for {
		conn, err := net.DialTimeout("tcp", addr, opts.DialTimeout)
		if err != nil {
			d, ok := retry.Backoff(uint64(tries))
			if !ok {
				return nil, err
			}

			time.Sleep(d)
			continue
		}

		return conn, nil
	}
}

package tcp

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"

	"github.com/blushift-io/hl7v2"
)

type echoHandler struct {
	received chan *hl7v2.RawMessage
}

func (h *echoHandler) OnConnect(c *Context) error {
	return nil
}

func (h *echoHandler) OnMessage(c *Context) error {
	if msg, ok := c.Data.(*hl7v2.RawMessage); ok {
		if h.received != nil {
			h.received <- msg
		}
		if c.Conn != nil {
			ackBytes, err := hl7v2.AckRawMessage(msg.Value().Bytes())
			if err == nil {
				rawAck, err := hl7v2.ParseRaw(ackBytes)
				if err == nil {
					_ = c.Conn.WriteMessage(rawAck)
				}
			}
		}
	}
	return nil
}

func (h *echoHandler) OnError(c *Context) error {
	return nil
}

func (h *echoHandler) OnClose(c *Context) error {
	return nil
}

type nackHandler struct {
	nackSent chan *hl7v2.RawMessage
}

func (h *nackHandler) OnConnect(c *Context) error { return nil }
func (h *nackHandler) OnMessage(c *Context) error {
	if msg, ok := c.Data.(*hl7v2.RawMessage); ok {
		if c.Conn != nil {
			nack, err := c.Conn.NackMessage(msg, hl7v2.AckApplicationError, "Invalid processing ID", hl7v2.MessageErrorApplicationInternalError)
			if err == nil && h.nackSent != nil {
				h.nackSent <- nack
			}
		}
	}
	return nil
}
func (h *nackHandler) OnError(c *Context) error { return nil }
func (h *nackHandler) OnClose(c *Context) error { return nil }

func TestServerLifecycle(t *testing.T) {
	h := &echoHandler{
		received: make(chan *hl7v2.RawMessage, 10),
	}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	addr := l.Addr().String()

	srv, err := NewServer(h, WithAddresses(addr))
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	go func() {
		_ = srv.Serve(l)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := Dial(ctx, addr, func(o *ConnOptions) {
		o.DialRetries = 1
	})
	if err != nil {
		t.Fatalf("failed to dial server: %v", err)
	}

	rawMsg, err := hl7v2.ParseRaw([]byte("MSH|^~\\&|SEND|FAC|REC|FAC|20260101||ADT^A01|123|P|2.5\rPID|1||12345\r"))
	if err != nil {
		t.Fatalf("failed to parse message: %v", err)
	}

	if err := conn.WriteMessage(rawMsg); err != nil {
		t.Fatalf("failed to write message: %v", err)
	}

	select {
	case msg := <-h.received:
		if msg == nil {
			t.Errorf("received nil message on server")
		}
	case <-time.After(2 * time.Second):
		t.Errorf("timeout waiting for server to receive message")
	}

	_ = conn.Close()

	if err := srv.Shutdown(); err != nil {
		t.Errorf("failed to shutdown server cleanly: %v", err)
	}
}

func TestNackMessage(t *testing.T) {
	nh := &nackHandler{
		nackSent: make(chan *hl7v2.RawMessage, 10),
	}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	addr := l.Addr().String()

	srv, err := NewServer(nh, WithAddresses(addr))
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	go func() {
		_ = srv.Serve(l)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := Dial(ctx, addr, func(o *ConnOptions) {
		o.DialRetries = 1
	})
	if err != nil {
		t.Fatalf("failed to dial server: %v", err)
	}

	rawMsg, err := hl7v2.ParseRaw([]byte("MSH|^~\\&|SEND|FAC|REC|FAC|20260101||ADT^A01|MSG999|P|2.5\rPID|1||12345\r"))
	if err != nil {
		t.Fatalf("failed to parse message: %v", err)
	}

	if err := conn.WriteMessage(rawMsg); err != nil {
		t.Fatalf("failed to write message: %v", err)
	}

	resp, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read NACK response: %v", err)
	}

	respStr := string(resp.Value().Bytes())
	if !bytes.Contains(resp.Value().Bytes(), []byte("MSA|AE|MSG999|Invalid processing ID|||207")) {
		t.Errorf("expected MSA|AE|MSG999|Invalid processing ID in NACK response, got: %s", respStr)
	}

	_ = conn.Close()
	_ = srv.Shutdown()
}

func TestNilHandlerValidation(t *testing.T) {
	_, err := NewServer(nil)
	if err == nil {
		t.Errorf("expected error when creating server with nil handler, got nil")
	}
}

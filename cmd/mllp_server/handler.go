package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/net/tcp"
)

type handler struct {
	savePath string
}

func (h *handler) OnConnect(ctx *tcp.Context) error {
	log.Printf("Client connected: %s", ctx.Conn.RemoteAddr())

	return nil
}

func (h *handler) OnMessage(ctx *tcp.Context) error {
	log.Printf("Received message from %s", ctx.Conn.RemoteAddr())

	msg, ok := ctx.Data.(*hl7v2.RawMessage)
	if !ok {
		return fmt.Errorf("invalid hl7 raw message, got type: %T", ctx.Data)
	}

	cid, err := msg.QueryValue("MSH.10")
	if err != nil {
		return fmt.Errorf("error querying MSH.10: %w", err)
	}

	b := msg.Value().Bytes()
	fn := filepath.Join(h.savePath, fmt.Sprintf("%s.hl7", cid.String()))

	if err := os.WriteFile(fn, b, 0644); err != nil {
		return fmt.Errorf("error writing hl7 message to file: %w", err)
	}

	ack, err := ctx.Conn.AckMessage(msg)
	if err != nil {
		return fmt.Errorf("error sending ack: %w", err)
	}

	ab := ack.Value().Bytes()
	afn := filepath.Join(h.savePath, fmt.Sprintf("%s-ACK.hl7", cid.String()))
	if err := os.WriteFile(afn, ab, 0644); err != nil {
		return fmt.Errorf("error writing hl7 ack message to file: %w", err)
	}

	return nil
}

func (h *handler) OnError(ctx *tcp.Context) error {
	log.Printf("Error: %s", ctx.Error)

	return nil
}

func (h *handler) OnClose(ctx *tcp.Context) error {
	log.Printf("Client disconnected: %s", ctx.Conn.RemoteAddr())

	return nil
}

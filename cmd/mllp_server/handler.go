package main

import (
	"errors"
	"log"

	"github.com/blushift-io/hl7v2/net/tcp"
)

var (
	ErrInvalidMsg = errors.New("invalid hl7 raw message")
)

type handler struct {
	actions []action
}

func newHandler(conf []Action) (*handler, error) {
	h := &handler{}
	for _, a := range conf {
		log.Printf("Creating action: %s", a.Type)
		a, err := newAction(a)
		if err != nil {
			return nil, err
		}
		h.actions = append(h.actions, a)
	}

	return h, nil
}

func (h *handler) OnConnect(ctx *tcp.Context) error {
	log.Printf("Client connected: %s", ctx.Conn.RemoteAddr())

	var errs []error
	for _, a := range h.actions {
		err := a.OnConnect(ctx)
		if err != nil {
			if aerr := a.OnError(ctx); aerr != nil {
				errs = append(errs, aerr)
			}
		}
	}

	if len(errs) > 1 {
		return errors.Join(errs...)
	}

	return nil
}

func (h *handler) OnMessage(ctx *tcp.Context) error {
	log.Printf("Received message from %s", ctx.Conn.RemoteAddr())

	var errs []error
	for _, a := range h.actions {
		err := a.OnMessage(ctx)
		if err != nil {
			if aerr := a.OnError(ctx); aerr != nil {
				errs = append(errs, aerr)
			}
		}
	}

	if len(errs) > 1 {
		return errors.Join(errs...)
	}

	return nil
}

func (h *handler) OnError(ctx *tcp.Context) error {
	log.Printf("Error: %s", ctx.Error)

	return nil
}

func (h *handler) OnClose(ctx *tcp.Context) error {
	log.Printf("Client disconnected: %s", ctx.Conn.RemoteAddr())

	var errs []error
	for _, a := range h.actions {
		err := a.OnClose(ctx)
		if err != nil {
			if aerr := a.OnError(ctx); aerr != nil {
				errs = append(errs, aerr)
			}
		}
	}

	if len(errs) > 1 {
		return errors.Join(errs...)
	}

	return nil
}

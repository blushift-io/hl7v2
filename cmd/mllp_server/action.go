package main

import (
	"fmt"

	"github.com/blushift-io/hl7v2/net/tcp"
)

type action interface {
	tcp.Handler
}

func newAction(conf Action) (action, error) {
	switch conf.Type {
	case ActionNoop:
		return newNoopAction(), nil
	case ActionSave:
		return newSaveAction(conf.Options)
	case ActionForward:
		return newForwardAction(conf.Options)
	default:
		return nil, fmt.Errorf("unknown action type: %s", conf.Type)
	}
}

type noopAction struct{}

func newNoopAction() *noopAction {
	return &noopAction{}
}

func (a *noopAction) OnConnect(ctx *tcp.Context) error {
	return nil
}

func (a *noopAction) OnMessage(ctx *tcp.Context) error {
	return nil
}

func (a *noopAction) OnError(ctx *tcp.Context) error {
	return nil
}

func (a *noopAction) OnClose(ctx *tcp.Context) error {
	return nil
}

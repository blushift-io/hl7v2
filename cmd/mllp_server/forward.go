package main

import (
	"context"

	"github.com/blushift-io/hl7v2/net/tcp"
	"github.com/go-viper/mapstructure/v2"
)

type forwardAction struct {
	options *ForwardOptions
	c       tcp.Conn
}

func newForwardAction(options map[string]any) (*forwardAction, error) {
	var opts ForwardOptions
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &opts,
		TagName:          "yaml",
		WeaklyTypedInput: true,
	})
	if err != nil {
		return nil, err
	}

	if err := dec.Decode(options); err != nil {
		return nil, err
	}

	c, err := tcp.Dial(context.Background(), opts.Endpoint)
	if err != nil {
		return nil, err
	}

	return &forwardAction{options: &opts, c: c}, nil
}

func (a *forwardAction) OnConnect(ctx *tcp.Context) error {
	return nil
}

func (a *forwardAction) OnMessage(ctx *tcp.Context) error {
	return nil
}

func (a *forwardAction) OnError(ctx *tcp.Context) error {
	return nil
}

func (a *forwardAction) OnClose(ctx *tcp.Context) error {
	return nil
}

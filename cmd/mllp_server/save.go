package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/net/tcp"
	"github.com/go-viper/mapstructure/v2"
)

const defaultFilenameTmpl = "{{.Header.SendingFacility}}/{{.Header.ControlID}}.hl7"

type saveAction struct {
	options *SaveOptions
	ft      *template.Template
}

func newSaveAction(options map[string]any) (*saveAction, error) {
	var c SaveOptions
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &c,
		TagName:          "yaml",
		WeaklyTypedInput: true,
	})
	if err != nil {
		return nil, err
	}
	if err := dec.Decode(options); err != nil {
		return nil, err
	}

	return &saveAction{options: &c}, nil
}

func (a *saveAction) OnConnect(ctx *tcp.Context) error {
	return nil
}

func (a *saveAction) OnMessage(ctx *tcp.Context) error {
	msg, ok := ctx.Data.(*hl7v2.RawMessage)
	if !ok {
		return fmt.Errorf("on_message: expected RawMessage, got: %T: %w", ctx.Data, ErrInvalidMsg)
	}

	fn, err := a.getFilename(ctx, msg)
	if err != nil {
		return fmt.Errorf("error getting filename: %w", err)
	}

	out := filepath.Join(a.options.OutputPath, fn)
	if err := os.MkdirAll(out, 0755); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	if err := os.WriteFile(out, msg.Value().Bytes(), 0644); err != nil {
		return fmt.Errorf("error writing hl7 message to file: %w", err)
	}

	ack, err := ctx.Conn.AckMessage(msg)
	if err != nil {
		return fmt.Errorf("error sending ack: %w", err)
	}

	ab := ack.Value().Bytes()
	bd := filepath.Dir(out)
	afn := filepath.Join(bd, "ACK_"+filepath.Base(fn))
	if err := os.WriteFile(afn, ab, 0644); err != nil {
		return fmt.Errorf("error writing hl7 ack message to file: %w", err)
	}

	return nil
}

func (a *saveAction) OnError(ctx *tcp.Context) error {
	if errors.Is(ctx.Error, ErrInvalidMsg) {
		return ctx.Error
	}

	msg, ok := ctx.Data.(*hl7v2.RawMessage)
	if !ok {
		return fmt.Errorf("on_message: expected RawMessage, got: %T: %w", ctx.Data, ErrInvalidMsg)
	}

	if _, err := ctx.Conn.NackMessage(msg, hl7v2.AckApplicationError, ctx.Error.Error()); err != nil {
		return fmt.Errorf("error nacking message: %w, original error: %w", err, ctx.Error)
	}

	return ctx.Error
}

func (a *saveAction) OnClose(ctx *tcp.Context) error {
	return nil
}

type templateContext struct {
	Message *hl7v2.RawMessage
	Header  *hl7v2.MessageHeader
	Meta    map[string]string
}

func (a *saveAction) getFilename(ctx *tcp.Context, msg *hl7v2.RawMessage) (string, error) {
	hdr, err := hl7v2.ParseHeader(msg.Value().Bytes())
	if err != nil {
		return "", fmt.Errorf("error parsing hl7 header: %w", err)
	}

	tctx := templateContext{
		Message: msg,
		Header:  hdr,
		Meta: map[string]string{
			"remote_address": ctx.Conn.RemoteAddr().String(),
			"local_address":  ctx.Conn.LocalAddr().String(),
		},
	}

	buf := bytes.NewBuffer(nil)
	if err := a.ft.Execute(buf, tctx); err != nil {
		return "", fmt.Errorf("error executing filename template: %w", err)
	}

	return buf.String(), nil
}

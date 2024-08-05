package tcp

import "context"

type Context struct {
	context.Context
	Conn  Conn
	Error error
	Data  any
}

type ContextOption func(*Context)

func SetContext(ctx context.Context) ContextOption {
	return func(c *Context) {
		c.Context = ctx
	}
}

func SetData(data any) ContextOption {
	return func(c *Context) {
		c.Data = data
	}
}

func SetError(err error) ContextOption {
	return func(c *Context) {
		c.Error = err
	}
}

func NewContext(conn Conn, opts ...ContextOption) *Context {
	ctx := &Context{
		Context: context.Background(),
		Conn:    conn,
	}

	ctx.Apply(opts...)
	return ctx
}

func (c *Context) Apply(opts ...ContextOption) {
	for _, opt := range opts {
		opt(c)
	}
}

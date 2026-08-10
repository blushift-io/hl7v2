package tcp

import "context"

// Context carries connection state, context metadata, messages, and errors through handler functions.
type Context struct {
	context.Context
	Conn  Conn
	Error error
	Data  any
}

// ContextOption defines a function signature for configuring a Context.
type ContextOption func(*Context)

// SetContext returns a ContextOption that sets the underlying context.Context.
func SetContext(ctx context.Context) ContextOption {
	return func(c *Context) {
		c.Context = ctx
	}
}

// SetData returns a ContextOption that sets custom data on the Context.
func SetData(data any) ContextOption {
	return func(c *Context) {
		c.Data = data
	}
}

// SetError returns a ContextOption that sets an error on the Context.
func SetError(err error) ContextOption {
	return func(c *Context) {
		c.Error = err
	}
}

// NewContext creates a Context wrapping the given Conn and applying options.
func NewContext(conn Conn, opts ...ContextOption) *Context {
	ctx := &Context{
		Context: context.Background(),
		Conn:    conn,
	}

	ctx.Apply(opts...)
	return ctx
}

// Apply executes a list of ContextOption functions on the Context.
func (c *Context) Apply(opts ...ContextOption) {
	for _, opt := range opts {
		opt(c)
	}
}

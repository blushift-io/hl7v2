package generator

import "text/template"

// Options configures code generation parameters.
type Options struct {
	Templates map[SpecType][]*template.Template
}

// Option defines a functional option for configuring Options.
type Option func(*Options)

// NewOptions constructs an Options instance by applying option functions.
func NewOptions(opts ...Option) *Options {
	o := &Options{}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

// WithTemplate returns an Option that adds custom templates for a SpecType.
func WithTemplate(t SpecType, templ ...*template.Template) Option {
	return func(o *Options) {
		o.Templates[t] = append(o.Templates[t], templ...)
	}
}

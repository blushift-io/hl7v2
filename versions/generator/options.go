package generator

import "text/template"

type Options struct {
	Templates map[SpecType][]*template.Template
}

type Option func(*Options)

func NewOptions(opts ...Option) *Options {
	o := &Options{}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func WithTemplate(t SpecType, templ ...*template.Template) Option {
	return func(o *Options) {
		o.Templates[t] = append(o.Templates[t], templ...)
	}
}

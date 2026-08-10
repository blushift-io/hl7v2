package field

import (
	"fmt"

	"github.com/blushift-io/hl7v2"
)

type ConditionalElse struct {
	builder *Builder
	cond    bool
}

func When(cond bool, val any) *ConditionalElse {
	if cond {
		switch v := val.(type) {
		case *Builder:
			return &ConditionalElse{builder: v, cond: true}
		default:
			return &ConditionalElse{builder: String(fmt.Sprintf("%v", val)), cond: true}
		}
	}
	return &ConditionalElse{builder: NewBuilder(nil), cond: false}
}

func If(cond bool, val any) *ConditionalElse {
	return When(cond, val)
}

func (ce *ConditionalElse) Else(val any) *Builder {
	if ce.cond {
		return ce.builder
	}
	switch v := val.(type) {
	case *Builder:
		return v
	default:
		return String(fmt.Sprintf("%v", val))
	}
}

func (ce *ConditionalElse) Build() hl7v2.RawField {
	return ce.builder.Build()
}

func Eval(fn func() any) *Builder {
	res := fn()
	if b, ok := res.(*Builder); ok {
		return b
	}
	return String(fmt.Sprintf("%v", res))
}

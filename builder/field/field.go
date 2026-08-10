package field

import (
	"fmt"
	"time"

	"github.com/blushift-io/hl7v2"
)

type Builder struct {
	rawhl7 hl7v2.RawField
	err    error
}

func NewBuilder(rf hl7v2.RawField) *Builder {
	return &Builder{rawhl7: rf}
}

func String(v string) *Builder {
	return NewBuilder(hl7v2.RawField{
		hl7v2.RawRepetition{
			hl7v2.RawComponent{
				hl7v2.RawSubcomponent(v),
			},
		},
	})
}

func Int(v int) *Builder {
	return String(fmt.Sprintf("%d", v))
}

func Bool(v bool) *Builder {
	if v {
		return String("Y")
	}
	return String("N")
}

type TimeBuilder struct {
	*Builder
	t      time.Time
	layout string
}

func Time(t time.Time) *TimeBuilder {
	tb := &TimeBuilder{
		t:      t,
		layout: "20060102150405",
	}
	tb.Builder = String(t.Format(tb.layout))
	return tb
}

func (tb *TimeBuilder) Format(layout string) *TimeBuilder {
	tb.layout = layout
	tb.Builder = String(tb.t.Format(layout))
	return tb
}

func Components(comps ...any) *Builder {
	rawComps := make(hl7v2.RawRepetition, len(comps))
	for i, c := range comps {
		str := fmt.Sprintf("%v", c)
		rawComps[i] = hl7v2.RawComponent{hl7v2.RawSubcomponent(str)}
	}

	return NewBuilder(hl7v2.RawField{rawComps})
}

func Repeated(reps ...any) *Builder {
	var rawReps hl7v2.RawField
	for _, r := range reps {
		switch v := r.(type) {
		case interface{ Build() hl7v2.RawField }:
			raw := v.Build()
			if len(raw) > 0 {
				rawReps = append(rawReps, raw[0])
			}
		case hl7v2.RawField:
			if len(v) > 0 {
				rawReps = append(rawReps, v[0])
			}
		default:
			str := fmt.Sprintf("%v", r)
			rawReps = append(rawReps, hl7v2.RawRepetition{
				hl7v2.RawComponent{hl7v2.RawSubcomponent(str)},
			})
		}
	}

	return NewBuilder(rawReps)
}

func (b *Builder) Build() hl7v2.RawField {
	if b == nil {
		return hl7v2.RawField{}
	}
	return b.rawhl7
}

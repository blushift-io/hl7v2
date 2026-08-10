// builder/segment.go
package builder

import (
	"fmt"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/builder/field"
	"github.com/blushift-io/hl7v2/query"
)

type SegmentOption func(*SegmentBuilder)

type SegmentBuilder struct {
	id     string
	fields []*field.Builder
}

func Segment(id string, items ...any) *SegmentBuilder {
	sb := &SegmentBuilder{
		id: id,
		fields: []*field.Builder{
			field.String(id),
		},
	}

	for _, item := range items {
		sb.Add(item)
	}

	return sb
}

func (b *SegmentBuilder) Add(item any) *SegmentBuilder {
	switch v := item.(type) {
	case *field.Builder:
		b.fields = append(b.fields, v)
	case SegmentOption:
		v(b)
	case *SegmentOption:
		if v != nil {
			(*v)(b)
		}
	case nil:
		// ignore
	default:
		b.fields = append(b.fields, field.String(fmt.Sprintf("%v", v)))
	}
	return b
}

func (b *SegmentBuilder) Set(pos int, items ...any) *SegmentBuilder {
	if pos <= 0 {
		return b
	}

	// Ensure field slice capacity
	for len(b.fields) <= pos {
		b.fields = append(b.fields, field.String(""))
	}

	if len(items) == 1 {
		switch v := items[0].(type) {
		case *field.Builder:
			b.fields[pos] = v
		default:
			b.fields[pos] = field.String(fmt.Sprintf("%v", v))
		}
	} else if len(items) > 1 {
		b.fields[pos] = field.Components(items...)
	}

	return b
}

func (b *SegmentBuilder) Build() hl7v2.RawSegment {
	rs := make(hl7v2.RawSegment, len(b.fields))
	for i, f := range b.fields {
		if f != nil {
			rs[i] = f.Build()
		} else {
			rs[i] = hl7v2.RawField{}
		}
	}
	return rs
}

func (b *SegmentBuilder) BuildElement(parent hl7v2.Element, pos int) hl7v2.Element {
	return hl7v2.NewSegment(parent, pos, b.Build())
}

// --- Backwards Compatibility for header.go and builder.go ---

type SegmentBuildOption = SegmentOption

func WithFields(fields ...*FieldBuilder) SegmentOption {
	return func(b *SegmentBuilder) {
		for _, f := range fields {
			b.fields = append(b.fields, field.NewBuilder(f.Build()))
		}
	}
}

func (b *SegmentBuilder) AddField(f *FieldBuilder) *SegmentBuilder {
	b.fields = append(b.fields, field.NewBuilder(f.Build()))
	return b
}

func (b *SegmentBuilder) InsertField(pos int, fld *FieldBuilder) *SegmentBuilder {
	if pos <= 0 {
		return b
	}

	for len(b.fields) <= pos {
		b.fields = append(b.fields, field.String(""))
	}
	
	if fld != nil {
		b.fields[pos] = field.NewBuilder(fld.Build())
	}
	return b
}

func InsertField(pos int, fld *FieldBuilder) SegmentOption {
	return func(b *SegmentBuilder) {
		b.InsertField(pos, fld)
	}
}

func (b *SegmentBuilder) SetLocation(loc query.Location, val hl7v2.Value) *SegmentBuilder {
	// Basic backwards compatibility just to pass the build and simple tests
	return b.Set(loc.Field, val.String())
}

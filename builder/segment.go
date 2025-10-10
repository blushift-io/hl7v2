package builder

import (
	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/query"
)

type SegmentBuildOption func(*SegmentBuilder)

func WithFields(fields ...*FieldBuilder) SegmentBuildOption {
	return func(b *SegmentBuilder) {
		b.fields = append(b.fields, fields...)
	}
}

func InsertField(pos int, fld *FieldBuilder) SegmentBuildOption {
	return func(b *SegmentBuilder) {
		b.InsertField(pos, fld)
	}
}

type SegmentBuilder struct {
	id     string
	fields []*FieldBuilder
}

func Segment(id string, opts ...SegmentBuildOption) *SegmentBuilder {
	f := []*FieldBuilder{
		SingleValueField(hl7v2.NewStringValue(id)),
	}

	b := &SegmentBuilder{
		id:     id,
		fields: f,
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func (b *SegmentBuilder) Field() *FieldBuilder {
	f := &FieldBuilder{}
	b.fields = append(b.fields, f)

	return f
}

func (b *SegmentBuilder) AddField(f *FieldBuilder) *SegmentBuilder {
	b.fields = append(b.fields, f)

	return b
}

func (b *SegmentBuilder) InsertField(pos int, fld *FieldBuilder) *SegmentBuilder {
	if pos <= 0 {
		return b
	}

	if pos >= len(b.fields) {
		delta := pos + 1 - len(b.fields)
		for range delta {
			b.AddField(EmptyField())
		}
	}

	if fld != nil {
		b.fields[pos] = fld
	}

	return b
}

func (b *SegmentBuilder) SetLocation(loc query.Location, val hl7v2.Value) *SegmentBuilder {
	if loc.Field <= 0 {
		return b
	}

	size := len(b.fields) - 1
	if loc.Field > size {
		delta := loc.Field - size
		for range delta {
			b.AddField(EmptyField())
		}
	}

	b.fields[loc.Field].SetLocation(loc, val)

	return b
}

func (b *SegmentBuilder) Build() hl7v2.RawSegment {
	rs := make(hl7v2.RawSegment, len(b.fields))
	for i, f := range b.fields {
		rs[i] = f.Build()
	}

	return rs
}

func (b *SegmentBuilder) BuildElement(parent hl7v2.Element, pos int) hl7v2.Element {
	return hl7v2.NewSegment(parent, pos, b.Build())
}

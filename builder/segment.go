package builder

import (
	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/query"
)

// SegmentBuildOption defines a function signature for configuring a SegmentBuilder.
type SegmentBuildOption func(*SegmentBuilder)

// WithFields returns a SegmentBuildOption that appends field builders to a segment.
func WithFields(fields ...*FieldBuilder) SegmentBuildOption {
	return func(b *SegmentBuilder) {
		b.fields = append(b.fields, fields...)
	}
}

// InsertField returns a SegmentBuildOption that inserts a field builder at a specific position.
func InsertField(pos int, fld *FieldBuilder) SegmentBuildOption {
	return func(b *SegmentBuilder) {
		b.InsertField(pos, fld)
	}
}

// SegmentBuilder builds HL7 v2 message segments consisting of fields.
type SegmentBuilder struct {
	id     string
	fields []*FieldBuilder
}

// Segment creates a new SegmentBuilder with the given segment ID and options.
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

// Field appends a new empty FieldBuilder and returns it.
func (b *SegmentBuilder) Field() *FieldBuilder {
	f := &FieldBuilder{}
	b.fields = append(b.fields, f)

	return f
}

// AddField appends a FieldBuilder to the SegmentBuilder.
func (b *SegmentBuilder) AddField(f *FieldBuilder) *SegmentBuilder {
	b.fields = append(b.fields, f)

	return b
}

// InsertField places a FieldBuilder at the specified field position, padding with empty fields if necessary.
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

// SetLocation sets a value at the specified location within the segment.
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

// Build constructs and returns the RawSegment.
func (b *SegmentBuilder) Build() hl7v2.RawSegment {
	rs := make(hl7v2.RawSegment, len(b.fields))
	for i, f := range b.fields {
		rs[i] = f.Build()
	}

	return rs
}

// BuildElement constructs an HL7 v2 Element representation of the segment.
func (b *SegmentBuilder) BuildElement(parent hl7v2.Element, pos int) hl7v2.Element {
	return hl7v2.NewSegment(parent, pos, b.Build())
}

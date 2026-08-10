// Package builder provides a fluent API for constructing HL7 v2 messages, segments,
// fields, repetitions, components, and subcomponents.
package builder

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/query"
	"github.com/blushift-io/hl7v2/builder/field"
	"github.com/google/uuid"
)

// BuildOption defines a function signature for configuring a Builder.
type BuildOption func(*Builder)

// WithHeader configures a Builder with a new header using the specified message type, version, and header options.
func WithHeader(typ hl7v2.MessageType, ver hl7v2.Version, opts ...HeaderBuildOption) BuildOption {
	return func(b *Builder) {
		b.Header(typ, ver, opts...)
	}
}

// ApplyHeaderOptions configures header options on the Builder's header, creating a default header builder if one does not exist.
func ApplyHeaderOptions(opts ...HeaderBuildOption) BuildOption {
	return func(b *Builder) {
		if b.hdr == nil {
			b.hdr = &HeaderBuilder{hdr: &hl7v2.MessageHeader{}}
		}

		for _, opt := range opts {
			opt(b.hdr)
		}
	}
}

// Builder constructs HL7 v2 messages with customizable delimiters, headers, and segments.
type Builder struct {
	delims   *hl7v2.Delimiters
	hdr      *HeaderBuilder
	segments []*SegmentBuilder
}

// New creates a new Builder initialized with default delimiters and any provided build options.
func New(opts ...BuildOption) *Builder {
	b := &Builder{
		delims:   hl7v2.DefaultDelimiters(),
		segments: make([]*SegmentBuilder, 0),
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

// SetDelimiters sets custom delimiters for the Builder and its header builder if present.
func (b *Builder) SetDelimiters(delims *hl7v2.Delimiters) *Builder {
	b.delims = delims
	if b.hdr != nil {
		b.hdr.SetDelimiters(delims)
	}

	return b
}

// Header configures or updates the header for the Builder and returns its HeaderBuilder.
func (b *Builder) Header(typ hl7v2.MessageType, ver hl7v2.Version, opts ...HeaderBuildOption) *Builder {
	b.hdr = BuildHeader(typ, ver, opts...)
	return b
}

// SetHeader sets the HeaderBuilder for the Builder.
func (b *Builder) SetHeader(hdr *HeaderBuilder) *Builder {
	b.hdr = hdr

	return b
}

// AddSegment appends an existing SegmentBuilder to the Builder.
func (b *Builder) AddSegment(seg *SegmentBuilder) *Builder {
	b.segments = append(b.segments, seg)
	return b
}

// Segment creates and appends a new SegmentBuilder with the given ID and options.
func (b *Builder) Segment(id string, items ...any) *SegmentBuilder {
	seg := Segment(id, items...)
	b.segments = append(b.segments, seg)
	return seg
}

// Set parses a path like "PV1-1" or "PV1-3.1" and sets the value.
func (b *Builder) Set(path string, val any) *Builder {
	parts := strings.Split(path, "-")
	if len(parts) != 2 {
		return b
	}

	segID := parts[0]
	locParts := strings.Split(parts[1], ".")

	fieldIdx, err := strconv.Atoi(locParts[0])
	if err != nil || fieldIdx <= 0 {
		return b
	}

	seg := b.GetSegment(segID)
	if seg == nil {
		seg = Segment(segID)
		b.AddSegment(seg)
	}

	if len(locParts) == 1 {
		seg.Set(fieldIdx, val)
	} else if len(locParts) == 2 {
		compIdx, err := strconv.Atoi(locParts[1])
		if err == nil && compIdx > 0 {
			for len(seg.fields) <= fieldIdx {
				seg.fields = append(seg.fields, field.String(""))
			}
			fld := seg.fields[fieldIdx]
			var raw hl7v2.RawField
			if fld != nil {
				raw = fld.Build()
			}
			if len(raw) == 0 {
				raw = append(raw, hl7v2.RawRepetition{})
			}
			for len(raw[0]) < compIdx {
				raw[0] = append(raw[0], hl7v2.RawComponent{})
			}
			raw[0][compIdx-1] = hl7v2.RawComponent{hl7v2.RawSubcomponent(fmt.Sprintf("%v", val))}
			seg.fields[fieldIdx] = field.NewBuilder(raw)
		}
	}

	return b
}

// GetSegment retrieves the SegmentBuilder matching the given segment ID and repetition index.
func (b *Builder) GetSegment(id string, rep ...int) *SegmentBuilder {
	segRep := 0
	if len(rep) > 0 {
		segRep = rep[0]
	}

	idx := 0
	for _, seg := range b.segments {
		if seg.id == id {
			if idx == segRep {
				return seg
			}

			idx++
		}
	}

	return nil
}

// SetLocation assigns a value to a specific query location within the message segments.
func (b *Builder) SetLocation(loc query.Location, val hl7v2.Value) *Builder {
	rep := 0
	if loc.SegmentRep != nil {
		rep = *loc.SegmentRep
	}

	seg := b.GetSegment(loc.Segment, rep)
	if seg != nil {
		seg.SetLocation(loc, val)

		return b
	}

	if rep == 0 {
		return b.AddSegment(Segment(loc.Segment).SetLocation(loc, val))
	}

	for range rep {
		b.AddSegment(Segment(loc.Segment).SetLocation(loc, val))
	}

	return b.SetLocation(loc, val)
}

// Build constructs and returns the finalized HL7 v2 Message.
func (b *Builder) Build() (*hl7v2.Message, error) {
	raw := b.BuildRaw()
	return raw.ToMessage()
}

// BuildRaw constructs and returns the finalized RawMessage.
func (b *Builder) BuildRaw() *hl7v2.RawMessage {
	delims := b.delims
	if b.hdr != nil {
		delims = b.hdr.hdr.Delimiters
	}

	return hl7v2.NewRawMessage(delims, b.BuildRawSegements()...)
}

// BuildRawSegements constructs and returns a slice of RawSegment items for all header and body segments.
func (b *Builder) BuildRawSegements() []hl7v2.RawSegment {
	var segs []hl7v2.RawSegment

	if b.hdr != nil {
		segs = append(segs, b.hdr.Build())
	}

	for _, seg := range b.segments {
		segs = append(segs, seg.Build())
	}

	return segs
}

func generateID() string {
	s := uuid.NewString()

	parts := strings.Split(s, "-")
	return parts[len(parts)-1]
}

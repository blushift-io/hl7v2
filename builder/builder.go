package builder

import (
	"strings"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/query"
	"github.com/google/uuid"
)

type BuildOption func(*Builder)

func WithHeader(typ hl7v2.MessageType, ver hl7v2.Version, opts ...HeaderBuildOption) BuildOption {
	return func(b *Builder) {
		b.Header(typ, ver, opts...)
	}
}

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

type Builder struct {
	delims   *hl7v2.Delimiters
	hdr      *HeaderBuilder
	segments []*SegmentBuilder
}

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

func (b *Builder) SetDelimiters(delims *hl7v2.Delimiters) *Builder {
	b.delims = delims
	if b.hdr != nil {
		b.hdr.SetDelimiters(delims)
	}

	return b
}

func (b *Builder) Header(typ hl7v2.MessageType, ver hl7v2.Version, opts ...HeaderBuildOption) *HeaderBuilder {
	if b.hdr != nil {
		b.hdr.SetMessageType(typ).SetVersion(ver)
		return b.hdr
	}

	b.hdr = BuildHeader(typ, ver)

	return b.hdr
}

func (b *Builder) SetHeader(hdr *HeaderBuilder) *Builder {
	b.hdr = hdr

	return b
}

func (b *Builder) AddSegment(seg *SegmentBuilder) *Builder {
	b.segments = append(b.segments, seg)
	return b
}

func (b *Builder) Segment(id string, opts ...SegmentBuildOption) *SegmentBuilder {
	seg := Segment(id, opts...)
	b.segments = append(b.segments, seg)

	return seg
}

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

func (b *Builder) Build() (*hl7v2.Message, error) {
	raw := b.BuildRaw()
	return raw.ToMessage()
}

func (b *Builder) BuildRaw() *hl7v2.RawMessage {
	delims := b.delims
	if b.hdr != nil {
		delims = b.hdr.hdr.Delimiters
	}

	return hl7v2.NewRawMessage(delims, b.BuildRawSegements()...)
}

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

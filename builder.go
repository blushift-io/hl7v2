package hl7v2

import (
	"strings"
	"time"

	"github.com/blushift-io/hl7v2/query"
	"github.com/google/uuid"
)

type Builder struct {
	delims   *Delimiters
	hdr      *HeaderBuilder
	segments []*SegmentBuilder
}

func NewBuilder() *Builder {
	return &Builder{
		delims:   DefaultDelimiters(),
		segments: make([]*SegmentBuilder, 0),
	}
}

func (b *Builder) SetDelimiters(delims *Delimiters) *Builder {
	b.delims = delims
	if b.hdr != nil {
		b.hdr.SetDelimiters(delims)
	}

	return b
}

func (b *Builder) Header(typ MessageType, ver string) *HeaderBuilder {
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

func (b *Builder) Segment(id string, flds ...*FieldBuilder) *SegmentBuilder {
	seg := BuildSegment(id, flds...)
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

func (b *Builder) SetLocation(loc query.Location, val Value) *Builder {
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
		return b.AddSegment(BuildSegment(loc.Segment).SetLocation(loc, val))
	}

	for range rep {
		b.AddSegment(BuildSegment(loc.Segment).SetLocation(loc, val))
	}

	return b.SetLocation(loc, val)
}

func (b *Builder) Build() (*Message, error) {
	return newMessage(nil, 0, b.BuildRaw())
}

func (b *Builder) BuildRaw() *RawMessage {
	delims := b.delims
	if b.hdr != nil {
		delims = b.hdr.hdr.Delimiters
	}

	return newRawMessage(delims, b.BuildRawSegements()...)
}

func (b *Builder) BuildRawSegements() []RawSegment {
	var segs []RawSegment

	if b.hdr != nil {
		segs = append(segs, b.hdr.Build())
	}

	for _, seg := range b.segments {
		segs = append(segs, seg.Build())
	}

	return segs
}

type HeaderBuilder struct {
	hdr *MessageHeader
}

func BuildHeader(typ MessageType, ver string) *HeaderBuilder {
	hdr := &MessageHeader{
		Delimiters:       DefaultDelimiters(),
		MessageDate:      time.Now(),
		MessageCode:      typ.Code,
		TriggerEvent:     typ.Event,
		MessageStructure: typ.Structure,
		ControlID:        generateID(),
		VersionID:        ver,
	}

	return &HeaderBuilder{hdr: hdr}
}

func NewHeaderBuilder(hdr *MessageHeader) *HeaderBuilder {
	return &HeaderBuilder{hdr: hdr}
}

func (b *HeaderBuilder) Set(hdr *MessageHeader) *HeaderBuilder {
	b.hdr = hdr

	return b
}

func (b *HeaderBuilder) SetDelimiters(delims *Delimiters) *HeaderBuilder {
	b.hdr.Delimiters = delims
	return b
}

func (b *HeaderBuilder) SetSendingApplication(app string) *HeaderBuilder {
	b.hdr.SendingApplication = app
	return b
}

func (b *HeaderBuilder) SetSendingFacility(facility string) *HeaderBuilder {
	b.hdr.SendingFacility = facility
	return b
}

func (b *HeaderBuilder) SetReceivingApplication(app string) *HeaderBuilder {
	b.hdr.ReceivingApplication = app
	return b
}

func (b *HeaderBuilder) SetReceivingFacility(facility string) *HeaderBuilder {
	b.hdr.ReceivingFacility = facility
	return b
}

func (b *HeaderBuilder) SetMessageDate(date time.Time) *HeaderBuilder {
	b.hdr.MessageDate = date
	return b
}

func (b *HeaderBuilder) SetMessageType(typ MessageType) *HeaderBuilder {
	b.hdr.MessageCode = typ.Code
	b.hdr.TriggerEvent = typ.Event
	b.hdr.MessageStructure = typ.Structure

	return b
}

func (b *HeaderBuilder) SetVersion(ver string) *HeaderBuilder {
	b.hdr.VersionID = ver
	return b
}

func (b *HeaderBuilder) SetSecurity(security string) *HeaderBuilder {
	b.hdr.Security = security
	return b
}

func (b *HeaderBuilder) SetControlID(id string) *HeaderBuilder {
	b.hdr.ControlID = id
	return b
}

func (b *HeaderBuilder) SetProcessingID(id string) *HeaderBuilder {
	b.hdr.ProcessingID = id
	return b
}

func (b *HeaderBuilder) SetSequenceNumber(seq int) *HeaderBuilder {
	b.hdr.SequenceNumber = seq
	return b
}

func (b *HeaderBuilder) SetContinuationPointer(pointer string) *HeaderBuilder {
	b.hdr.ContinuationPointer = pointer
	return b
}

func (b *HeaderBuilder) Build() RawSegment {
	return b.Segment().Build()
}

func (b *HeaderBuilder) BuildElement(parent Element) Element {
	return newSegment(parent, 0, b.Build())
}

func (b *HeaderBuilder) Segment() *SegmentBuilder {
	msgTypeFld := ComponentField(
		NewStringValue(b.hdr.MessageCode),
		NewStringValue(b.hdr.TriggerEvent),
		NewStringValue(b.hdr.MessageStructure),
	)

	return BuildSegment("MSH").
		AddField(SingleValueField(b.hdr.Delimiters.fieldValue())).
		AddField(SingleValueField(b.hdr.Delimiters.encodingCharsValue())).
		AddField(SingleValueField(NewStringValue(b.hdr.SendingApplication))).
		AddField(SingleValueField(NewStringValue(b.hdr.SendingFacility))).
		AddField(SingleValueField(NewStringValue(b.hdr.ReceivingApplication))).
		AddField(SingleValueField(NewStringValue(b.hdr.ReceivingFacility))).
		AddField(SingleValueField(NewStringValue(b.hdr.MessageDate.Format(time.RFC3339)))).
		AddField(SingleValueField(NewStringValue(b.hdr.Security))).
		AddField(msgTypeFld).
		AddField(SingleValueField(NewStringValue(b.hdr.ControlID))).
		AddField(SingleValueField(NewStringValue(b.hdr.ProcessingID))).
		AddField(SingleValueField(NewStringValue(b.hdr.VersionID))).
		AddField(SingleValueField(NewIntValue(b.hdr.SequenceNumber))).
		AddField(SingleValueField(NewStringValue(b.hdr.ContinuationPointer))).
		AddField(SingleValueField(NewStringValue(b.hdr.AcceptAcknowledgmentType))).
		AddField(SingleValueField(NewStringValue(b.hdr.ApplicationAcknowledgmentType))).
		AddField(SingleValueField(NewStringValue(b.hdr.CountryCode))).
		AddField(SingleValueField(NewStringValue(b.hdr.CharacterSet))).
		AddField(SingleValueField(NewStringValue(b.hdr.PrincipalLanguage))).
		AddField(SingleValueField(NewStringValue(b.hdr.AlternateCharacterSet))).
		AddField(SingleValueField(NewStringValue(b.hdr.ProfileID))).
		AddField(SingleValueField(NewStringValue(b.hdr.SendingResponsibleOrgCode))).
		AddField(SingleValueField(NewStringValue(b.hdr.ReceivingResponsibleOrgCode))).
		AddField(SingleValueField(NewStringValue(b.hdr.SendingNetworkAddress))).
		AddField(SingleValueField(NewStringValue(b.hdr.ReceivingNetworkAddress)))
}

type SegmentBuilder struct {
	id     string
	fields []*FieldBuilder
}

func BuildSegment(id string, fields ...*FieldBuilder) *SegmentBuilder {
	f := []*FieldBuilder{
		SingleValueField(NewStringValue(id)),
	}

	f = append(f, fields...)

	return &SegmentBuilder{
		id:     id,
		fields: f,
	}
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

func (b *SegmentBuilder) SetLocation(loc query.Location, val Value) *SegmentBuilder {
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

func (b *SegmentBuilder) Build() RawSegment {
	rs := make(RawSegment, len(b.fields))
	for i, f := range b.fields {
		rs[i] = f.Build()
	}

	return rs
}

func (b *SegmentBuilder) BuildElement(parent Element, pos int) Element {
	return newSegment(parent, pos, b.Build())
}

type FieldBuilder struct {
	reps []*RepetitionBuilder
}

func BuildField(reps ...*RepetitionBuilder) *FieldBuilder {
	return &FieldBuilder{
		reps: reps,
	}
}

func SingleValueField(val Value) *FieldBuilder {
	return &FieldBuilder{
		reps: []*RepetitionBuilder{
			BuildRepetition(BuildComponent(BuildSubcomponent(val))),
		},
	}
}

func RepeatingField(vals ...Value) *FieldBuilder {
	reps := make([]*RepetitionBuilder, len(vals))

	for i, v := range vals {
		reps[i] = BuildRepetition(
			BuildComponent(
				BuildSubcomponent(v),
			),
		)

	}

	return BuildField(reps...)
}

func ComponentField(vals ...Value) *FieldBuilder {
	comps := make([]*ComponentBuilder, len(vals))

	for i, v := range vals {
		comps[i] = BuildComponent(
			BuildSubcomponent(v),
		)
	}

	return BuildField(BuildRepetition(comps...))
}

func EmptyField() *FieldBuilder {
	return &FieldBuilder{}
}

func (f *FieldBuilder) Repetition() *RepetitionBuilder {
	r := &RepetitionBuilder{}
	f.reps = append(f.reps, r)

	return r
}

func (f *FieldBuilder) AddRepetition(r *RepetitionBuilder) *FieldBuilder {
	f.reps = append(f.reps, r)

	return f
}

func (f *FieldBuilder) SetLocation(loc query.Location, val Value) *FieldBuilder {
	rep := 0
	if loc.FieldRep != nil {
		rep = *loc.FieldRep
	}

	size := len(f.reps)
	if size == 0 {
		if delta := rep - size; delta == 0 {
			return f.AddRepetition(EmptyRepetition()).SetLocation(loc, val)
		} else {
			for range delta {
				f.AddRepetition(EmptyRepetition())
			}
		}
	} else if rep > size {
		for range rep - size {
			f.AddRepetition(EmptyRepetition())
		}
	}

	f.reps[rep].SetLocation(loc, val)

	return f
}

func (f *FieldBuilder) Build() RawField {
	rf := make(RawField, len(f.reps))
	for i, r := range f.reps {
		rf[i] = r.Build()
	}

	return rf
}

func (f *FieldBuilder) BuildElement(parent Element, pos int) Element {
	return newField(parent, pos, f.Build())
}

type RepetitionBuilder struct {
	comps []*ComponentBuilder
}

func BuildRepetition(comps ...*ComponentBuilder) *RepetitionBuilder {
	return &RepetitionBuilder{
		comps: comps,
	}
}

func EmptyRepetition() *RepetitionBuilder {
	return &RepetitionBuilder{}
}

func (r *RepetitionBuilder) Component() *ComponentBuilder {
	c := &ComponentBuilder{}
	r.comps = append(r.comps, c)

	return c
}

func (r *RepetitionBuilder) AddComponent(c *ComponentBuilder) *RepetitionBuilder {
	r.comps = append(r.comps, c)

	return r
}

func (r *RepetitionBuilder) SetLocation(loc query.Location, val Value) *RepetitionBuilder {
	if loc.Component < 0 {
		return r
	}

	size := len(r.comps)
	if size == 0 {
		if delta := loc.Component - size; delta == 0 {
			return r.AddComponent(EmptyComponent()).SetLocation(loc, val)
		} else {
			for range delta {
				r.AddComponent(EmptyComponent())
			}
		}
	} else if loc.Component > size {
		delta := loc.Component - size
		for range delta {
			r.AddComponent(EmptyComponent())
		}
	}
	r.comps[loc.Component].SetLocation(loc, val)

	return r
}

func (r *RepetitionBuilder) Build() RawRepetition {
	rr := make(RawRepetition, len(r.comps))
	for i, c := range r.comps {
		rr[i] = c.Build()
	}
	return rr
}

func (r *RepetitionBuilder) BuildElement(parent Element, pos int) Element {
	return newRepetition(parent, pos, r.Build())
}

type ComponentBuilder struct {
	subs []*SubcomponentBuilder
}

func BuildComponent(subs ...*SubcomponentBuilder) *ComponentBuilder {
	return &ComponentBuilder{
		subs: subs,
	}
}

func SingleComponent(val Value) *ComponentBuilder {
	return &ComponentBuilder{
		subs: []*SubcomponentBuilder{
			BuildSubcomponent(val),
		},
	}
}

func EmptyComponent() *ComponentBuilder {
	return &ComponentBuilder{}
}

func (c *ComponentBuilder) Subcomponent() *ComponentBuilder {
	s := &SubcomponentBuilder{}
	c.subs = append(c.subs, s)

	return c
}

func (c *ComponentBuilder) AddSubcomponent(s *SubcomponentBuilder) *ComponentBuilder {
	c.subs = append(c.subs, s)

	return c
}

func (c *ComponentBuilder) SetLocation(loc query.Location, v Value) *ComponentBuilder {
	if loc.Subcomponent < 0 {
		return c
	}

	size := len(c.subs)
	if size == 0 {
		if delta := loc.Subcomponent - size; delta == 0 {
			return c.AddSubcomponent(BuildSubcomponent(v))
		} else {
			for range delta {
				c.AddSubcomponent(EmptySubcomponent())
			}
		}
	} else if loc.Component > size {
		for range loc.Component - size {
			c.AddSubcomponent(EmptySubcomponent())
		}
	}

	c.subs[loc.Subcomponent].SetValue(v)

	return c
}

func (c *ComponentBuilder) Build() RawComponent {
	rc := make(RawComponent, len(c.subs))
	for i, s := range c.subs {
		rc[i] = s.Build()
	}

	return rc
}

func (c *ComponentBuilder) BuildElement(parent Element, pos int) Element {
	return newComponent(parent, pos, c.Build())
}

type SubcomponentBuilder struct {
	val []byte
}

func BuildSubcomponent(val Value) *SubcomponentBuilder {
	return &SubcomponentBuilder{
		val: val.Bytes(),
	}
}

func EmptySubcomponent() *SubcomponentBuilder {
	return &SubcomponentBuilder{
		val: []byte{},
	}
}

func (s *SubcomponentBuilder) SetValue(v Value) *SubcomponentBuilder {
	s.val = v.Bytes()

	return s
}

func (s *SubcomponentBuilder) Build() RawSubcomponent {
	return RawSubcomponent(s.val)
}

func (s *SubcomponentBuilder) BuildElement(parent Element, pos int) Element {
	return newSubcomponent(parent, pos, NewValue(s.val))
}

func generateID() string {
	s := uuid.NewString()

	parts := strings.Split(s, "-")
	return parts[len(parts)-1]
}

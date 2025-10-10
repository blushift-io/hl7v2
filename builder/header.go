package builder

import (
	"time"

	"github.com/blushift-io/hl7v2"
)

type HeaderBuildOption func(*HeaderBuilder)

type HeaderBuilder struct {
	includeMsgStructure bool
	hdr                 *hl7v2.MessageHeader
}

func BuildHeader(typ hl7v2.MessageType, ver hl7v2.Version, opts ...HeaderBuildOption) *HeaderBuilder {
	hdr := &hl7v2.MessageHeader{
		Delimiters:       hl7v2.DefaultDelimiters(),
		MessageDate:      time.Now(),
		MessageCode:      typ.Code,
		TriggerEvent:     typ.Event,
		MessageStructure: typ.Structure,
		ControlID:        generateID(),
		VersionID:        ver.String(),
	}

	b := &HeaderBuilder{includeMsgStructure: true, hdr: hdr}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func NewHeaderBuilder(hdr *hl7v2.MessageHeader) *HeaderBuilder {
	return &HeaderBuilder{hdr: hdr}
}

func (b *HeaderBuilder) Set(hdr *hl7v2.MessageHeader) *HeaderBuilder {
	b.hdr = hdr

	return b
}

func IncludeMsgStructure(b bool) HeaderBuildOption {
	return func(bld *HeaderBuilder) {
		bld.includeMsgStructure = b
	}
}

func (b *HeaderBuilder) IncludeMsgStructure(bInclude bool) *HeaderBuilder {
	b.includeMsgStructure = bInclude
	return b
}

func SetDelimiters(delims *hl7v2.Delimiters) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.Delimiters = delims
	}
}

func (b *HeaderBuilder) SetDelimiters(delims *hl7v2.Delimiters) *HeaderBuilder {
	b.hdr.Delimiters = delims
	return b
}

func SetSendingApplication(app string) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.SendingApplication = app
	}
}

func (b *HeaderBuilder) SetSendingApplication(app string) *HeaderBuilder {
	b.hdr.SendingApplication = app
	return b
}

func SetSendingFacility(facility string) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.SendingFacility = facility
	}
}

func (b *HeaderBuilder) SetSendingFacility(facility string) *HeaderBuilder {
	b.hdr.SendingFacility = facility
	return b
}

func SetReceivingApplication(app string) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.ReceivingApplication = app
	}
}

func (b *HeaderBuilder) SetReceivingApplication(app string) *HeaderBuilder {
	b.hdr.ReceivingApplication = app
	return b
}

func SetReceivingFacility(facility string) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.ReceivingFacility = facility
	}
}

func (b *HeaderBuilder) SetReceivingFacility(facility string) *HeaderBuilder {
	b.hdr.ReceivingFacility = facility
	return b
}

func SetMessageDate(date time.Time) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.MessageDate = date
	}
}

func (b *HeaderBuilder) SetMessageDate(date time.Time) *HeaderBuilder {
	b.hdr.MessageDate = date
	return b
}

func SetMessageType(typ hl7v2.MessageType) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.MessageCode = typ.Code
		b.hdr.TriggerEvent = typ.Event
		b.hdr.MessageStructure = typ.Structure
	}
}

func (b *HeaderBuilder) SetMessageType(typ hl7v2.MessageType) *HeaderBuilder {
	b.hdr.MessageCode = typ.Code
	b.hdr.TriggerEvent = typ.Event
	b.hdr.MessageStructure = typ.Structure

	return b
}

func SetVersion(ver hl7v2.Version) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.VersionID = ver.String()
	}
}

func (b *HeaderBuilder) SetVersion(ver hl7v2.Version) *HeaderBuilder {
	b.hdr.VersionID = ver.String()
	return b
}

func SetSecurity(security string) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.Security = security
	}
}

func (b *HeaderBuilder) SetSecurity(security string) *HeaderBuilder {
	b.hdr.Security = security
	return b
}

func SetControlID(id string) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.ControlID = id
	}
}

func (b *HeaderBuilder) SetControlID(id string) *HeaderBuilder {
	b.hdr.ControlID = id
	return b
}

func SetProcessingID(id string) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.ProcessingID = id
	}
}

func (b *HeaderBuilder) SetProcessingID(id string) *HeaderBuilder {
	b.hdr.ProcessingID = id
	return b
}

func SetSequenceNumber(seq int) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.SequenceNumber = seq
	}
}

func (b *HeaderBuilder) SetSequenceNumber(seq int) *HeaderBuilder {
	b.hdr.SequenceNumber = seq
	return b
}

func SetContinuationPointer(pointer string) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.ContinuationPointer = pointer
	}
}

func (b *HeaderBuilder) SetContinuationPointer(pointer string) *HeaderBuilder {
	b.hdr.ContinuationPointer = pointer
	return b
}

func (b *HeaderBuilder) Build() hl7v2.RawSegment {
	return b.Segment().Build()
}

func (b *HeaderBuilder) BuildElement(parent hl7v2.Element) hl7v2.Element {
	return hl7v2.NewSegment(parent, 0, b.Build())
}

func (b *HeaderBuilder) Segment(opts ...SegmentBuildOption) *SegmentBuilder {
	msh9vals := []hl7v2.Value{
		hl7v2.NewStringValue(b.hdr.MessageCode),
		hl7v2.NewStringValue(b.hdr.TriggerEvent),
	}

	if b.includeMsgStructure && b.hdr.MessageStructure != "" {
		msh9vals = append(msh9vals, hl7v2.NewStringValue(b.hdr.MessageStructure))
	}

	opts = append([]SegmentBuildOption{
		WithFields(
			SingleValueField(b.hdr.Delimiters.FieldSeparatorValue()),
			SingleValueField(b.hdr.Delimiters.EncodingCharsValue()),
			SingleValueField(hl7v2.NewStringValue(b.hdr.SendingApplication)),
			SingleValueField(hl7v2.NewStringValue(b.hdr.SendingFacility)),
			SingleValueField(hl7v2.NewStringValue(b.hdr.ReceivingApplication)),
			SingleValueField(hl7v2.NewStringValue(b.hdr.ReceivingFacility)),
			SingleValueField(hl7v2.NewStringValue(b.hdr.MessageDate.Format("20060102150405"))),
			SingleValueField(hl7v2.NewStringValue(b.hdr.Security)),
			ComponentField(msh9vals...),
			SingleValueField(hl7v2.NewStringValue(b.hdr.ControlID)),
			SingleValueField(hl7v2.NewStringValue(b.hdr.ProcessingID)),
			SingleValueField(hl7v2.NewStringValue(b.hdr.VersionID)),
			SingleValueField(hl7v2.NewIntValue(b.hdr.SequenceNumber)),
			SingleValueField(hl7v2.NewStringValue(b.hdr.ContinuationPointer)),
			SingleValueField(hl7v2.NewStringValue(b.hdr.AcceptAcknowledgmentType)),
			SingleValueField(hl7v2.NewStringValue(b.hdr.ApplicationAcknowledgmentType)),
			SingleValueField(hl7v2.NewStringValue(b.hdr.CountryCode)),
			SingleValueField(hl7v2.NewStringValue(b.hdr.CharacterSet)),
			SingleValueField(hl7v2.NewStringValue(b.hdr.PrincipalLanguage)),
			SingleValueField(hl7v2.NewStringValue(b.hdr.AlternateCharacterSet)),
			SingleValueField(hl7v2.NewStringValue(b.hdr.ProfileID)),
			SingleValueField(hl7v2.NewStringValue(b.hdr.SendingResponsibleOrgCode)),
			SingleValueField(hl7v2.NewStringValue(b.hdr.ReceivingResponsibleOrgCode)),
			SingleValueField(hl7v2.NewStringValue(b.hdr.SendingNetworkAddress)),
			SingleValueField(hl7v2.NewStringValue(b.hdr.ReceivingNetworkAddress)),
		),
	}, opts...)

	return Segment("MSH", opts...)
}

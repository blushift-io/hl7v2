package builder

import (
	"time"

	"github.com/blushift-io/hl7v2"
)

// HeaderBuildOption defines a function signature for configuring a HeaderBuilder.
type HeaderBuildOption func(*HeaderBuilder)

// HeaderBuilder builds HL7 v2 MSH header segments.
type HeaderBuilder struct {
	includeMsgStructure bool
	hdr                 *hl7v2.MessageHeader
}

// BuildHeader creates a HeaderBuilder with default header settings for the given message type and version.
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

// NewHeaderBuilder creates a HeaderBuilder wrapping an existing MessageHeader.
func NewHeaderBuilder(hdr *hl7v2.MessageHeader) *HeaderBuilder {
	return &HeaderBuilder{hdr: hdr}
}

// Set updates the underlying MessageHeader of the HeaderBuilder.
func (b *HeaderBuilder) Set(hdr *hl7v2.MessageHeader) *HeaderBuilder {
	b.hdr = hdr

	return b
}

// IncludeMsgStructure returns a HeaderBuildOption that configures whether to include the message structure in MSH-9.
func IncludeMsgStructure(b bool) HeaderBuildOption {
	return func(bld *HeaderBuilder) {
		bld.includeMsgStructure = b
	}
}

// IncludeMsgStructure configures whether to include the message structure in MSH-9.
func (b *HeaderBuilder) IncludeMsgStructure(bInclude bool) *HeaderBuilder {
	b.includeMsgStructure = bInclude
	return b
}

// SetDelimiters returns a HeaderBuildOption that sets message delimiters on the header.
func SetDelimiters(delims *hl7v2.Delimiters) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.Delimiters = delims
	}
}

// SetDelimiters sets message delimiters on the header.
func (b *HeaderBuilder) SetDelimiters(delims *hl7v2.Delimiters) *HeaderBuilder {
	b.hdr.Delimiters = delims
	return b
}

// SetSendingApplication returns a HeaderBuildOption that sets the sending application in MSH-3.
func SetSendingApplication(app string) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.SendingApplication = app
	}
}

// SetSendingApplication sets the sending application in MSH-3.
func (b *HeaderBuilder) SetSendingApplication(app string) *HeaderBuilder {
	b.hdr.SendingApplication = app
	return b
}

// SetSendingFacility returns a HeaderBuildOption that sets the sending facility in MSH-4.
func SetSendingFacility(facility string) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.SendingFacility = facility
	}
}

// SetSendingFacility sets the sending facility in MSH-4.
func (b *HeaderBuilder) SetSendingFacility(facility string) *HeaderBuilder {
	b.hdr.SendingFacility = facility
	return b
}

// SetReceivingApplication returns a HeaderBuildOption that sets the receiving application in MSH-5.
func SetReceivingApplication(app string) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.ReceivingApplication = app
	}
}

// SetReceivingApplication sets the receiving application in MSH-5.
func (b *HeaderBuilder) SetReceivingApplication(app string) *HeaderBuilder {
	b.hdr.ReceivingApplication = app
	return b
}

// SetReceivingFacility returns a HeaderBuildOption that sets the receiving facility in MSH-6.
func SetReceivingFacility(facility string) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.ReceivingFacility = facility
	}
}

// SetReceivingFacility sets the receiving facility in MSH-6.
func (b *HeaderBuilder) SetReceivingFacility(facility string) *HeaderBuilder {
	b.hdr.ReceivingFacility = facility
	return b
}

// SetMessageDate returns a HeaderBuildOption that sets the message timestamp in MSH-7.
func SetMessageDate(date time.Time) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.MessageDate = date
	}
}

// SetMessageDate sets the message timestamp in MSH-7.
func (b *HeaderBuilder) SetMessageDate(date time.Time) *HeaderBuilder {
	b.hdr.MessageDate = date
	return b
}

// SetMessageType returns a HeaderBuildOption that sets the message type in MSH-9.
func SetMessageType(typ hl7v2.MessageType) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.MessageCode = typ.Code
		b.hdr.TriggerEvent = typ.Event
		b.hdr.MessageStructure = typ.Structure
	}
}

// SetMessageType sets the message type in MSH-9.
func (b *HeaderBuilder) SetMessageType(typ hl7v2.MessageType) *HeaderBuilder {
	b.hdr.MessageCode = typ.Code
	b.hdr.TriggerEvent = typ.Event
	b.hdr.MessageStructure = typ.Structure

	return b
}

// SetVersion returns a HeaderBuildOption that sets the HL7 v2 version in MSH-12.
func SetVersion(ver hl7v2.Version) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.VersionID = ver.String()
	}
}

// SetVersion sets the HL7 v2 version in MSH-12.
func (b *HeaderBuilder) SetVersion(ver hl7v2.Version) *HeaderBuilder {
	b.hdr.VersionID = ver.String()
	return b
}

// SetSecurity returns a HeaderBuildOption that sets the security field in MSH-8.
func SetSecurity(security string) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.Security = security
	}
}

// SetSecurity sets the security field in MSH-8.
func (b *HeaderBuilder) SetSecurity(security string) *HeaderBuilder {
	b.hdr.Security = security
	return b
}

// SetControlID returns a HeaderBuildOption that sets the message control ID in MSH-10.
func SetControlID(id string) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.ControlID = id
	}
}

// SetControlID sets the message control ID in MSH-10.
func (b *HeaderBuilder) SetControlID(id string) *HeaderBuilder {
	b.hdr.ControlID = id
	return b
}

// SetProcessingID returns a HeaderBuildOption that sets the processing ID in MSH-11.
func SetProcessingID(id string) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.ProcessingID = id
	}
}

// SetProcessingID sets the processing ID in MSH-11.
func (b *HeaderBuilder) SetProcessingID(id string) *HeaderBuilder {
	b.hdr.ProcessingID = id
	return b
}

// SetSequenceNumber returns a HeaderBuildOption that sets the sequence number in MSH-13.
func SetSequenceNumber(seq int) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.SequenceNumber = seq
	}
}

// SetSequenceNumber sets the sequence number in MSH-13.
func (b *HeaderBuilder) SetSequenceNumber(seq int) *HeaderBuilder {
	b.hdr.SequenceNumber = seq
	return b
}

// SetContinuationPointer returns a HeaderBuildOption that sets the continuation pointer in MSH-14.
func SetContinuationPointer(pointer string) HeaderBuildOption {
	return func(b *HeaderBuilder) {
		b.hdr.ContinuationPointer = pointer
	}
}

// SetContinuationPointer sets the continuation pointer in MSH-14.
func (b *HeaderBuilder) SetContinuationPointer(pointer string) *HeaderBuilder {
	b.hdr.ContinuationPointer = pointer
	return b
}

// Build constructs and returns the RawSegment for the MSH header.
func (b *HeaderBuilder) Build() hl7v2.RawSegment {
	return b.Segment().Build()
}

// BuildElement constructs an HL7 v2 Element representation of the MSH segment.
func (b *HeaderBuilder) BuildElement(parent hl7v2.Element) hl7v2.Element {
	return hl7v2.NewSegment(parent, 0, b.Build())
}

// Segment generates a SegmentBuilder representing the MSH segment.
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

	anyOpts := make([]any, len(opts))
	for i, o := range opts {
		anyOpts[i] = o
	}

	return Segment("MSH", anyOpts...)
}

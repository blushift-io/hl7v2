package hl7v2

import (
	"fmt"
	"time"
)

//TODO: Clean all this mess up

//go:generate enumer -type=HeaderType,TrailerType --linecomment -output=header_enums.go

// HeaderType is an enum type for valid header types
type HeaderType int

const (
	// HeaderTypeMessage represents an MSH header type.
	HeaderTypeMessage HeaderType = iota //MSH
	// HeaderTypeFile represents an FHS header type.
	HeaderTypeFile                      //FHS
	// HeaderTypeBatch represents a BHS header type.
	HeaderTypeBatch                     //BHS
)

// Trailer returns the corresponding TrailerType for the HeaderType.
func (t HeaderType) Trailer() TrailerType {
	switch t {
	case HeaderTypeFile:
		return TrailerTypeFile
	case HeaderTypeBatch:
		return TrailerTypeBatch
	default:
		return TrailerTypeNone
	}
}

// TrailerType is an enum for valid trailer segment types.
type TrailerType int

const (
	// TrailerTypeNone indicates no trailer segment.
	TrailerTypeNone  TrailerType = iota //_
	// TrailerTypeFile represents an FTS trailer segment type.
	TrailerTypeFile                     //FTS
	// TrailerTypeBatch represents a BTS trailer segment type.
	TrailerTypeBatch                    //BTS
)

// Header is an interface representing an HL7 message, file, or batch header.
type Header interface {
	Type() HeaderType
	Delimiters() *Delimiters
	SendingApplication() string
	SendingFacility() string
	ReceivingApplication() string
	ReceivingFacility() string
	DateTime() time.Time
	Security() string
}

// NewHeader returns a new Header for the given HeaderType.
func NewHeader(typ HeaderType) Header {
	return &header{
		typ:    typ,
		delims: DefaultDelimiters(),
	}
}

type header struct {
	typ          HeaderType
	delims       *Delimiters
	sendingApp   string
	receivingApp string
	sendingFac   string
	receivingFac string
	dateTime     time.Time
	security     string
}

func (h *header) Type() HeaderType {
	return h.typ
}

func (h *header) Delimiters() *Delimiters {
	return h.delims
}

func (h *header) SendingApplication() string {
	return h.sendingApp
}

func (h *header) SendingFacility() string {
	return h.sendingFac
}

func (h *header) ReceivingApplication() string {
	return h.receivingApp
}

func (h *header) ReceivingFacility() string {
	return h.receivingFac
}

func (h *header) DateTime() time.Time {
	return h.dateTime
}

func (h *header) Security() string {
	return h.security
}

// MessageType represents the HL7 message type, trigger event, and structure.
type MessageType struct {
	Code      string `json:"code"`
	Event     string `json:"event"`
	Structure string `json:"structure"`
}

// ParseMessageType parses a message type string into a MessageType struct.
func ParseMessageType(s string, delims *Delimiters) (MessageType, error) {
	parts := delims.Split([]byte(s), ComponentDelimiter)
	var code string
	var evt string
	var str string

	switch len(parts) {
	case 1:
		if string(parts[0]) != "ACK" {
			return MessageType{}, fmt.Errorf("invalid message type: %s", s)
		}

		code = string(parts[0])
	case 2:
		code = string(parts[0])
		evt = string(parts[1])
		str = code + "_" + evt
	case 3:
		code = string(parts[0])
		evt = string(parts[1])
		str = string(parts[2])
	}

	mt := MessageType{
		Code:      code,
		Event:     evt,
		Structure: str,
	}

	return mt, nil
}

// String returns the string representation of the MessageType.
func (t MessageType) String() string {
	if len(t.Structure) > 0 {
		return t.Structure
	}

	return fmt.Sprintf("%s_%s", t.Code, t.Event)
}

// Encode encodes the MessageType into its byte representation using the given delimiters.
func (t MessageType) Encode(delims *Delimiters) []byte {
	return delims.Join([][]byte{
		[]byte(t.Code),
		[]byte(t.Event),
		[]byte(t.Structure),
	}, ComponentDelimiter)
}

// MessageHeader represents the parsed MSH header fields of an HL7 message.
type MessageHeader struct {
	Delimiters                    *Delimiters `json:"delimiters" hl7:"MSH.1"`
	SendingApplication            string      `json:"sending_application" hl7:"MSH.3"`
	SendingFacility               string      `json:"sending_facility" hl7:"MSH.4"`
	ReceivingApplication          string      `json:"receiving_application" hl7:"MSH.5"`
	ReceivingFacility             string      `json:"receiving_facility" hl7:"MSH.6"`
	MessageDate                   time.Time   `json:"message_date" hl7:"MSH.7"`
	Security                      string      `json:"security" hl7:"MSH.8"`
	MessageCode                   string      `json:"type" hl7:"MSH.9.1"`
	TriggerEvent                  string      `json:"event_trigger" hl7:"MSH.9.2"`
	MessageStructure              string      `json:"message_structure" hl7:"MSH.9.3"`
	ControlID                     string      `json:"control_id" hl7:"MSH.10"`
	ProcessingID                  string      `json:"processing_id" hl7:"MSH.11"`
	VersionID                     string      `json:"version_id" hl7:"MSH.12"`
	SequenceNumber                int         `json:"sequence_number" hl7:"MSH.13"`
	ContinuationPointer           string      `json:"continuation_pointer" hl7:"MSH.14"`
	AcceptAcknowledgmentType      string      `json:"accept_acknowledgment" hl7:"MSH.15"`
	ApplicationAcknowledgmentType string      `json:"application_acknowledgment" hl7:"MSH.16"`
	CountryCode                   string      `json:"country_code" hl7:"MSH.17"`
	CharacterSet                  string      `json:"character_set" hl7:"MSH.18"`
	PrincipalLanguage             string      `json:"principal_language" hl7:"MSH.19"`
	AlternateCharacterSet         string      `json:"alternate_character_set" hl7:"MSH.20"`
	ProfileID                     string      `json:"profile_id" hl7:"MSH.21"`
	SendingResponsibleOrgCode     string      `json:"sending_responsible_org_code" hl7:"MSH.22"`
	ReceivingResponsibleOrgCode   string      `json:"receiving_responsible_org_code" hl7:"MSH.23"`
	SendingNetworkAddress         string      `json:"sending_network_address" hl7:"MSH.24"`
	ReceivingNetworkAddress       string      `json:"receiving_network_address" hl7:"MSH.25"`
}

// ParseHeader parses an MSH message header from raw message bytes.
func ParseHeader(b []byte) (*MessageHeader, error) {
	p, err := newRawParser(b, OnlyHeader())
	if err != nil {
		return nil, err
	}

	msg, err := p.Parse()
	if err != nil {
		return nil, err
	}

	return newMessageHeader(msg)
}

func newMessageHeader(m *RawMessage) (*MessageHeader, error) {
	seg := m.Segments("MSH")
	if len(seg) == 0 {
		return nil, fmt.Errorf("missing MSH segment")
	}

	if len(seg) > 1 {
		return nil, fmt.Errorf("multiple MSH segments found")
	}

	//msh := seg[0]

	var h MessageHeader
	if err := Unmarshal(m.v, &h); err != nil {
		return nil, err
	}

	return &h, nil
}

// MessageType constructs a MessageType from the MSH header.
func (h *MessageHeader) MessageType() MessageType {
	return MessageType{
		Code:      h.MessageCode,
		Event:     h.TriggerEvent,
		Structure: h.MessageStructure,
	}
}

// MarshalHL7 marshals the MessageHeader into HL7 wire format bytes.
func (h *MessageHeader) MarshalHL7() ([]byte, error) {
	panic("not implemented")
}

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
	HeaderTypeMessage HeaderType = iota //MSH
	HeaderTypeFile                      //FHS
	HeaderTypeBatch                     //BHS
)

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

type TrailerType int

const (
	TrailerTypeNone  TrailerType = iota //_
	TrailerTypeFile                     //FTS
	TrailerTypeBatch                    //BTS
)

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

type MessageType struct {
	Code      string
	Event     string
	Structure string
}

func (t MessageType) String() string {
	if len(t.Structure) > 0 {
		return t.Structure
	}

	return fmt.Sprintf("%s_%s", t.Code, t.Event)
}

func (t MessageType) Encode(delims *Delimiters) []byte {
	return delims.Join([][]byte{
		[]byte(t.Code),
		[]byte(t.Event),
		[]byte(t.Structure),
	}, ComponentDelimiter)
}

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
	if err := Unmarshal(m, &h); err != nil {
		return nil, err
	}

	return &h, nil
}

func (h *MessageHeader) MessageType() MessageType {
	return MessageType{
		Code:      h.MessageCode,
		Event:     h.TriggerEvent,
		Structure: h.MessageStructure,
	}
}

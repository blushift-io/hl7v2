package hl7v2

import (
	"fmt"
	"time"

	"github.com/spf13/cast"
)

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
	return t.Code + t.Event + t.Structure
}

type MessageHeader struct {
	*header
	Type                        MessageType
	ControlID                   string
	ProcessingID                string
	VersionID                   Version
	SequenceNumber              int
	ContinuationPointer         string
	AcceptAcknowledgment        string
	CountryCode                 string
	CharacterSet                string
	PrincipalLanguage           string
	AlternateCharacterSet       string
	ProfileID                   string
	SendingResponsibleOrgCode   string
	ReceivingResponsibleOrgCode string
	SendingNetworkAddress       string
	ReceivingNetworkAddress     string
}

func newMessageHeader(m *RawMessage) (*MessageHeader, error) {
	seg := m.Segments("MSH")
	if len(seg) == 0 {
		return nil, fmt.Errorf("missing MSH segment")
	}

	if len(seg) > 1 {
		return nil, fmt.Errorf("multiple MSH segments found")
	}

	msh := seg[0]

	h := &header{
		typ:          HeaderTypeMessage,
		delims:       m.delims,
		sendingApp:   msh[2].String(),
		receivingApp: msh[5].String(),
		sendingFac:   msh[4].String(),
		receivingFac: msh[6].String(),
		dateTime:     cast.ToTime(msh[7].String()),
		security:     msh[8].String(),
	}

	return &MessageHeader{
		header: h,
	}, nil
}

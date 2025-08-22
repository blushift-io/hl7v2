package hl7v2

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	ErrEncodingTooShort = errors.New("invalid delimiters, encoding chars must be at least 5 bytes")
	ErrEncodingInvalid  = errors.New("invalid encoding chars")
)

// Type is an enum of valid delimiters
//
//go:generate enumer -type=DelimiterType --linecomment -output=delimiter_enums.go
type DelimiterType int

const (
	InvalidDelimiter      DelimiterType = iota //invalid
	SegmentDelimiter                           //segment
	FieldDelimiter                             //field
	RepetitionDelimiter                        //repetition
	ComponentDelimiter                         //component
	SubcomponentDelimiter                      //subcomponent
	EscapeDelimiter                            //escape
	TruncationDelimiter                        //truncation

	segmentDelimiter = byte('\r')
)

// Delimiter is a byte representation of a single delimiter
type Delimiter byte

func (d Delimiter) String() string {
	return string(d)
}

func (d Delimiter) Byte() byte {
	return byte(d)
}

func (d Delimiter) Bytes() []byte {
	return []byte{byte(d)}
}

// Delimiters holds the message delimiters
type Delimiters struct {
	Segment      Delimiter
	Field        Delimiter
	Repetition   Delimiter
	Component    Delimiter
	Subcomponent Delimiter
	Escape       Delimiter
	Truncation   Delimiter
}

type SetDelimiter func(*Delimiters)

// DefaultDelimiters returns the default delimiters
func DefaultDelimiters() *Delimiters {
	return &Delimiters{
		Segment:      '\r',
		Field:        '|',
		Repetition:   '~',
		Component:    '^',
		Subcomponent: '&',
		Escape:       '\\',
		Truncation:   '#',
	}
}

func NewDelimiters(setters ...SetDelimiter) *Delimiters {
	delims := DefaultDelimiters()

	for _, setter := range setters {
		setter(delims)
	}

	return delims
}

func SetSegment(delim Delimiter) SetDelimiter {
	return func(d *Delimiters) {
		d.Segment = delim
	}
}

func SetField(delim Delimiter) SetDelimiter {
	return func(d *Delimiters) {
		d.Field = delim
	}
}

func SetComponent(delim Delimiter) SetDelimiter {
	return func(d *Delimiters) {
		d.Component = delim
	}
}

func SetRepetition(delim Delimiter) SetDelimiter {
	return func(d *Delimiters) {
		d.Repetition = delim
	}
}

func SetSubcomponent(delim Delimiter) SetDelimiter {
	return func(d *Delimiters) {
		d.Subcomponent = delim
	}
}

func SetTruncation(delim Delimiter) SetDelimiter {
	return func(d *Delimiters) {
		d.Truncation = delim
	}
}

func ParseDelimiters(b []byte) (*Delimiters, error) {
	delims := DefaultDelimiters()

	if len(b) < 5 {
		return delims, ErrEncodingTooShort
	}

	chk := make(map[byte]bool)
	for _, v := range b {
		_, ok := chk[v]
		if ok {
			return delims, ErrEncodingInvalid
		}

		chk[v] = true
	}

	delims.Field = Delimiter(b[0])
	delims.Component = Delimiter(b[1])
	delims.Repetition = Delimiter(b[2])
	delims.Escape = Delimiter(b[3])
	delims.Subcomponent = Delimiter(b[4])

	if len(b) >= 6 {
		delims.Truncation = Delimiter(b[5])
	}

	return delims, nil
}

func (d *Delimiters) fieldValue() Value {
	return NewValue(d.Field.Bytes())
}

func (d *Delimiters) encodingCharsValue() Value {
	b := &bytes.Buffer{}
	b.Write(d.Component.Bytes())
	b.Write(d.Repetition.Bytes())
	b.Write(d.Escape.Bytes())
	b.Write(d.Subcomponent.Bytes())

	return NewValue(b.Bytes())
}

func (d *Delimiters) Join(b [][]byte, typ DelimiterType) []byte {
	var delim Delimiter
	switch typ {
	case SegmentDelimiter:
		delim = d.Segment
	case FieldDelimiter:
		delim = d.Field
	case RepetitionDelimiter:
		delim = d.Repetition
	case ComponentDelimiter:
		delim = d.Component
	case SubcomponentDelimiter:
		delim = d.Subcomponent
	case EscapeDelimiter:
		delim = d.Escape
	case TruncationDelimiter:
		delim = d.Truncation
	default:
		panic("invalid delimiter type")
	}

	return bytes.Join(b, delim.Bytes())
}

func (d *Delimiters) Split(b []byte, typ DelimiterType) [][]byte {
	var delim Delimiter
	switch typ {
	case SegmentDelimiter:
		delim = d.Segment
	case FieldDelimiter:
		delim = d.Field
	case RepetitionDelimiter:
		delim = d.Repetition
	case ComponentDelimiter:
		delim = d.Component
	case SubcomponentDelimiter:
		delim = d.Subcomponent
	case EscapeDelimiter:
		delim = d.Escape
	case TruncationDelimiter:
		delim = d.Truncation
	default:
		panic("invalid delimiter type")
	}

	return bytes.Split(b, delim.Bytes())
}

func getDelimiters(delims ...*Delimiters) *Delimiters {
	if len(delims) > 0 && delims[0] != nil {
		return delims[0]
	}

	return DefaultDelimiters()
}

func getEncodingChars(b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(b)

	hdr := make([]byte, 3)
	_, err := buf.Read(hdr)
	if err != nil {
		return nil, fmt.Errorf("error reading message header: %w", err)
	}

	if string(hdr) != "MSH" {
		return nil, fmt.Errorf("invalid message: expected header 'MSH', got '%s'", hdr)
	}

	fs, err := buf.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("error reading segment delimiter: %w", err)
	}

	encChars, err := buf.ReadBytes(fs)
	if err != nil {
		return nil, fmt.Errorf("error reading encoding characters: %w", err)
	}

	return append([]byte{fs}, encChars[:len(encChars)-1]...), nil
}

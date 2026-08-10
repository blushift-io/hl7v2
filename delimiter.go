package hl7v2

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrEncodingTooShort is returned when encoding characters are fewer than 5 bytes.
	ErrEncodingTooShort = errors.New("invalid delimiters, encoding chars must be at least 5 bytes")
	// ErrEncodingInvalid is returned when duplicate characters are found in encoding chars.
	ErrEncodingInvalid  = errors.New("invalid encoding chars")
)

// Type is an enum of valid delimiters
//
//go:generate enumer -type=DelimiterType --linecomment -output=delimiter_enums.go
type DelimiterType int

const (
	// InvalidDelimiter represents an invalid delimiter.
	InvalidDelimiter      DelimiterType = iota //invalid
	// SegmentDelimiter represents the segment delimiter character (\r).
	SegmentDelimiter                           //segment
	// FieldDelimiter represents the field delimiter character (|).
	FieldDelimiter                             //field
	// RepetitionDelimiter represents the field repetition delimiter character (~).
	RepetitionDelimiter                        //repetition
	// ComponentDelimiter represents the component delimiter character (^).
	ComponentDelimiter                         //component
	// SubcomponentDelimiter represents the subcomponent delimiter character (&).
	SubcomponentDelimiter                      //subcomponent
	// EscapeDelimiter represents the escape character (\).
	EscapeDelimiter                            //escape
	// TruncationDelimiter represents the truncation delimiter character (#).
	TruncationDelimiter                        //truncation

	segmentDelimiter = byte('\r')
)

// Delimiter is a byte representation of a single delimiter
type Delimiter byte

// String returns the delimiter as a string.
func (d Delimiter) String() string {
	return string(d)
}

// Byte returns the delimiter as a byte.
func (d Delimiter) Byte() byte {
	return byte(d)
}

// Bytes returns the delimiter as a byte slice.
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

// SetDelimiter is a functional option for configuring Delimiters.
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

// NewDelimiters creates a new Delimiters instance with default settings and applies any functional options.
func NewDelimiters(setters ...SetDelimiter) *Delimiters {
	delims := DefaultDelimiters()

	for _, setter := range setters {
		setter(delims)
	}

	return delims
}

// SetSegment returns an option function setting the segment delimiter.
func SetSegment(delim Delimiter) SetDelimiter {
	return func(d *Delimiters) {
		d.Segment = delim
	}
}

// SetField returns an option function setting the field delimiter.
func SetField(delim Delimiter) SetDelimiter {
	return func(d *Delimiters) {
		d.Field = delim
	}
}

// SetComponent returns an option function setting the component delimiter.
func SetComponent(delim Delimiter) SetDelimiter {
	return func(d *Delimiters) {
		d.Component = delim
	}
}

// SetRepetition returns an option function setting the repetition delimiter.
func SetRepetition(delim Delimiter) SetDelimiter {
	return func(d *Delimiters) {
		d.Repetition = delim
	}
}

// SetSubcomponent returns an option function setting the subcomponent delimiter.
func SetSubcomponent(delim Delimiter) SetDelimiter {
	return func(d *Delimiters) {
		d.Subcomponent = delim
	}
}

// SetTruncation returns an option function setting the truncation delimiter.
func SetTruncation(delim Delimiter) SetDelimiter {
	return func(d *Delimiters) {
		d.Truncation = delim
	}
}

// ParseDelimiters parses delimiters from an HL7 encoding characters byte slice.
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

// Join joins byte slices using the specified delimiter type.
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

// Split splits a byte slice using the specified delimiter type.
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

// Escaper creates a new Escaper configured with these delimiters.
func (d *Delimiters) Escaper() *Escaper {
	return NewEscaper(d)
}

// FieldSeparatorValue returns the field delimiter as a Value.
func (d *Delimiters) FieldSeparatorValue() Value {
	return NewValue(d.Field.Bytes())
}

// EncodingCharsValue returns the encoding characters as a Value.
func (d *Delimiters) EncodingCharsValue() Value {
	b := &bytes.Buffer{}
	b.Write(d.Component.Bytes())
	b.Write(d.Repetition.Bytes())
	b.Write(d.Escape.Bytes())
	b.Write(d.Subcomponent.Bytes())

	return NewValue(b.Bytes())
}

// Escaper handles escaping and unescaping of HL7 delimiter characters.
type Escaper struct {
	delims   *Delimiters
	escape   strings.Replacer
	unescape strings.Replacer
}

// NewEscaper returns a new Escaper configured with the provided delimiters.
func NewEscaper(delims *Delimiters) *Escaper {
	wrap := func(s string) string {
		return fmt.Sprintf("%s%s%s", delims.Escape.String(), s, delims.Escape.String())
	}
	return &Escaper{
		delims: delims,
		escape: *strings.NewReplacer(
			delims.Field.String(), wrap("F"),
			delims.Repetition.String(), wrap("R"),
			delims.Component.String(), wrap("S"),
			delims.Subcomponent.String(), wrap("T"),
			delims.Escape.String(), wrap("E"),
			"\n", wrap("X0A"),
			"\r", wrap("X0D"),
		),
		unescape: *strings.NewReplacer(
			wrap("F"), delims.Field.String(),
			wrap("R"), delims.Repetition.String(),
			wrap("S"), delims.Component.String(),
			wrap("T"), delims.Subcomponent.String(),
			wrap("E"), delims.Escape.String(),
			wrap("X0A"), "\n",
			wrap("X0D"), "\r",
			wrap(".br"), "\r",
		),
	}
}

// Escape escapes HL7 delimiter characters in the input byte slice.
func (e *Escaper) Escape(b []byte) []byte {
	return []byte(e.escape.Replace(string(b)))
}

// Unescape unescapes HL7 escape sequences in the input byte slice.
func (e *Escaper) Unescape(b []byte) []byte {
	return []byte(e.unescape.Replace(string(b)))
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

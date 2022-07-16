package hl7v2

import "errors"

var (
	ErrEncodingTooShort = errors.New("invalid delimiters, encoding chars must be at least 5 bytes")
	ErrEncodingInvalid  = errors.New("invalid encoding chars")
)

//go:generate enumer -type=DelimiterType --linecomment
//DelimiterType is an enum of valid delimiters
type DelimiterType int

const (
	DelimiterInvalid      DelimiterType = iota //invalid
	DelimiterSegment                           //segment
	DelimiterField                             //field
	DelimiterRepetition                        //repetition
	DelimiterComponent                         //component
	DelimiterSubcomponent                      //subcomponent
	DelimiterEscape                            //escape
	DelimiterTruncation                        //truncation
)

//Delimiter is a byte representation of a single delimiter
type Delimiter byte

func (d Delimiter) String() string {
	return string(d)
}

//Delimiters holds the message delimiters
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

//DefaultDelimiters returns the default delimiters
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

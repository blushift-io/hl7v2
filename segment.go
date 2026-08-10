package hl7v2

import (
	"bytes"
	"fmt"

	"github.com/blushift-io/hl7v2/query"
)

// RawSegmentGroup represents an unparsed group of segments.
type RawSegmentGroup []RawSegment

// RawSegment represents an unparsed segment containing raw fields.
type RawSegment []RawField

func (s RawSegment) collect(delims ...*Delimiters) []byte {
	d := getDelimiters(delims...)

	var b [][]byte
	for _, field := range s {
		b = append(b, field.collect(d))
	}

	var res []byte
	if s.ID() == "MSH" {
		res = append(res, b[0]...)
		res = append(res, b[1]...)
		res = append(res, d.Join(b[2:], FieldDelimiter)...)
	} else {
		res = d.Join(b, FieldDelimiter)
	}

	return res

}

// Query queries the RawSegment at the specified location.
func (s RawSegment) Query(loc query.Location, delims ...*Delimiters) (*Value, error) {
	if loc.Field == 0 {
		val := NewValue(s.collect(delims...))
		return &val, nil
	}

	if int(loc.Field) > len(s)-1 {
		return nil, ErrValueNotFound
	}

	return s[loc.Field].Query(loc, delims...)
}

// ID returns the segment identifier name (e.g., "MSH", "PID").
func (s RawSegment) ID() string {
	if len(s) == 0 {
		return ""
	}

	return s[0].String()
}

// Value wraps the RawSegment bytes in a Value.
func (s RawSegment) Value(delims ...*Delimiters) Value {
	return NewValue(s.collect(delims...))
}

// String converts the RawSegment to its string representation.
func (s RawSegment) String(delims ...*Delimiters) string {
	if len(s) == 0 {
		return ""
	}

	return string(s.collect(delims...))
}

// Segment represents a parsed HL7 segment containing fields.
type Segment struct {
	id       string
	pos      int
	parent   Element
	children []*Field
}

// NewSegment creates a new Segment instance.
func NewSegment(parent Element, pos int, raw RawSegment) *Segment {
	seg := &Segment{
		id:       raw.ID(),
		pos:      pos,
		parent:   parent,
		children: make([]*Field, len(raw)),
	}

	for i, field := range raw {
		seg.children[i] = NewField(seg, i, field)
	}

	return seg
}

// Type returns the element type for Segment.
func (s *Segment) Type() ElementType {
	return ElementSegment
}

// Name returns the segment identifier name.
func (s *Segment) Name() string {
	return s.id
}

// Delimiters returns the message delimiters.
func (s *Segment) Delimiters() *Delimiters {
	if s.parent == nil {
		return DefaultDelimiters()
	}

	return s.parent.Delimiters()
}

// Header returns the message header.
func (s *Segment) Header() *MessageHeader {
	if s.parent == nil {
		return nil
	}

	return s.parent.Header()
}

// Parent returns the parent element.
func (s *Segment) Parent() Element {
	return s.parent
}

// Children returns the fields as child elements.
func (s *Segment) Children() []Element {
	return makeElements(s.children...)
}

// Length returns the number of fields in the segment.
func (s *Segment) Length() int {
	return len(s.children)
}

// Position returns the 1-indexed position of the segment in the message.
func (s *Segment) Position() int {
	return s.pos
}

// Index returns the segment set ID index if present.
func (s *Segment) Index() int {
	if len(s.children) < 2 {
		return 0
	}

	return s.children[1].Value().Int()
}

// Location returns the query Location of the segment.
func (s *Segment) Location() query.Location {
	loc := query.Location{}
	if s == nil {
		return loc
	}

	loc.Segment = s.id

	rep := 0
	if len(s.children) > 1 {
		rep = s.children[1].Value().Int()
		if rep > 0 {
			loc.SegmentRep = &rep
		}
	}

	return loc
}

// Value returns the Value of the segment.
func (s *Segment) Value(escape ...bool) Value {
	var b [][]byte

	start := 0
	if s.id == "MSH" {
		start = 3
		fd := [][]byte{
			s.children[0].Value().Bytes(),
			s.children[1].Value().Bytes(),
			s.children[2].Value().Bytes(),
		}

		b = append(b, bytes.Join(fd, []byte{}))
	}

	for i := start; i < len(s.children); i++ {
		v := s.children[i].Value(escape...)
		b = append(b, v.Bytes())
	}

	return NewValue(s.Delimiters().Join(b, FieldDelimiter))
}

// GetLocation resolves an element inside the segment by its location query.
func (s *Segment) GetLocation(loc query.Location) (Element, error) {
	if loc.Segment != s.id {
		return nil, fmt.Errorf("segment mismatch: querying %s, got %s", s.id, loc.Segment)

	}
	if int(loc.Field) > len(s.children) {
		return nil, fmt.Errorf("field %d not found", loc.Field)
	}

	return s.children[loc.Field].GetLocation(loc)
}

// SetLocation updates the value at the given location inside the segment.
func (s *Segment) SetLocation(loc query.Location, val Value) error {
	if loc.Segment != s.id {
		return fmt.Errorf("segment mismatch: querying %s, got %s", s.id, loc.Segment)
	}

	if int(loc.Field) > len(s.children) {
		return fmt.Errorf("field %d not found", loc.Field)
	}

	return s.children[loc.Field].SetLocation(loc, val)
}

// Encode encodes the segment into its byte representation.
func (s *Segment) Encode() ([]byte, error) {
	return s.Value().Bytes(), nil
}

// Append appends a field element to the segment.
func (s *Segment) Append(el Element) error {
	if el.Type() != ElementField {
		return fmt.Errorf("cannot append type %s to segment", el.Type())
	}

	ins, ok := el.(*Field)
	if !ok {
		return fmt.Errorf("cannot append type %s to segment", el.Type())
	}

	ins.parent = s
	ins.pos = len(s.children) + 1
	s.children = append(s.children, ins)

	return nil
}

// Fields returns the child fields in the segment.
func (s *Segment) Fields() []*Field {
	return s.children
}

// Field returns the field at the specified 0-indexed position.
func (s *Segment) Field(index int) (*Field, error) {
	if index < 0 || index > len(s.children)-1 {
		return nil, fmt.Errorf("field %d not found", index)
	}

	return s.children[index], nil
}

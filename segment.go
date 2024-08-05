package hl7v2

import (
	"fmt"

	"github.com/blushift-io/hl7v2/query"
)

type RawSegmentGroup []RawSegment
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

func (s RawSegment) ID() string {
	if len(s) == 0 {
		return ""
	}

	return s[0].String()
}

func (s RawSegment) Value(delims ...*Delimiters) Value {
	return NewValue(s.collect(delims...))
}

func (s RawSegment) String(delims ...*Delimiters) string {
	if len(s) == 0 {
		return ""
	}

	return string(s.collect(delims...))
}

type Segment struct {
	id       string
	pos      int
	parent   Element
	children []*Field
}

func newSegment(parent Element, pos int, raw RawSegment) *Segment {
	seg := &Segment{
		id:       raw.ID(),
		pos:      pos,
		parent:   parent,
		children: make([]*Field, len(raw)),
	}

	for i, field := range raw {
		seg.children[i] = newField(seg, i+1, field)
	}

	return seg
}

func (s *Segment) Type() ElementType {
	return ElementSegment
}

func (s *Segment) Name() string {
	return s.id
}

func (s *Segment) Delimiters() *Delimiters {
	return s.parent.Delimiters()
}

func (s *Segment) Parent() Element {
	return s.parent
}

func (s *Segment) Children() []Element {
	return makeElements(s.children...)
}

func (s *Segment) Length() int {
	return len(s.children)
}

func (s *Segment) Position() int {
	return s.pos
}

func (s *Segment) Index() int {
	if len(s.children) < 2 {
		return 0
	}

	return s.children[1].Value().Int()
}

func (s *Segment) Location() query.Location {
	rep := s.children[1].Value().Int()

	return query.Location{
		Segment:    s.Name(),
		SegmentRep: rep,
	}
}

func (s *Segment) GetLocation(loc query.Location) (Element, error) {
	if loc.Segment != s.id {
		return nil, fmt.Errorf("segment mismatch: querying %s, got %s", s.id, loc.Segment)

	}
	if int(loc.Field) > len(s.children) {
		return nil, fmt.Errorf("field %d not found", loc.Field)
	}

	return s.children[loc.Field].GetLocation(loc)
}

func (s *Segment) Value() Value {
	var b [][]byte

	for _, fld := range s.children {
		b = append(b, fld.Value().Bytes())
	}

	return NewValue(s.Delimiters().Join(b, FieldDelimiter))
}

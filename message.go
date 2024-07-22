package hl7v2

import (
	"fmt"
	"io"
	"os"

	"github.com/blushift-io/hl7v2/query"
)

type Message struct {
	raw      *RawMessage
	parent   Element
	pos      int
	delims   *Delimiters
	segments []*Segment
	header   *MessageHeader
}

func NewMessage(r io.Reader, opts ...ParserOption) (*Message, error) {
	raw, err := NewRawMessage(r, opts...)
	if err != nil {
		return nil, err
	}

	return newMessage(nil, 0, raw)
}

func NewMessageFromBytes(b []byte, opts ...ParserOption) (*Message, error) {
	raw, err := NewRawMessageFromBytes(b, opts...)
	if err != nil {
		return nil, err
	}

	return newMessage(nil, 0, raw)
}

func NewMessageFromFile(f string, opts ...ParserOption) (*Message, error) {
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	return NewMessageFromBytes(b, opts...)
}

func newMessage(parent Element, pos int, raw *RawMessage) (*Message, error) {
	msg := &Message{
		raw:      raw,
		parent:   parent,
		pos:      pos,
		delims:   raw.delims,
		segments: make([]*Segment, len(raw.segs)),
	}

	for i, seg := range raw.segs {
		msg.segments[i] = newSegment(msg, i+1, seg)
	}

	h, err := newMessageHeader(raw)
	if err != nil {
		return nil, err
	}

	msg.header = h

	return msg, nil
}

func (m *Message) Type() ElementType {
	return ElementMessage
}

func (m *Message) Name() string {
	return m.header.Type.String()
}

func (m *Message) Header() *MessageHeader {
	return m.header
}

func (m *Message) Delimiters() *Delimiters {
	return m.delims
}

func (m *Message) Parent() Element {
	return m.parent
}

func (m *Message) Children() []Element {
	return makeElements(m.segments...)
}

func (m *Message) Position() int {
	return m.pos
}

func (m *Message) Value() Value {
	var b [][]byte

	for _, seg := range m.segments {
		b = append(b, seg.Value().Bytes())
	}

	return NewValue(m.delims.Join(b, SegmentDelimiter))
}

func (m *Message) Query(q string) (Element, error) {
	loc, err := query.ParseLocation(q)
	if err != nil {
		return nil, err
	}

	return m.GetLocation(loc)
}

func (m *Message) Location() query.Location {
	return query.Location{}
}

func (m *Message) GetLocation(loc query.Location) (Element, error) {
	if loc.Segment == "" {
		return nil, fmt.Errorf("invalid message query: missing segment")
	}

	segs := m.getSegment(loc.Segment)
	if len(segs) == 0 {
		return nil, fmt.Errorf("segment '%s' not found", loc.Segment)
	}

	if loc.SegmentRep == 0 {
		return segs[0].GetLocation(loc)
	}

	if loc.SegmentRep > len(segs) {
		return nil, fmt.Errorf("segment '%s' repetition %d not found", loc.Segment, loc.SegmentRep)
	}

	return segs[loc.SegmentRep-1].GetLocation(loc)
}

func (m *Message) getSegment(id string) []*Segment {
	var segs []*Segment

	for _, seg := range m.segments {
		if seg.Name() == id {
			segs = append(segs, seg)
		}
	}

	return segs
}

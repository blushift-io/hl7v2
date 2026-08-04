package hl7v2

import (
	"fmt"
	"io"
	"os"

	"github.com/blushift-io/hl7v2/query"
	"github.com/blushift-io/hl7v2/schema"
)

type Message struct {
	raw *RawMessage
	sch *schema.Message

	parent   Element
	pos      int
	delims   *Delimiters
	segments []*Segment
	segCount map[string]int
	header   *MessageHeader
}

func ReadMessage(r io.Reader, opts ...ParserOption) (*Message, error) {
	raw, err := ReadRaw(r, opts...)
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

	return NewMessage(b, opts...)
}

func NewMessage(b []byte, opts ...ParserOption) (*Message, error) {
	raw, err := ParseRaw(b, opts...)
	if err != nil {
		return nil, err
	}

	return newMessage(nil, 0, raw)
}

func newMessage(parent Element, pos int, raw *RawMessage) (*Message, error) {
	msg := &Message{
		raw:      raw,
		parent:   parent,
		pos:      pos,
		delims:   raw.delims,
		segments: make([]*Segment, len(raw.segs)),
		segCount: make(map[string]int),
	}

	h, err := newMessageHeader(raw)
	if err != nil {
		return nil, err
	}

	msg.header = h

	for i, seg := range raw.segs {
		msg.segments[i] = NewSegment(msg, i+1, seg)
		msg.segCount[seg.ID()]++
	}

	return msg, nil
}

func (m *Message) Type() ElementType {
	return ElementMessage
}

func (m *Message) Name() string {
	return m.header.MessageType().String()
}

func (m *Message) Delimiters() *Delimiters {
	return m.delims
}

func (m *Message) Header() *MessageHeader {
	return m.header
}

func (m *Message) Parent() Element {
	return m.parent
}

func (m *Message) Children() []Element {
	return makeElements(m.segments...)
}

func (m *Message) Length() int {
	return len(m.segments)
}

func (m *Message) Position() int {
	return m.pos
}

func (m *Message) Location() query.Location {
	return query.Location{}
}

func (m *Message) Value(escape ...bool) Value {
	var b [][]byte

	for _, seg := range m.segments {
		b = append(b, seg.Value(escape...).Bytes())
	}

	return NewValue(m.delims.Join(b, SegmentDelimiter))
}

func (m *Message) GetLocation(loc query.Location) (Element, error) {
	if loc.Segment == "" {
		return nil, fmt.Errorf("invalid message query: missing segment")
	}

	segs := m.getSegment(loc.Segment)
	if len(segs) == 0 {
		return nil, fmt.Errorf("segment '%s' not found", loc.Segment)
	}

	rep := 0
	if loc.SegmentRep != nil {
		rep = *loc.SegmentRep
	}

	if rep == 0 {
		return segs[0].GetLocation(loc)
	}

	if rep > len(segs) {
		return nil, fmt.Errorf("segment '%s' repetition %d not found", loc.Segment, loc.SegmentRep)
	}

	return segs[rep-1].GetLocation(loc)
}

func (m *Message) SetLocation(loc query.Location, val Value) error {
	el, err := m.GetLocation(loc)
	if err != nil {
		return err
	}

	return el.SetLocation(loc, val)
}

func (m *Message) Append(el Element) error {
	if el.Type() != ElementSegment {
		return fmt.Errorf("cannot append %s to message", el.Type())
	}

	ins, ok := el.(*Segment)
	if !ok {
		return fmt.Errorf("cannot append %s to message", el.Type())
	}

	m.segments = append(m.segments, ins)
	m.segCount[ins.Name()]++
	ins.parent = m
	ins.pos = len(m.segments) + 1

	return nil
}

func (m *Message) Encode() ([]byte, error) {
	return m.Value(true).Bytes(), nil
}

func (m *Message) Raw() (*RawMessage, error) {
	return ParseRaw(m.Value(true).Bytes())
}

func (m *Message) Query(q string) (Element, error) {
	loc, err := query.ParseLocation(q)
	if err != nil {
		return nil, err
	}

	return m.GetLocation(loc)
}

func (m *Message) QueryValue(q string) (*Value, error) {
	el, err := m.Query(q)
	if err != nil {
		return nil, err
	}

	if el == nil {
		return nil, nil
	}

	v := el.Value()

	return &v, nil
}

func (m *Message) Schema() (*schema.Message, error) {
	sch := schema.Open(m.header.VersionID)
	if sch == nil {
		return nil, schema.ErrSchemaNotFound
	}

	typ := m.header.MessageType().String()
	msg := sch.Message(typ)
	if msg == nil {
		return nil, schema.ErrMessageTypeNotFound
	}

	return msg, nil
}

func (m *Message) Select(q string) ([]Element, error) {
	g, err := query.ParseGrammars(q)
	if err != nil {
		return nil, err
	}

	return m.doSelect(g)
}

func (m *Message) Segments(id ...string) []*Segment {
	if len(id) > 0 && id[0] != "" {
		return m.getSegment(id[0])
	}

	return m.segments
}

func (m *Message) Segment(name string, idx ...int) (*Segment, error) {
	setID := 0
	if len(idx) > 0 {
		setID = idx[0]
	}

	segs := m.getSegment(name)
	if len(segs) == 0 {
		return nil, fmt.Errorf("segment '%s' not found", name)
	}

	if setID == 0 {
		return segs[0], nil
	}

	for _, seg := range segs {
		if seg.Position() == setID {
			return seg, nil
		}
	}

	return nil, fmt.Errorf("segment '%s' repetition %d not found", name, setID)
}

func (m *Message) HasSegment(id string) bool {
	return m.segCount[id] > 0
}

func (m *Message) SegmentList() []string {
	var segs []string

	for _, seg := range m.segments {
		segs = append(segs, seg.Name())
	}

	return segs
}

func (m *Message) doSelect(g query.Grammars) ([]Element, error) {
	segs := m.SegmentList()
	if err := g.Validate(segs); err != nil {
		return nil, err
	}

	sel, err := g.Select(segs)
	if err != nil {
		return nil, err
	}

	var res []Element

	for _, seg := range sel {
		res = append(res, makeElements(m.getSegment(seg)...)...)
	}

	return res, nil
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

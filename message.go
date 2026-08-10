package hl7v2

import (
	"fmt"
	"io"
	"os"

	"github.com/blushift-io/hl7v2/query"
	"github.com/blushift-io/hl7v2/schema"
)

// Message represents a parsed HL7 message containing structured segments and fields.
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

// ReadMessage parses an HL7 message from an io.Reader.
func ReadMessage(r io.Reader, opts ...ParserOption) (*Message, error) {
	raw, err := ReadRaw(r, opts...)
	if err != nil {
		return nil, err
	}

	return newMessage(nil, 0, raw)
}

// NewMessageFromFile reads and parses an HL7 message from a file path.
func NewMessageFromFile(f string, opts ...ParserOption) (*Message, error) {
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	return NewMessage(b, opts...)
}

// NewMessage parses an HL7 message from a byte slice.
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

// Type returns the element type for Message.
func (m *Message) Type() ElementType {
	return ElementMessage
}

// Name returns the message type name.
func (m *Message) Name() string {
	return m.header.MessageType().String()
}

// Delimiters returns the message delimiters.
func (m *Message) Delimiters() *Delimiters {
	return m.delims
}

// Header returns the parsed MessageHeader.
func (m *Message) Header() *MessageHeader {
	return m.header
}

// Parent returns the parent element.
func (m *Message) Parent() Element {
	return m.parent
}

// Children returns the segments of the message as child elements.
func (m *Message) Children() []Element {
	return makeElements(m.segments...)
}

// Length returns the number of segments in the message.
func (m *Message) Length() int {
	return len(m.segments)
}

// Position returns the position of the message.
func (m *Message) Position() int {
	return m.pos
}

// Location returns the location query of the message.
func (m *Message) Location() query.Location {
	return query.Location{}
}

// Value returns the combined Value of the message segments.
func (m *Message) Value(escape ...bool) Value {
	var b [][]byte

	for _, seg := range m.segments {
		b = append(b, seg.Value(escape...).Bytes())
	}

	return NewValue(m.delims.Join(b, SegmentDelimiter))
}

// GetLocation resolves an element in the message by a query Location.
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

// SetLocation updates the value at the specified query Location.
func (m *Message) SetLocation(loc query.Location, val Value) error {
	el, err := m.GetLocation(loc)
	if err != nil {
		return err
	}

	return el.SetLocation(loc, val)
}

// Append appends a segment element to the message.
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

// Encode encodes the message into HL7 wire format bytes.
func (m *Message) Encode() ([]byte, error) {
	return m.Value(true).Bytes(), nil
}

// Raw returns a RawMessage parsed from the current Message bytes.
func (m *Message) Raw() (*RawMessage, error) {
	return ParseRaw(m.Value(true).Bytes())
}

// Query finds an element in the message using a location query string.
func (m *Message) Query(q string) (Element, error) {
	loc, err := query.ParseLocation(q)
	if err != nil {
		return nil, err
	}

	return m.GetLocation(loc)
}

// QueryValue retrieves the Value at the given location query string.
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

// Schema loads the schema specification matching the message header version and type.
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

// Select filters elements in the message using a grammar query string.
func (m *Message) Select(q string) ([]Element, error) {
	g, err := query.ParseGrammars(q)
	if err != nil {
		return nil, err
	}

	return m.doSelect(g)
}

// Segments returns segments matching the specified ID, or all segments if no ID is given.
func (m *Message) Segments(id ...string) []*Segment {
	if len(id) > 0 && id[0] != "" {
		return m.getSegment(id[0])
	}

	return m.segments
}

// Segment returns the segment matching the given name and optional position index.
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

// HasSegment checks if the message contains at least one segment with the given ID.
func (m *Message) HasSegment(id string) bool {
	return m.segCount[id] > 0
}

// SegmentList returns a slice of all segment names present in the message.
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

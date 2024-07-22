package hl7v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/blushift-io/hl7v2/query"
)

type RawMessage struct {
	v      []byte
	delims *Delimiters
	segIdx map[string]int
	segs   []RawSegment
}

func NewRawMessage(r io.Reader, opts ...ParserOption) (*RawMessage, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return NewRawMessageFromBytes(b, opts...)

}

func NewRawMessageFromBytes(b []byte, opts ...ParserOption) (*RawMessage, error) {
	p, err := newRawParser(b, opts...)
	if err != nil {
		return nil, err
	}

	return p.Parse()
}

func (m *RawMessage) Delimiters() *Delimiters {
	return m.delims
}

func (m *RawMessage) Segments(id ...string) []RawSegment {
	if len(id) == 0 {
		return m.segs
	}

	var segs []RawSegment
	for _, s := range m.segs {
		if s.ID() == id[0] {
			segs = append(segs, s)
		}
	}

	return segs
}

func (m *RawMessage) hasSegment(id string) bool {
	_, ok := m.segIdx[id]

	return ok
}

func (m *RawMessage) Query(loc query.Location) (*Value, error) {
	if loc.Segment == "" {
		return nil, fmt.Errorf("invalid query: missing segment")
	}

	if !m.hasSegment(loc.Segment) {
		return nil, fmt.Errorf("segment '%s' not found", loc.Segment)
	}

	segs := m.Segments(loc.Segment)
	segCnt := m.segIdx[loc.Segment]

	if segCnt > len(segs) {
		return nil, fmt.Errorf("segment '%s' repetition %d not found", loc.Segment, segCnt)
	}

	seg := segs[segCnt-1]

	return seg.Query(loc, m.delims)
}

func (m *RawMessage) append(seg RawSegment) {
	id := seg.ID()
	if _, ok := m.segIdx[id]; ok {
		m.segIdx[id]++
	} else {
		m.segIdx[id] = 1
	}

	m.segs = append(m.segs, seg)
}

func (m *RawMessage) JSON(pretty ...bool) ([]byte, error) {
	msg := jsonMessage{}

	for _, seg := range m.segs {
		msg.Segments = append(msg.Segments, extractJSONSegment(seg))
	}

	if len(pretty) > 0 && pretty[0] {
		return json.MarshalIndent(msg, "", "  ")
	}

	return json.Marshal(msg)
}

type rawParser struct {
	data   []byte
	delims *Delimiters
	opts   *ParserOptions

	msg *RawMessage
	seg RawSegment
	fld RawField
	rep RawRepetition
	cmp RawComponent
	sub RawSubcomponent
}

func newRawParser(data []byte, opts ...ParserOption) (*rawParser, error) {
	if len(data) < 8 {
		return nil, ErrMsgLength
	}

	return &rawParser{
		data: data,
		opts: NewParserOptions(opts...),
		msg: &RawMessage{
			v:      data,
			delims: DefaultDelimiters(),
			segIdx: make(map[string]int),
		},
	}, nil
}

func (p *rawParser) Parse() (*RawMessage, error) {
	for _, preParse := range p.opts.preParse {
		b, err := preParse(p.data)
		if err != nil {
			return nil, fmt.Errorf("error pre-parsing data: %w", err)
		}

		p.data = b
	}

	buf := bytes.NewBuffer(p.data)

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

	encChars = append([]byte{fs}, encChars[:len(encChars)-1]...)

	delims, err := ParseDelimiters(encChars)
	if err != nil {
		return nil, fmt.Errorf("error parsing delimiters: %w", err)
	}

	p.delims = delims
	p.msg.delims = delims

	p.seg = RawSegment{
		RawField{RawRepetition{RawComponent{RawSubcomponent(hdr)}}},
		RawField{RawRepetition{RawComponent{RawSubcomponent([]byte{fs})}}},
		RawField{RawRepetition{RawComponent{RawSubcomponent(encChars[1:])}}},
	}

	rem, err := buf.ReadBytes(delims.Segment.Byte())
	if err != nil {
		return nil, fmt.Errorf("error reading segment: %w", err)
	}

	p.parseSeg(rem)

	for {
		b, err := buf.ReadBytes(delims.Segment.Byte())
		if err != nil {
			if err == io.EOF {
				p.parseSeg(b)
				break
			}

			return nil, err
		}

		p.parseSeg(b)
	}

	return p.msg, nil
}

func (p *rawParser) parseSeg(b []byte) {
	b = bytes.TrimSpace(b)
	for _, f := range p.delims.Split(b, FieldDelimiter) {
		for _, r := range p.delims.Split(f, RepetitionDelimiter) {
			for _, c := range p.delims.Split(r, ComponentDelimiter) {
				for _, s := range p.delims.Split(c, SubcomponentDelimiter) {
					p.commitBuffer(s)
					p.commitSub()
				}
				p.commitComp()
			}
			p.commitRep()
		}
		p.commitField()
	}

	p.commitSeg()
}

func (p *rawParser) commitSeg() {
	p.msg.append(p.seg)
	p.seg = RawSegment{}
}

func (p *rawParser) commitField() {
	p.seg = append(p.seg, p.fld)
	p.fld = RawField{}
}

func (p *rawParser) commitRep() {
	p.fld = append(p.fld, p.rep)
	p.rep = RawRepetition{}
}

func (p *rawParser) commitComp() {
	p.rep = append(p.rep, p.cmp)
	p.cmp = RawComponent{}
}

func (p *rawParser) commitSub() {
	p.cmp = append(p.cmp, p.sub)
	p.sub = RawSubcomponent{}
}

func (p *rawParser) commitBuffer(b []byte) {
	p.sub = append(p.sub, b...)
}

type jsonMessage struct {
	Segments []jsonSegment `json:"segments"`
}

type jsonSegment struct {
	ID     string         `json:"id"`
	Fields map[string]any `json:"fields"`
}

func extractJSONSegment(seg RawSegment) jsonSegment {
	id := seg[0][0][0][0]
	flds := make(map[string]any)

	for fi, fld := range seg {
		if fi == 0 {
			continue
		}

		idx := fmt.Sprintf("%d", fi)
		for ri, rep := range fld {
			if len(fld) > 1 {
				idx = fmt.Sprintf("%s[%d]", idx, ri+1)
			}

			if len(rep) > 1 {
				idx = fmt.Sprintf("%s.%d", idx, ri+1)
			}

			for ci, cmp := range rep {
				if len(cmp) > 1 {
					idx = fmt.Sprintf("%s.%d", idx, ci+1)
				}

				for si, sub := range cmp {
					if len(cmp) > 1 {
						idx = fmt.Sprintf("%s-%d", idx, si+1)
					}

					flds[idx] = string(sub)
				}
			}
		}
	}

	return jsonSegment{
		ID:     string(id),
		Fields: flds,
	}
}

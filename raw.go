package hl7v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/blushift-io/hl7v2/query"
)

type RawMessage struct {
	v      []byte
	delims *Delimiters
	segIdx map[string]int
	segs   []RawSegment
}

func NewRawMessage(delims *Delimiters, segs ...RawSegment) *RawMessage {
	msg := &RawMessage{
		delims: delims,
		segIdx: make(map[string]int),
	}

	for _, s := range segs {
		msg.append(s)
	}

	msg.v = msg.Value().Bytes()

	return msg
}

func ReadRaw(r io.Reader, opts ...ParserOption) (*RawMessage, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return ParseRaw(b, opts...)

}

func ParseRaw(b []byte, opts ...ParserOption) (*RawMessage, error) {
	p, err := newRawParser(b, opts...)
	if err != nil {
		return nil, err
	}

	return p.Parse()
}

func (m *RawMessage) Delimiters() *Delimiters {
	return m.delims
}

func (m *RawMessage) ToMessage() (*Message, error) {
	return newMessage(nil, 0, m)
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

func (m *RawMessage) QueryValue(q string) (*Value, error) {
	loc, err := query.ParseLocation(q)
	if err != nil {
		return nil, fmt.Errorf("error parsing query: %w", err)
	}

	if loc.Segment == "" {
		if len(m.segs) > 0 {
			loc.Segment = m.segs[0].ID()
		} else {
			return nil, ErrElementNotFound
		}
	}

	if !m.hasSegment(loc.Segment) {
		return nil, ErrElementNotFound
	}

	segs := m.Segments(loc.Segment)
	segCnt := m.segIdx[loc.Segment]

	if segCnt > len(segs) {
		return nil, ErrElementNotFound
	}

	seg := segs[segCnt-1]

	v, err := seg.Query(loc, m.delims)
	if err != nil {
		return nil, fmt.Errorf("failed to query location %s: %w", q, err)
	}

	return v, nil
}

func (m *RawMessage) Value() Value {
	var b [][]byte

	for _, seg := range m.segs {
		b = append(b, seg.Value().Bytes())
	}

	return NewValue(m.delims.Join(b, SegmentDelimiter))
}

func (m *RawMessage) append(seg RawSegment) {
	id := seg.ID()

	m.segIdx[id] += 1

	m.segs = append(m.segs, seg)
}

func (m *RawMessage) toJSONMessage() jsonMessage {
	jmsg := jsonMessage{
		Segments: make([]jsonSegment, 0, len(m.segs)),
	}
	if m.delims != nil {
		jmsg.Delimiters = string(m.delims.Field) + m.delims.EncodingCharsValue().String()
	}

	for _, seg := range m.segs {
		jseg := extractJSONSegment(seg)
		if jseg.ID != "" {
			jmsg.Segments = append(jmsg.Segments, jseg)
		}
	}

	return jmsg
}

func (m *RawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}

	return json.Marshal(m.toJSONMessage())
}

func (m *RawMessage) UnmarshalJSON(b []byte) error {
	msg, err := ParseJSON(b)
	if err != nil {
		return err
	}

	*m = *msg
	return nil
}

func (m *RawMessage) JSON(pretty ...bool) ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}

	if len(pretty) > 0 && pretty[0] {
		return json.MarshalIndent(m.toJSONMessage(), "", "  ")
	}

	return m.MarshalJSON()
}

func (m *RawMessage) Decode() (*Message, error) {
	return newMessage(nil, 0, m)
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
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("error reading segment: %w", err)
	}

	p.parseSeg(rem)

	if err == io.EOF {
		return p.msg, nil
	}

	if p.opts.onlyHeader {
		return p.msg, nil
	}

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
	if len(b) == 0 {
		return
	}
	for _, f := range p.delims.Split(b, FieldDelimiter) {
		for _, r := range p.delims.Split(f, RepetitionDelimiter) {
			for _, c := range p.delims.Split(r, ComponentDelimiter) {
				for _, s := range p.delims.Split(c, SubcomponentDelimiter) {
					s = p.delims.Escaper().Unescape(s)

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
	Delimiters string        `json:"delimiters,omitempty"`
	Segments   []jsonSegment `json:"segments"`
}

type jsonSegment struct {
	ID     string         `json:"id"`
	Fields map[string]any `json:"fields"`
}

func extractJSONSegment(seg RawSegment) jsonSegment {
	if len(seg) == 0 {
		return jsonSegment{Fields: make(map[string]any)}
	}

	id := seg.ID()
	flds := make(map[string]any)

	for fi := 1; fi < len(seg); fi++ {
		fld := seg[fi]
		numReps := len(fld)
		if numReps == 0 {
			continue
		}

		for ri := 0; ri < numReps; ri++ {
			rep := fld[ri]
			numComps := len(rep)
			if numComps == 0 {
				continue
			}

			for ci := 0; ci < numComps; ci++ {
				cmp := rep[ci]
				numSubs := len(cmp)
				if numSubs == 0 {
					continue
				}

				for si := 0; si < numSubs; si++ {
					subStr := string(cmp[si])
					if len(subStr) == 0 {
						continue
					}

					key := formatFieldKey(fi, ri, ci, si, numReps, numComps, numSubs)
					flds[key] = subStr
				}
			}
		}
	}

	return jsonSegment{
		ID:     id,
		Fields: flds,
	}
}

func formatFieldKey(fi, ri, ci, si, numReps, numComps, numSubs int) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(fi))

	if numReps > 1 {
		sb.WriteString("[")
		sb.WriteString(strconv.Itoa(ri))
		sb.WriteString("]")
	}

	if numComps > 1 || ci > 0 || numSubs > 1 || si > 0 {
		sb.WriteString(".")
		sb.WriteString(strconv.Itoa(ci + 1))
	}

	if numSubs > 1 || si > 0 {
		sb.WriteString(".")
		sb.WriteString(strconv.Itoa(si + 1))
	}

	return sb.String()
}

func ReadJSON(r io.Reader) (*RawMessage, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return ParseJSON(b)
}

var keyRegex = regexp.MustCompile(`^(\d+)(?:\[(\d+)\])?(?:\.(\d+)(?:[\.-](\d+))?)?$`)

func ParseJSON(b []byte) (*RawMessage, error) {
	var jmsg jsonMessage
	if err := json.Unmarshal(b, &jmsg); err != nil {
		return nil, fmt.Errorf("hl7v2: failed to unmarshal json: %w", err)
	}

	var delims *Delimiters
	var err error
	if jmsg.Delimiters != "" {
		delims, err = ParseDelimiters([]byte(jmsg.Delimiters))
		if err != nil {
			delims = DefaultDelimiters()
		}
	} else {
		var msh1, msh2 string
		for _, seg := range jmsg.Segments {
			if seg.ID == "MSH" {
				if v, ok := seg.Fields["1"]; ok {
					msh1 = fmt.Sprint(v)
				}
				if v, ok := seg.Fields["2"]; ok {
					msh2 = fmt.Sprint(v)
				}
				break
			}
		}
		if msh1 != "" && msh2 != "" {
			delims, err = ParseDelimiters([]byte(msh1 + msh2))
			if err != nil {
				delims = DefaultDelimiters()
			}
		} else {
			delims = DefaultDelimiters()
		}
	}

	rawSegs := make([]RawSegment, 0, len(jmsg.Segments))
	for _, jseg := range jmsg.Segments {
		seg, err := buildRawSegment(jseg)
		if err != nil {
			return nil, err
		}
		rawSegs = append(rawSegs, seg)
	}

	return NewRawMessage(delims, rawSegs...), nil
}

func buildRawSegment(jseg jsonSegment) (RawSegment, error) {
	var seg RawSegment

	ensureField(&seg, 0)
	ensureRep(&seg[0], 0)
	ensureComp(&seg[0][0], 0)
	ensureSub(&seg[0][0][0], 0)
	seg[0][0][0][0] = RawSubcomponent(jseg.ID)

	for key, valAny := range jseg.Fields {
		valStr := fmt.Sprint(valAny)

		fieldIdx, repIdx, compIdx, subIdx, err := parseFieldKey(key)
		if err != nil {
			return nil, fmt.Errorf("hl7v2: invalid field key %q in segment %s: %w", key, jseg.ID, err)
		}

		ensureField(&seg, fieldIdx)
		ensureRep(&seg[fieldIdx], repIdx)
		ensureComp(&seg[fieldIdx][repIdx], compIdx-1)
		ensureSub(&seg[fieldIdx][repIdx][compIdx-1], subIdx-1)

		seg[fieldIdx][repIdx][compIdx-1][subIdx-1] = RawSubcomponent(valStr)
	}

	return seg, nil
}

func parseFieldKey(key string) (fieldIdx, repIdx, compIdx, subIdx int, err error) {
	matches := keyRegex.FindStringSubmatch(key)
	if matches == nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid key format")
	}

	fieldIdx, _ = strconv.Atoi(matches[1])

	repIdx = 0
	if matches[2] != "" {
		repIdx, _ = strconv.Atoi(matches[2])
	}

	compIdx = 1
	if matches[3] != "" {
		compIdx, _ = strconv.Atoi(matches[3])
	}

	subIdx = 1
	if matches[4] != "" {
		subIdx, _ = strconv.Atoi(matches[4])
	}

	if fieldIdx < 1 || repIdx < 0 || compIdx < 1 || subIdx < 1 {
		return 0, 0, 0, 0, fmt.Errorf("indices must be >= 1 for fields/components/subcomponents, >= 0 for reps")
	}

	return fieldIdx, repIdx, compIdx, subIdx, nil
}

func ensureField(seg *RawSegment, fieldIdx int) {
	for len(*seg) <= fieldIdx {
		*seg = append(*seg, RawField{})
	}
}

func ensureRep(fld *RawField, repIdx int) {
	for len(*fld) <= repIdx {
		*fld = append(*fld, RawRepetition{})
	}
}

func ensureComp(rep *RawRepetition, compIdx int) {
	for len(*rep) <= compIdx {
		*rep = append(*rep, RawComponent{})
	}
}

func ensureSub(cmp *RawComponent, subIdx int) {
	for len(*cmp) <= subIdx {
		*cmp = append(*cmp, RawSubcomponent{})
	}
}

package message

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/blushift-io/hl7v2"
)

var (
	ErrMsgLength = errors.New("invalid message: message length must be at least 8 bytes")
)

type Context struct {
	Message *Message
}

type Transform func([]byte) ([]byte, error)
type Processor func(*Context) error

type Options struct {
	preParse  []Transform
	onCommit  []Processor
	postParse []Processor
}

type Option func(*Options)

func DefaultOptions() *Options {
	return &Options{}
}

func NewOptions(opts ...Option) *Options {
	options := DefaultOptions()

	for _, opt := range opts {
		opt(options)
	}

	return options
}

func FixLineEndings() Option {
	return func(options *Options) {
		options.preParse = append(options.preParse, fixLineEndings)
	}
}

func PreProcess(t Transform) Option {
	return func(options *Options) {
		options.preParse = append(options.preParse, t)
	}
}

func PostProcess(p Processor) Option {
	return func(options *Options) {
		options.postParse = append(options.postParse, p)
	}
}

func CommitTransform(p Processor) Option {
	return func(options *Options) {
		options.onCommit = append(options.onCommit, p)
	}
}

type Parser struct {
	data   []byte
	delims *hl7v2.Delimiters
	opts   *Options

	msg *RawMessage
	seg RawSegment
	fld RawField
	rep RawRepetition
	cmp RawComponent
	sub RawSubcomponent
}

func NewParser(data []byte, opts ...Option) (*Parser, error) {
	if len(data) < 8 {
		return nil, ErrMsgLength
	}

	return &Parser{
		data: data,
		opts: NewOptions(opts...),
		msg: &RawMessage{
			delims: hl7v2.DefaultDelimiters(),
		},
	}, nil
}

func (p *Parser) Parse() (*RawMessage, error) {
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

	delims, err := hl7v2.ParseDelimiters(encChars)
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

func (p *Parser) parseSeg(b []byte) {
	b = bytes.TrimSpace(b)
	for _, f := range p.delims.Split(b, hl7v2.FieldDelimiter) {
		for _, r := range p.delims.Split(f, hl7v2.RepetitionDelimiter) {
			for _, c := range p.delims.Split(r, hl7v2.ComponentDelimiter) {
				for _, s := range p.delims.Split(c, hl7v2.SubcomponentDelimiter) {
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

func (p *Parser) commitSeg() {
	p.msg.Append(p.seg)
	p.seg = RawSegment{}
}

func (p *Parser) commitField() {
	p.seg = append(p.seg, p.fld)
	p.fld = RawField{}
}

func (p *Parser) commitRep() {
	p.fld = append(p.fld, p.rep)
	p.rep = RawRepetition{}
}

func (p *Parser) commitComp() {
	p.rep = append(p.rep, p.cmp)
	p.cmp = RawComponent{}
}

func (p *Parser) commitSub() {
	p.cmp = append(p.cmp, p.sub)
	p.sub = RawSubcomponent{}
}

func (p *Parser) commitBuffer(b []byte) {
	p.sub = append(p.sub, b...)
}

func fixLineEndings(b []byte) ([]byte, error) {
	b = bytes.ReplaceAll(b, []byte("\r\n"), []byte("\r"))
	return bytes.ReplaceAll(b, []byte{'\n'}, []byte{'\r'}), nil
}

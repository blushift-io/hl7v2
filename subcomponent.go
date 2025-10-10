package hl7v2

import (
	"fmt"

	"github.com/blushift-io/hl7v2/query"
)

type RawSubcomponent []byte

func (s RawSubcomponent) String() string {
	return string(s)
}

func (s RawSubcomponent) Value() Value {
	return NewValue(s)
}

type Subcomponent struct {
	v      Value
	pos    int
	parent Element
}

func NewSubcomponentValue(val Value) *Subcomponent {
	return NewSubcomponent(nil, 0, val)
}

func NewEmptySubcomponent() *Subcomponent {
	return NewSubcomponent(nil, 0, NewEmptyValue())
}

func NewSubcomponent(parent Element, pos int, v Value) *Subcomponent {
	delims := parent.Delimiters()
	v = v.Unescape(delims)

	return &Subcomponent{
		v:      v,
		pos:    pos,
		parent: parent,
	}
}

func (el *Subcomponent) Type() ElementType {
	return ElementSubcomponent
}

func (el *Subcomponent) Name() string {
	return fmt.Sprintf("%s.%d", el.parent.Name(), el.pos)
}

func (el *Subcomponent) Delimiters() *Delimiters {
	if el.parent == nil {
		return DefaultDelimiters()
	}

	return el.parent.Delimiters()
}

func (el *Subcomponent) Header() *MessageHeader {
	if el.parent == nil {
		return nil
	}

	return el.parent.Header()
}

func (el *Subcomponent) Parent() Element {
	return el.parent
}

func (el *Subcomponent) Children() []Element {
	return nil
}

func (el *Subcomponent) Length() int {
	return 1
}

func (el *Subcomponent) Position() int {
	return el.pos
}

func (el *Subcomponent) Location() query.Location {
	loc := query.Location{}
	if el == nil {
		return loc
	}

	if el.parent != nil {
		loc = el.parent.Location()
	}

	loc.Subcomponent = el.pos

	return loc
}

func (el *Subcomponent) Value(escape ...bool) Value {
	if len(escape) > 0 && escape[0] {
		return el.v.Escape(el.Delimiters())
	}

	return el.v
}

func (el *Subcomponent) GetLocation(loc query.Location) (Element, error) {
	return el, nil
}

func (el *Subcomponent) SetLocation(loc query.Location, val Value) error {
	delims := el.Delimiters()
	val = val.Unescape(delims)

	el.v = val

	return nil
}

func (el *Subcomponent) Append(_ Element) error {
	return fmt.Errorf("cannot append to subcomponent")
}

func (el *Subcomponent) Encode() ([]byte, error) {
	return el.v.Bytes(), nil
}

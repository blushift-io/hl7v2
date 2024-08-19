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

func NewSubcomponent(val Value) *Subcomponent {
	return newSubcomponent(nil, 0, val)
}

func NewEmptySubcomponent() *Subcomponent {
	return newSubcomponent(nil, 0, NewEmptyValue())
}

func newSubcomponent(parent Element, pos int, v Value) *Subcomponent {
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
	if el.parent == nil {
		return query.Location{}
	}

	loc := el.parent.Location()
	loc.Subcomponent = el.pos

	return loc
}

func (el *Subcomponent) Value() Value {
	return el.v
}

func (el *Subcomponent) GetLocation(loc query.Location) (Element, error) {
	return el, nil
}

func (el *Subcomponent) SetLocation(loc query.Location, val Value) error {
	el.v = val

	return nil
}

func (el *Subcomponent) Append(_ Element) error {
	return fmt.Errorf("cannot append to subcomponent")
}

func (el *Subcomponent) Encode() ([]byte, error) {
	return el.v.Bytes(), nil
}

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
	return el.parent.Delimiters()
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
	loc := el.parent.Location()
	loc.Subcomponent = el.pos

	return loc
}

func (el *Subcomponent) GetLocation(loc query.Location) (Element, error) {
	return el, nil
}

func (el *Subcomponent) Value() Value {
	return el.v
}

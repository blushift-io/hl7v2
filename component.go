package hl7v2

import (
	"fmt"

	"github.com/blushift-io/hl7v2/query"
)

type RawComponent []RawSubcomponent

func (c RawComponent) collect(delims ...*Delimiters) []byte {
	d := getDelimiters(delims...)

	var b [][]byte
	for _, sub := range c {
		b = append(b, sub)
	}

	return d.Join(b, SubcomponentDelimiter)
}

func (c RawComponent) String(delims ...*Delimiters) string {
	if len(c) == 0 {
		return ""
	}

	return string(c.collect(delims...))
}

func (c RawComponent) Value() Value {
	return NewValue(c.collect())
}

func (c RawComponent) Query(loc query.Location, delims ...*Delimiters) (*Value, error) {
	if loc.Subcomponent == 0 {
		val := NewValue(c.collect(delims...))
		return &val, nil
	}

	if int(loc.Subcomponent) > len(c) {
		return nil, fmt.Errorf("subcomponent %d not found", loc.Subcomponent)
	}

	val := c[loc.Subcomponent-1].Value()

	return &val, nil
}

type Component struct {
	parent   Element
	children []*Subcomponent
	pos      int
}

func newComponent(parent Element, pos int, raw RawComponent) *Component {
	comp := &Component{
		parent:   parent,
		children: make([]*Subcomponent, len(raw)),
		pos:      pos,
	}

	for i, sub := range raw {
		comp.children[i] = newSubcomponent(comp, i+1, sub.Value())
	}

	return comp
}

func (el *Component) Type() ElementType {
	return ElementComponent
}

func (el *Component) Name() string {
	return ""
}

func (el *Component) Delimiters() *Delimiters {
	return el.parent.Delimiters()
}

func (el *Component) Parent() Element {
	return el.parent
}

func (el *Component) Children() []Element {
	return makeElements(el.children...)
}

func (el *Component) Position() int {
	return el.pos
}

func (el *Component) Value() Value {
	var b [][]byte

	for _, sub := range el.children {
		b = append(b, sub.v.Bytes())
	}

	return NewValue(el.Delimiters().Join(b, SubcomponentDelimiter))
}

func (el *Component) GetLocation(loc query.Location) (Element, error) {
	if loc.Subcomponent == 0 {
		return el, nil
	}

	if int(loc.Subcomponent) > len(el.children) {
		return nil, fmt.Errorf("subcomponent %d not found", loc.Subcomponent)
	}

	return el.children[loc.Subcomponent-1], nil
}

func (el *Component) Location() query.Location {
	return query.Location{
		Component: el.pos,
	}
}

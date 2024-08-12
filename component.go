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

func (c RawComponent) Value(delims ...*Delimiters) Value {
	return NewValue(c.collect(delims...))
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

func (c *Component) Type() ElementType {
	return ElementComponent
}

func (c *Component) Name() string {
	return ""
}

func (c *Component) Delimiters() *Delimiters {
	return c.parent.Delimiters()
}

func (c *Component) Parent() Element {
	return c.parent
}

func (c *Component) Children() []Element {
	return makeElements(c.children...)
}

func (c *Component) Length() int {
	return len(c.children)
}

func (c *Component) Position() int {
	return c.pos
}

func (c *Component) Location() query.Location {
	if c.parent == nil {
		return query.Location{}
	}

	loc := c.parent.Location()
	loc.Component = c.pos
	return loc
}

func (c *Component) Value() Value {
	var b [][]byte

	for _, sub := range c.children {
		b = append(b, sub.v.Bytes())
	}

	return NewValue(c.Delimiters().Join(b, SubcomponentDelimiter))
}

func (c *Component) GetLocation(loc query.Location) (Element, error) {
	if loc.Subcomponent == 0 {
		return c.children[0], nil
	}

	if int(loc.Subcomponent) > len(c.children) {
		return nil, fmt.Errorf("subcomponent %d not found", loc.Subcomponent)
	}

	return c.children[loc.Subcomponent-1], nil
}

func (c *Component) SetLocation(loc query.Location, val Value) error {
	el, err := c.GetLocation(loc)
	if err != nil {
		return err
	}

	return el.SetLocation(loc, val)
}

func (c *Component) Append(ne Element) error {
	if ne.Type() != ElementSubcomponent {
		return fmt.Errorf("cannot append type %s to component", c.Type())
	}

	ins, ok := ne.(*Subcomponent)
	if !ok {
		return fmt.Errorf("cannot append type %s to component", c.Type())
	}

	ins.parent = c
	ins.pos = len(c.children) + 1
	c.children = append(c.children, ne.(*Subcomponent))

	return nil
}

func (c *Component) Encode() ([]byte, error) {
	return c.Value().Bytes(), nil
}

func (c *Component) Subcomponent(index int) (*Subcomponent, error) {
	if index < 0 || index > len(c.children) {
		return nil, fmt.Errorf("subcomponent %d not found", index)
	}

	return c.children[index-1], nil
}

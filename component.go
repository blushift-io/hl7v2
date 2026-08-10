package hl7v2

import (
	"fmt"

	"github.com/blushift-io/hl7v2/query"
)

// RawComponent represents an unparsed HL7 component containing raw subcomponents.
type RawComponent []RawSubcomponent

func (c RawComponent) collect(delims ...*Delimiters) []byte {
	d := getDelimiters(delims...)

	var b [][]byte
	for _, sub := range c {
		b = append(b, sub)
	}

	return d.Join(b, SubcomponentDelimiter)
}

// String converts the RawComponent to its string representation.
func (c RawComponent) String(delims ...*Delimiters) string {
	if len(c) == 0 {
		return ""
	}

	return string(c.collect(delims...))
}

// Value wraps the RawComponent bytes in a Value.
func (c RawComponent) Value(delims ...*Delimiters) Value {
	return NewValue(c.collect(delims...))
}

// Query queries the RawComponent at the specified location.
func (c RawComponent) Query(loc query.Location, delims ...*Delimiters) (*Value, error) {
	if loc.Subcomponent == 0 {
		val := NewValue(c.collect(delims...))
		return &val, nil
	}

	if int(loc.Subcomponent) > len(c) {
		return nil, ErrElementNotFound
	}

	val := c[loc.Subcomponent-1].Value()

	return &val, nil
}

// Component represents an HL7 component element within a repetition or field.
type Component struct {
	parent   Element
	children []*Subcomponent
	pos      int
}

// NewComponent creates a new Component instance with the given parent, position, and raw content.
func NewComponent(parent Element, pos int, raw RawComponent) *Component {
	comp := &Component{
		parent:   parent,
		children: make([]*Subcomponent, len(raw)),
		pos:      pos,
	}

	for i, sub := range raw {
		comp.children[i] = NewSubcomponent(comp, i+1, sub.Value())
	}

	return comp
}

// Type returns the element type for Component.
func (c *Component) Type() ElementType {
	return ElementComponent
}

// Name returns the name of the component.
func (c *Component) Name() string {
	return ""
}

// Delimiters returns the message delimiters.
func (c *Component) Delimiters() *Delimiters {
	return c.parent.Delimiters()
}

// Header returns the message header.
func (c *Component) Header() *MessageHeader {
	if c.parent == nil {
		return nil
	}

	return c.parent.Header()
}

// Parent returns the parent element.
func (c *Component) Parent() Element {
	return c.parent
}

// Children returns the subcomponents as child elements.
func (c *Component) Children() []Element {
	return makeElements(c.children...)
}

// Length returns the number of subcomponents in the component.
func (c *Component) Length() int {
	return len(c.children)
}

// Position returns the 1-indexed position of the component within its parent.
func (c *Component) Position() int {
	return c.pos
}

// Location returns the query Location of the component.
func (c *Component) Location() query.Location {
	loc := query.Location{}
	if c == nil {
		return loc
	}

	if c.parent != nil {
		loc = c.parent.Location()
	}

	loc.Component = c.pos
	return loc
}

// Value returns the Value of the component.
func (c *Component) Value(escape ...bool) Value {
	var b [][]byte

	for _, sub := range c.children {
		var v []byte
		if len(escape) > 0 && escape[0] {
			v = sub.Value().Escape(c.Delimiters()).Bytes()
		} else {
			v = sub.Value().Bytes()
		}

		b = append(b, v)
	}

	return NewValue(c.Delimiters().Join(b, SubcomponentDelimiter))
}

// GetLocation resolves an element inside the component by its location query.
func (c *Component) GetLocation(loc query.Location) (Element, error) {
	if c == nil {
		return nil, fmt.Errorf("component is nil")
	}

	if loc.Subcomponent == 0 {
		return c.children[0], nil
	}

	if int(loc.Subcomponent) > len(c.children) {
		return nil, fmt.Errorf("subcomponent %d not found", loc.Subcomponent)
	}

	return c.children[loc.Subcomponent-1], nil
}

// SetLocation updates the value at the given location inside the component.
func (c *Component) SetLocation(loc query.Location, val Value) error {
	el, err := c.GetLocation(loc)
	if err != nil {
		return err
	}

	return el.SetLocation(loc, val)
}

// Append appends a subcomponent element to the component.
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

// Encode encodes the component into its byte representation.
func (c *Component) Encode() ([]byte, error) {
	return c.Value().Bytes(), nil
}

// Subcomponents returns the child subcomponents.
func (c *Component) Subcomponents() []*Subcomponent {
	return c.children
}

// Subcomponent returns the subcomponent at the specified 1-indexed position.
func (c *Component) Subcomponent(index int) (*Subcomponent, error) {
	if index < 0 || index > len(c.children) {
		return nil, fmt.Errorf("subcomponent %d not found", index)
	}

	return c.children[index-1], nil
}

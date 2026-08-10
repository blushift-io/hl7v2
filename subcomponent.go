package hl7v2

import (
	"fmt"

	"github.com/blushift-io/hl7v2/query"
)

// RawSubcomponent represents an unparsed byte slice of a subcomponent.
type RawSubcomponent []byte

// String converts the RawSubcomponent to a string.
func (s RawSubcomponent) String() string {
	return string(s)
}

// Value wraps the RawSubcomponent in a Value.
func (s RawSubcomponent) Value() Value {
	return NewValue(s)
}

// Subcomponent represents an HL7 subcomponent element.
type Subcomponent struct {
	v      Value
	pos    int
	parent Element
}

// NewSubcomponentValue creates a standalone Subcomponent with the given Value.
func NewSubcomponentValue(val Value) *Subcomponent {
	return NewSubcomponent(nil, 0, val)
}

// NewEmptySubcomponent creates an empty Subcomponent.
func NewEmptySubcomponent() *Subcomponent {
	return NewSubcomponent(nil, 0, NewEmptyValue())
}

// NewSubcomponent creates a new Subcomponent instance.
func NewSubcomponent(parent Element, pos int, v Value) *Subcomponent {
	delims := parent.Delimiters()
	v = v.Unescape(delims)

	return &Subcomponent{
		v:      v,
		pos:    pos,
		parent: parent,
	}
}

// Type returns the element type for Subcomponent.
func (el *Subcomponent) Type() ElementType {
	return ElementSubcomponent
}

// Name returns the name of the subcomponent.
func (el *Subcomponent) Name() string {
	return fmt.Sprintf("%s.%d", el.parent.Name(), el.pos)
}

// Delimiters returns the message delimiters.
func (el *Subcomponent) Delimiters() *Delimiters {
	if el.parent == nil {
		return DefaultDelimiters()
	}

	return el.parent.Delimiters()
}

// Header returns the message header.
func (el *Subcomponent) Header() *MessageHeader {
	if el.parent == nil {
		return nil
	}

	return el.parent.Header()
}

// Parent returns the parent element.
func (el *Subcomponent) Parent() Element {
	return el.parent
}

// Children returns nil as subcomponents have no children.
func (el *Subcomponent) Children() []Element {
	return nil
}

// Length returns 1 for a subcomponent.
func (el *Subcomponent) Length() int {
	return 1
}

// Position returns the 1-indexed position of the subcomponent.
func (el *Subcomponent) Position() int {
	return el.pos
}

// Location returns the query Location of the subcomponent.
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

// Value returns the Value of the subcomponent.
func (el *Subcomponent) Value(escape ...bool) Value {
	if len(escape) > 0 && escape[0] {
		return el.v.Escape(el.Delimiters())
	}

	return el.v
}

// GetLocation returns the subcomponent itself.
func (el *Subcomponent) GetLocation(loc query.Location) (Element, error) {
	return el, nil
}

// SetLocation updates the value of the subcomponent.
func (el *Subcomponent) SetLocation(loc query.Location, val Value) error {
	delims := el.Delimiters()
	val = val.Unescape(delims)

	el.v = val

	return nil
}

// Append returns an error as subcomponents cannot have child elements appended.
func (el *Subcomponent) Append(_ Element) error {
	return fmt.Errorf("cannot append to subcomponent")
}

// Encode encodes the subcomponent into its byte slice.
func (el *Subcomponent) Encode() ([]byte, error) {
	return el.v.Bytes(), nil
}

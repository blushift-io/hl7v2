package hl7v2

import (
	"fmt"

	"github.com/blushift-io/hl7v2/query"
)

// RawRepetition represents an unparsed field repetition containing raw components.
type RawRepetition []RawComponent

func (r RawRepetition) collect(delims ...*Delimiters) []byte {
	d := getDelimiters(delims...)

	var b [][]byte
	for _, comp := range r {
		b = append(b, comp.collect(d))
	}

	return d.Join(b, ComponentDelimiter)
}

// String converts the RawRepetition to its string representation.
func (r RawRepetition) String(delims ...*Delimiters) string {
	if len(r) == 0 {
		return ""
	}

	return string(r.collect(delims...))
}

// Value wraps the RawRepetition bytes in a Value.
func (r RawRepetition) Value(delims ...*Delimiters) Value {
	return NewValue(r.collect(delims...))
}

// Query queries the RawRepetition at the specified location.
func (r RawRepetition) Query(loc query.Location, delims ...*Delimiters) (*Value, error) {
	if loc.Component == 0 {
		val := NewValue(r.collect(delims...))
		return &val, nil
	}

	if int(loc.Component) > len(r) {
		return nil, ErrElementNotFound
	}

	return r[loc.Component-1].Query(loc, delims...)
}

// Repetition represents an HL7 field repetition element.
type Repetition struct {
	parent   Element
	children []*Component
	pos      int
}

// NewRepetition creates a new Repetition instance.
func NewRepetition(parent Element, pos int, raw RawRepetition) *Repetition {
	rep := &Repetition{
		parent:   parent,
		children: make([]*Component, len(raw)),
		pos:      pos,
	}

	for i, comp := range raw {
		rep.children[i] = NewComponent(rep, i+1, comp)
	}

	return rep
}

// Type returns the element type for Repetition.
func (r *Repetition) Type() ElementType {
	return ElementRepetition
}

// Name returns the repetition name.
func (r *Repetition) Name() string {
	return fmt.Sprintf("%s[%d]", r.parent.Name(), r.pos)
}

// Delimiters returns the message delimiters.
func (r *Repetition) Delimiters() *Delimiters {
	if r.parent == nil {
		return DefaultDelimiters()
	}

	return r.parent.Delimiters()
}

// Header returns the message header.
func (r *Repetition) Header() *MessageHeader {
	if r.parent == nil {
		return nil
	}

	return r.parent.Header()
}

// Parent returns the parent element.
func (r *Repetition) Parent() Element {
	return r.parent
}

// Children returns the components as child elements.
func (r *Repetition) Children() []Element {
	return makeElements(r.children...)
}

// Length returns the number of components in the repetition.
func (r *Repetition) Length() int {
	return len(r.children)
}

// Position returns the repetition index position.
func (r *Repetition) Position() int {
	return r.pos
}

// Location returns the query Location of the repetition.
func (r *Repetition) Location() query.Location {
	loc := query.Location{}
	if r == nil {
		return loc
	}

	if r.parent != nil {
		loc = r.parent.Location()
	}

	loc.FieldRep = &r.pos

	return loc
}

// Value returns the Value of the repetition.
func (r *Repetition) Value(escape ...bool) Value {
	var b [][]byte

	for _, comp := range r.children {
		b = append(b, comp.Value(escape...).Bytes())
	}

	return NewValue(r.Delimiters().Join(b, ComponentDelimiter))
}

// GetLocation resolves an element inside the repetition by its location query.
func (r *Repetition) GetLocation(loc query.Location) (Element, error) {
	if loc.Component == 0 {
		return r.children[0].GetLocation(loc)
	}

	if int(loc.Component) > len(r.children) {
		return nil, fmt.Errorf("component %d not found", loc.Component)
	}

	return r.children[loc.Component-1].GetLocation(loc)
}

// SetLocation updates the value at the given location inside the repetition.
func (r *Repetition) SetLocation(loc query.Location, val Value) error {
	el, err := r.GetLocation(loc)
	if err != nil {
		return err
	}

	return el.SetLocation(loc, val)
}

// Append appends a component element to the repetition.
func (r *Repetition) Append(el Element) error {
	if el.Type() != ElementComponent {
		return fmt.Errorf("cannot append type %s to repetition", el.Type())
	}

	ins, ok := el.(*Component)
	if !ok {
		return fmt.Errorf("cannot append type %s to repetition", el.Type())
	}

	ins.parent = r
	ins.pos = len(r.children) + 1
	r.children = append(r.children, ins)

	return nil
}

// Encode encodes the repetition into its byte representation.
func (r *Repetition) Encode() ([]byte, error) {
	return r.Value().Bytes(), nil
}

// Components returns the child components.
func (r *Repetition) Components() []*Component {
	return r.children
}

// Component returns the component at the specified 0-indexed position.
func (r *Repetition) Component(idx int) (*Component, error) {
	if idx < 0 || idx > len(r.children) {
		return nil, fmt.Errorf("component %d not found", idx)
	}

	return r.children[idx], nil
}

package hl7v2

import (
	"fmt"

	"github.com/blushift-io/hl7v2/query"
)

// RawField represents an unparsed field containing raw repetitions.
type RawField []RawRepetition

func (f RawField) collect(delims ...*Delimiters) []byte {
	d := getDelimiters(delims...)

	var b [][]byte
	for _, rep := range f {
		b = append(b, rep.collect(d))
	}

	return d.Join(b, RepetitionDelimiter)
}

// String converts the RawField to its string representation.
func (f RawField) String(delims ...*Delimiters) string {
	if len(f) == 0 {
		return ""
	}

	return string(f.collect(delims...))
}

// Value wraps the RawField bytes in a Value.
func (f RawField) Value(delims ...*Delimiters) Value {
	return NewValue(f.collect(delims...))
}

// Query queries the RawField at the specified location.
func (f RawField) Query(loc query.Location, delims ...*Delimiters) (*Value, error) {
	if loc.FieldRep == nil {
		if loc.Component != 0 {
			return f.subquery(loc, delims...)
		}

		v := NewValue(f.collect(delims...))
		return &v, nil
	}

	rep := int(*loc.FieldRep)
	if rep > len(f)-1 {
		return nil, ErrElementNotFound
	}

	return f[rep].Query(loc, delims...)
}

func (f RawField) subquery(loc query.Location, delims ...*Delimiters) (*Value, error) {
	d := getDelimiters(delims...)
	var b [][]byte

	for _, rep := range f {
		v, err := rep.Query(loc, delims...)
		if err != nil {
			return nil, err
		}

		b = append(b, v.Bytes())
	}

	res := NewValue(d.Join(b, RepetitionDelimiter))
	return &res, nil
}

// Field represents an HL7 field element within a segment.
type Field struct {
	parent   Element
	children []*Repetition
	pos      int
}

// NewField creates a new Field instance.
func NewField(parent Element, pos int, raw RawField) *Field {
	fld := &Field{
		parent:   parent,
		children: make([]*Repetition, len(raw)),
		pos:      pos,
	}

	for i, rep := range raw {
		fld.children[i] = NewRepetition(fld, i, rep)
	}

	return fld
}

// Type returns the element type for Field.
func (f *Field) Type() ElementType {
	return ElementField
}

// Name returns the name of the field.
func (f *Field) Name() string {
	return fmt.Sprintf("%s.%d", f.parent.Name(), f.pos)
}

// Delimiters returns the message delimiters.
func (f *Field) Delimiters() *Delimiters {
	if f.parent == nil {
		return DefaultDelimiters()
	}

	return f.parent.Delimiters()
}

// Header returns the message header.
func (f *Field) Header() *MessageHeader {
	if f.parent == nil {
		return nil
	}

	return f.parent.Header()
}

// Parent returns the parent segment element.
func (f *Field) Parent() Element {
	return f.parent
}

// Children returns the repetitions as child elements.
func (f *Field) Children() []Element {
	return makeElements(f.children...)
}

// Length returns the number of repetitions in the field.
func (f *Field) Length() int {
	return len(f.children)
}

// Position returns the 1-indexed field position within the segment.
func (f *Field) Position() int {
	return f.pos
}

// Location returns the query Location of the field.
func (f *Field) Location() query.Location {
	loc := query.Location{}
	if f == nil {
		return loc
	}

	if f.parent == nil {
		loc = f.parent.Location()
	}

	loc.Field = f.pos

	return loc
}

// Value returns the Value of the field.
func (f *Field) Value(escape ...bool) Value {
	var b [][]byte

	for _, comp := range f.children {
		b = append(b, comp.Value(escape...).Bytes())
	}

	return NewValue(f.Delimiters().Join(b, RepetitionDelimiter))
}

// GetLocation resolves an element inside the field by its location query.
func (f *Field) GetLocation(loc query.Location) (Element, error) {
	hasReps := len(f.children) > 1
	rep := 0

	if loc.FieldRep == nil {
		if loc.Component == 0 {
			return f.children[0], nil
		}

		if loc.Component > 0 {
			if hasReps {
				return nil, fmt.Errorf("ambiguous location: field has repetitions and component specified")
			}
		}
	} else {
		rep = int(*loc.FieldRep)
	}

	if rep > len(f.children)-1 {
		return nil, fmt.Errorf("repetition %d not found", loc.FieldRep)
	}

	return f.children[rep].GetLocation(loc)
}

// SetLocation updates the value at the given location inside the field.
func (f *Field) SetLocation(loc query.Location, val Value) error {
	el, err := f.GetLocation(loc)
	if err != nil {
		return err
	}

	return el.SetLocation(loc, val)
}

// Append appends a repetition element to the field.
func (f *Field) Append(el Element) error {
	if el.Type() != ElementRepetition {
		return fmt.Errorf("cannot append type %s to field", el.Type())
	}

	ins, ok := el.(*Repetition)
	if !ok {
		return fmt.Errorf("cannot append type %s to field", el.Type())
	}

	ins.parent = f
	ins.pos = len(f.children)
	f.children = append(f.children, ins)

	return nil
}

// Encode encodes the field into its byte representation.
func (f *Field) Encode() ([]byte, error) {
	return f.Value().Bytes(), nil
}

// Repetitions returns the field repetitions.
func (f *Field) Repetitions() []*Repetition {
	return f.children
}

// Repetition returns the repetition at the specified 0-indexed position.
func (f *Field) Repetition(index int) (*Repetition, error) {
	if index < 0 || index > len(f.children)-1 {
		return nil, fmt.Errorf("repetition %d not found", index)
	}

	return f.children[index], nil
}

package hl7v2

import (
	"fmt"

	"github.com/blushift-io/hl7v2/query"
)

type RawField []RawRepetition

func (f RawField) collect(delims ...*Delimiters) []byte {
	d := getDelimiters(delims...)

	var b [][]byte
	for _, rep := range f {
		b = append(b, rep.collect(d))
	}

	return d.Join(b, RepetitionDelimiter)
}

func (f RawField) String(delims ...*Delimiters) string {
	if len(f) == 0 {
		return ""
	}

	return string(f.collect(delims...))
}

func (f RawField) Value(delims ...*Delimiters) Value {
	return NewValue(f.collect(delims...))
}

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

type Field struct {
	parent   Element
	children []*Repetition
	pos      int
}

func NewField(vals ...Value) *Field {
	raw := make(RawField, len(vals))

	for i, val := range vals {
		raw[i] = RawRepetition{RawComponent{RawSubcomponent(val.Bytes())}}
	}

	return newField(nil, 0, raw)
}

func newField(parent Element, pos int, raw RawField) *Field {
	fld := &Field{
		parent:   parent,
		children: make([]*Repetition, len(raw)),
		pos:      pos,
	}

	for i, rep := range raw {
		fld.children[i] = newRepetition(fld, i, rep)
	}

	return fld
}

func (f *Field) Type() ElementType {
	return ElementField
}

func (f *Field) Name() string {
	return fmt.Sprintf("%s.%d", f.parent.Name(), f.pos)
}

func (f *Field) Delimiters() *Delimiters {
	if f.parent == nil {
		return DefaultDelimiters()
	}

	return f.parent.Delimiters()
}

func (f *Field) Header() *MessageHeader {
	if f.parent == nil {
		return nil
	}

	return f.parent.Header()
}

func (f *Field) Parent() Element {
	return f.parent
}

func (f *Field) Children() []Element {
	return makeElements(f.children...)
}

func (f *Field) Length() int {
	return len(f.children)
}

func (f *Field) Position() int {
	return f.pos
}

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

func (f *Field) Value() Value {
	var b [][]byte

	for _, comp := range f.children {
		b = append(b, comp.Value().Bytes())
	}

	return NewValue(f.Delimiters().Join(b, RepetitionDelimiter))
}

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

func (f *Field) SetLocation(loc query.Location, val Value) error {
	el, err := f.GetLocation(loc)
	if err != nil {
		return err
	}

	return el.SetLocation(loc, val)
}

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

func (f *Field) Encode() ([]byte, error) {
	return f.Value().Bytes(), nil
}

func (f *Field) Repetitions() []*Repetition {
	return f.children
}

func (f *Field) Repetition(index int) (*Repetition, error) {
	if index < 0 || index > len(f.children)-1 {
		return nil, fmt.Errorf("repetition %d not found", index)
	}

	return f.children[index], nil
}

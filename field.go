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
	if int(loc.FieldRep) > len(f)-1 {
		return nil, fmt.Errorf("repetition %d not found", loc.FieldRep)
	}

	return f[loc.FieldRep].Query(loc, delims...)
}

type Field struct {
	parent   Element
	children []*Repetition
	pos      int
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
	return f.parent.Delimiters()
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
	loc := f.parent.Location()
	loc.Field = f.pos

	return loc
}

func (f *Field) GetLocation(loc query.Location) (Element, error) {
	if int(loc.FieldRep) > len(f.children)-1 {
		return nil, fmt.Errorf("repetition %d not found", loc.FieldRep)
	}

	return f.children[loc.FieldRep].GetLocation(loc)
}

func (f *Field) Value() Value {
	var b [][]byte

	for _, comp := range f.children {
		b = append(b, comp.Value().Bytes())
	}

	return NewValue(f.Delimiters().Join(b, RepetitionDelimiter))
}

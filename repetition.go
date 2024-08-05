package hl7v2

import (
	"fmt"

	"github.com/blushift-io/hl7v2/query"
)

type RawRepetition []RawComponent

func (r RawRepetition) collect(delims ...*Delimiters) []byte {
	d := getDelimiters(delims...)

	var b [][]byte
	for _, comp := range r {
		b = append(b, comp.collect(d))
	}

	return d.Join(b, ComponentDelimiter)
}

func (r RawRepetition) String(delims ...*Delimiters) string {
	if len(r) == 0 {
		return ""
	}

	return string(r.collect(delims...))
}

func (r RawRepetition) Value(delims ...*Delimiters) Value {
	return NewValue(r.collect(delims...))
}

func (r RawRepetition) Query(loc query.Location, delims ...*Delimiters) (*Value, error) {
	if loc.Component == 0 {
		val := NewValue(r.collect(delims...))
		return &val, nil
	}

	if int(loc.Component) > len(r) {
		return nil, fmt.Errorf("component %d not found", loc.Component)
	}

	return r[loc.Component-1].Query(loc, delims...)
}

type Repetition struct {
	parent   Element
	children []*Component
	pos      int
}

func newRepetition(parent Element, pos int, raw RawRepetition) *Repetition {
	rep := &Repetition{
		parent:   parent,
		children: make([]*Component, len(raw)),
		pos:      pos,
	}

	for i, comp := range raw {
		rep.children[i] = newComponent(rep, i+1, comp)
	}

	return rep
}

func (r *Repetition) Type() ElementType {
	return ElementRepetition
}

func (r *Repetition) Name() string {
	return fmt.Sprintf("%s[%d]", r.parent.Name(), r.pos)
}

func (r *Repetition) Delimiters() *Delimiters {
	return r.parent.Delimiters()
}

func (r *Repetition) Parent() Element {
	return r.parent
}

func (r *Repetition) Children() []Element {
	return makeElements(r.children...)
}

func (r *Repetition) Length() int {
	return len(r.children)
}

func (r *Repetition) Position() int {
	return r.pos
}

func (r *Repetition) Value() Value {
	var b [][]byte

	for _, comp := range r.children {
		b = append(b, comp.Value().Bytes())
	}

	return NewValue(r.Delimiters().Join(b, ComponentDelimiter))
}

func (r *Repetition) Location() query.Location {
	loc := r.parent.Location()
	loc.FieldRep = r.pos

	return loc
}

func (r *Repetition) GetLocation(loc query.Location) (Element, error) {
	if loc.Component == 0 {
		return r, nil
	}

	if int(loc.Component) > len(r.children) {
		return nil, fmt.Errorf("component %d not found", loc.Component)
	}

	return r.children[loc.Component-1].GetLocation(loc)
}

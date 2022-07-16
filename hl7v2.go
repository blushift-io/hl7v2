package hl7v2

import "fmt"

//go:generate enumer -type=ElementType --linecomment
//ElementType is an enum of valid element types
type ElementType int

const (
	ElementInvalid      ElementType = iota //invalid
	ElementBatch                           //batch
	ElementFile                            //file
	ElementMessage                         //message
	ElementSegment                         //segment
	ElementField                           //field
	ElementRepetition                      //repetition
	ElementComponent                       //component
	ElementSubcomponent                    //subcomponent
)

//Element interface for HL7 elements
type Element interface {
	Type() ElementType
	Name() string
	Delimiters() *Delimiters
	Parent() Element
	Children() []Element
	Position() any
	Value() Value
}

type element struct {
	typ      ElementType
	delims   *Delimiters
	parent   Element
	children []Element
	pos      any
}

func (el *element) Type() ElementType {
	return el.typ
}

func (el *element) Name() string {
	return ""
}

func (el *element) Delimiters() *Delimiters {
	if el.delims == nil {
		return el.parent.Delimiters()
	}

	return el.delims
}

func (el *element) Parent() Element {
	return el.parent
}

func (el *element) Children() []Element {
	return el.children
}

func (el *element) Position() any {
	return el.pos
}

func (el *element) Value() Value {
	return Value{}
}

func makeElements[T Element](els ...T) []Element {
	res := make([]Element, len(els))
	for i := range els {
		res[i] = els[i]
	}

	return res
}

type Value struct {
	v any
}

func (v Value) String() string {
	s, ok := v.v.(string)
	if !ok {
		return fmt.Sprintf("%v", v.v)
	}

	return s
}

func (v Value) Bytes() []byte {
	b, ok := v.v.([]byte)
	if !ok {
		return []byte(v.String())
	}

	return b
}

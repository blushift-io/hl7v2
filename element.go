package hl7v2

import "github.com/blushift-io/hl7v2/query"

//go:generate enumer -type=ElementType --linecomment -output=element_enums.go

type Queryable interface {
	Delimiters() *Delimiters
	QueryValue(q string) (*Value, error)
}

// ElementType is an enum of valid element types
type ElementType int

const (
	ElementInvalid      ElementType = iota //invalid
	ElementBatch                           //batch
	ElementFile                            //file
	ElementMessage                         //message
	ElementSegmentGroup                    //group
	ElementSegment                         //segment
	ElementField                           //field
	ElementRepetition                      //repetition
	ElementComponent                       //component
	ElementSubcomponent                    //subcomponent
)

// Element interface for HL7 elements
type Element interface {
	Type() ElementType
	Name() string
	Delimiters() *Delimiters
	Header() *MessageHeader
	Parent() Element
	Children() []Element
	Length() int
	Position() int
	Location() query.Location
	Value() Value
	GetLocation(query.Location) (Element, error)
	SetLocation(query.Location, Value) error
	Append(Element) error
	Encode() ([]byte, error)
}

type ElementIterator struct {
	root   Element
	cursor int
}

func NewIterator(el Element) *ElementIterator {
	return &ElementIterator{
		root: el,
	}
}

func (i *ElementIterator) Next() Element {
	ch := i.root.Children()
	if i.cursor > len(ch)-1 {
		return nil
	}

	el := ch[i.cursor]
	i.cursor++

	return el
}

func makeElements[T Element](els ...T) []Element {
	res := make([]Element, len(els))
	for i := range els {
		res[i] = els[i]
	}

	return res
}

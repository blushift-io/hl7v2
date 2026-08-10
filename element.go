package hl7v2

import "github.com/blushift-io/hl7v2/query"

//go:generate enumer -type=ElementType --linecomment -output=element_enums.go

// Queryable is an interface implemented by elements that can be queried by a location string.
type Queryable interface {
	Delimiters() *Delimiters
	QueryValue(q string) (*Value, error)
}

// ElementType is an enum of valid element types
type ElementType int

const (
	// ElementInvalid represents an invalid or uninitialized element type.
	ElementInvalid      ElementType = iota //invalid
	// ElementBatch represents an HL7 batch element type.
	ElementBatch                           //batch
	// ElementFile represents an HL7 file element type.
	ElementFile                            //file
	// ElementMessage represents an HL7 message element type.
	ElementMessage                         //message
	// ElementSegmentGroup represents an HL7 segment group element type.
	ElementSegmentGroup                    //group
	// ElementSegment represents an HL7 segment element type.
	ElementSegment                         //segment
	// ElementField represents an HL7 field element type.
	ElementField                           //field
	// ElementRepetition represents an HL7 field repetition element type.
	ElementRepetition                      //repetition
	// ElementComponent represents an HL7 component element type.
	ElementComponent                       //component
	// ElementSubcomponent represents an HL7 subcomponent element type.
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
	Value(escape ...bool) Value
	GetLocation(query.Location) (Element, error)
	SetLocation(query.Location, Value) error
	Append(Element) error
	Encode() ([]byte, error)
}

// ElementIterator provides sequential iteration over child elements of an HL7 Element.
type ElementIterator struct {
	root   Element
	cursor int
}

// NewIterator creates a new ElementIterator for traversing the children of el.
func NewIterator(el Element) *ElementIterator {
	return &ElementIterator{
		root: el,
	}
}

// Next advances the iterator and returns the next child element, or nil if done.
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

package hl7v2

import "github.com/blushift-io/hl7v2/query"

// ElementType is an enum of valid element types
//
//go:generate enumer -type=ElementType --linecomment -output=element_enums.go
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
	Parent() Element
	Children() []Element
	Position() int
	Location() query.Location
	GetLocation(query.Location) (Element, error)
	Value() Value
}

func makeElements[T Element](els ...T) []Element {
	res := make([]Element, len(els))
	for i := range els {
		res[i] = els[i]
	}

	return res
}

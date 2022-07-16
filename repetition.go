package hl7v2

type Repetition struct {
	parent   Element
	children []*Component
	pos      int
}

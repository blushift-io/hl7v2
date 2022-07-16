package hl7v2

type Component struct {
	parent   Element
	children []*Subcomponent
	pos      int
}

func (el *Component) Type() ElementType {
	return ElementComponent
}

func (el *Component) Name() string {
	return ""
}

func (el *Component) Delimiters() *Delimiters {
	return el.parent.Delimiters()
}

func (el *Component) Parent() Element {
	return el.parent
}

func (el *Component) Children() []Element {
	return makeElements(el.children...)
}

func (el *Component) Position() any {
	return el.pos
}

func (el *Component) Value() Value {
	panic("not implemented") // TODO: Implement
}

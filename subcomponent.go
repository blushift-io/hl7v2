package hl7v2

type Subcomponent struct {
	v      Value
	pos    int
	parent Element
}

func (el *Subcomponent) Type() ElementType {
	return ElementSubcomponent
}

func (el *Subcomponent) Name() string {
	return ""
}

func (el *Subcomponent) Delimiters() *Delimiters {
	return el.parent.Delimiters()
}

func (el *Subcomponent) Parent() Element {
	return el.parent
}

func (el *Subcomponent) Children() []Element {
	return nil
}

func (el *Subcomponent) Position() any {
	return el.pos
}

func (el *Subcomponent) Value() Value {
	return el.v
}

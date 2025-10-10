package builder

import "github.com/blushift-io/hl7v2"

type SubcomponentBuilder struct {
	val []byte
}

func BuildSubcomponent(val hl7v2.Value) *SubcomponentBuilder {
	return &SubcomponentBuilder{
		val: val.Bytes(),
	}
}

func EmptySubcomponent() *SubcomponentBuilder {
	return &SubcomponentBuilder{
		val: []byte{},
	}
}

func (s *SubcomponentBuilder) SetValue(v hl7v2.Value) *SubcomponentBuilder {
	s.val = v.Bytes()

	return s
}

func (s *SubcomponentBuilder) Build() hl7v2.RawSubcomponent {
	return hl7v2.RawSubcomponent(s.val)
}

func (s *SubcomponentBuilder) BuildElement(parent hl7v2.Element, pos int) hl7v2.Element {
	return hl7v2.NewSubcomponent(parent, pos, hl7v2.NewValue(s.val))
}

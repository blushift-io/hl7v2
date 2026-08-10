package builder

import "github.com/blushift-io/hl7v2"

// SubcomponentBuilder builds HL7 v2 subcomponents.
type SubcomponentBuilder struct {
	val []byte
}

// BuildSubcomponent creates a SubcomponentBuilder initialized with the given value.
func BuildSubcomponent(val hl7v2.Value) *SubcomponentBuilder {
	return &SubcomponentBuilder{
		val: val.Bytes(),
	}
}

// EmptySubcomponent creates an empty SubcomponentBuilder.
func EmptySubcomponent() *SubcomponentBuilder {
	return &SubcomponentBuilder{
		val: []byte{},
	}
}

// SetValue sets the raw byte value of the subcomponent.
func (s *SubcomponentBuilder) SetValue(v hl7v2.Value) *SubcomponentBuilder {
	s.val = v.Bytes()

	return s
}

// Build constructs and returns the RawSubcomponent.
func (s *SubcomponentBuilder) Build() hl7v2.RawSubcomponent {
	return hl7v2.RawSubcomponent(s.val)
}

// BuildElement constructs an HL7 v2 Element representation of the subcomponent.
func (s *SubcomponentBuilder) BuildElement(parent hl7v2.Element, pos int) hl7v2.Element {
	return hl7v2.NewSubcomponent(parent, pos, hl7v2.NewValue(s.val))
}

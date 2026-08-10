package builder

import (
	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/query"
)

// ComponentBuilder builds HL7 v2 message components consisting of subcomponents.
type ComponentBuilder struct {
	subs []*SubcomponentBuilder
}

// BuildComponent creates a ComponentBuilder from a list of subcomponent builders.
func BuildComponent(subs ...*SubcomponentBuilder) *ComponentBuilder {
	return &ComponentBuilder{
		subs: subs,
	}
}

// SingleComponent creates a ComponentBuilder containing a single subcomponent with the given value.
func SingleComponent(val hl7v2.Value) *ComponentBuilder {
	return &ComponentBuilder{
		subs: []*SubcomponentBuilder{
			BuildSubcomponent(val),
		},
	}
}

// EmptyComponent creates an empty ComponentBuilder.
func EmptyComponent() *ComponentBuilder {
	return &ComponentBuilder{}
}

// Subcomponent appends a new empty SubcomponentBuilder and returns the ComponentBuilder.
func (c *ComponentBuilder) Subcomponent() *ComponentBuilder {
	s := &SubcomponentBuilder{}
	c.subs = append(c.subs, s)

	return c
}

// AddSubcomponent appends a SubcomponentBuilder to the ComponentBuilder.
func (c *ComponentBuilder) AddSubcomponent(s *SubcomponentBuilder) *ComponentBuilder {
	c.subs = append(c.subs, s)

	return c
}

// SetLocation sets a value at the specified subcomponent location.
func (c *ComponentBuilder) SetLocation(loc query.Location, v hl7v2.Value) *ComponentBuilder {
	if loc.Subcomponent < 0 {
		return c
	}

	size := len(c.subs)
	if size == 0 {
		if delta := loc.Subcomponent - size; delta == 0 {
			return c.AddSubcomponent(BuildSubcomponent(v))
		} else {
			for range delta {
				c.AddSubcomponent(EmptySubcomponent())
			}
		}
	} else if loc.Component > size {
		for range loc.Component - size {
			c.AddSubcomponent(EmptySubcomponent())
		}
	}

	c.subs[loc.Subcomponent].SetValue(v)

	return c
}

// Build constructs and returns the RawComponent.
func (c *ComponentBuilder) Build() hl7v2.RawComponent {
	rc := make(hl7v2.RawComponent, len(c.subs))
	for i, s := range c.subs {
		rc[i] = s.Build()
	}

	return rc
}

// BuildElement constructs an HL7 v2 Element representation of the component.
func (c *ComponentBuilder) BuildElement(parent hl7v2.Element, pos int) hl7v2.Element {
	return hl7v2.NewComponent(parent, pos, c.Build())
}

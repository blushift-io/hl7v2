package builder

import (
	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/query"
)

type ComponentBuilder struct {
	subs []*SubcomponentBuilder
}

func BuildComponent(subs ...*SubcomponentBuilder) *ComponentBuilder {
	return &ComponentBuilder{
		subs: subs,
	}
}

func SingleComponent(val hl7v2.Value) *ComponentBuilder {
	return &ComponentBuilder{
		subs: []*SubcomponentBuilder{
			BuildSubcomponent(val),
		},
	}
}

func EmptyComponent() *ComponentBuilder {
	return &ComponentBuilder{}
}

func (c *ComponentBuilder) Subcomponent() *ComponentBuilder {
	s := &SubcomponentBuilder{}
	c.subs = append(c.subs, s)

	return c
}

func (c *ComponentBuilder) AddSubcomponent(s *SubcomponentBuilder) *ComponentBuilder {
	c.subs = append(c.subs, s)

	return c
}

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

func (c *ComponentBuilder) Build() hl7v2.RawComponent {
	rc := make(hl7v2.RawComponent, len(c.subs))
	for i, s := range c.subs {
		rc[i] = s.Build()
	}

	return rc
}

func (c *ComponentBuilder) BuildElement(parent hl7v2.Element, pos int) hl7v2.Element {
	return hl7v2.NewComponent(parent, pos, c.Build())
}

package builder

import (
	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/query"
)

// RepetitionBuilder builds HL7 v2 field repetitions containing components.
type RepetitionBuilder struct {
	comps []*ComponentBuilder
}

// BuildRepetition creates a RepetitionBuilder initialized with component builders.
func BuildRepetition(comps ...*ComponentBuilder) *RepetitionBuilder {
	return &RepetitionBuilder{
		comps: comps,
	}
}

// EmptyRepetition creates an empty RepetitionBuilder.
func EmptyRepetition() *RepetitionBuilder {
	return &RepetitionBuilder{}
}

// Component appends a new empty ComponentBuilder and returns it.
func (r *RepetitionBuilder) Component() *ComponentBuilder {
	c := &ComponentBuilder{}
	r.comps = append(r.comps, c)

	return c
}

// AddComponent appends a ComponentBuilder to the RepetitionBuilder.
func (r *RepetitionBuilder) AddComponent(c *ComponentBuilder) *RepetitionBuilder {
	r.comps = append(r.comps, c)

	return r
}

// SetLocation sets a value at the specified component location.
func (r *RepetitionBuilder) SetLocation(loc query.Location, val hl7v2.Value) *RepetitionBuilder {
	if loc.Component < 0 {
		return r
	}

	size := len(r.comps)
	if size == 0 {
		if delta := loc.Component - size; delta == 0 {
			return r.AddComponent(EmptyComponent()).SetLocation(loc, val)
		} else {
			for range delta {
				r.AddComponent(EmptyComponent())
			}
		}
	} else if loc.Component > size {
		delta := loc.Component - size
		for range delta {
			r.AddComponent(EmptyComponent())
		}
	}
	r.comps[loc.Component].SetLocation(loc, val)

	return r
}

// Build constructs and returns the RawRepetition.
func (r *RepetitionBuilder) Build() hl7v2.RawRepetition {
	rr := make(hl7v2.RawRepetition, len(r.comps))
	for i, c := range r.comps {
		rr[i] = c.Build()
	}
	return rr
}

// BuildElement constructs an HL7 v2 Element representation of the repetition.
func (r *RepetitionBuilder) BuildElement(parent hl7v2.Element, pos int) hl7v2.Element {
	return hl7v2.NewRepetition(parent, pos, r.Build())
}

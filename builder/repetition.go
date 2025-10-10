package builder

import (
	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/query"
)

type RepetitionBuilder struct {
	comps []*ComponentBuilder
}

func BuildRepetition(comps ...*ComponentBuilder) *RepetitionBuilder {
	return &RepetitionBuilder{
		comps: comps,
	}
}

func EmptyRepetition() *RepetitionBuilder {
	return &RepetitionBuilder{}
}

func (r *RepetitionBuilder) Component() *ComponentBuilder {
	c := &ComponentBuilder{}
	r.comps = append(r.comps, c)

	return c
}

func (r *RepetitionBuilder) AddComponent(c *ComponentBuilder) *RepetitionBuilder {
	r.comps = append(r.comps, c)

	return r
}

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

func (r *RepetitionBuilder) Build() hl7v2.RawRepetition {
	rr := make(hl7v2.RawRepetition, len(r.comps))
	for i, c := range r.comps {
		rr[i] = c.Build()
	}
	return rr
}

func (r *RepetitionBuilder) BuildElement(parent hl7v2.Element, pos int) hl7v2.Element {
	return hl7v2.NewRepetition(parent, pos, r.Build())
}

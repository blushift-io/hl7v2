package builder

import (
	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/query"
)

type FieldBuilder struct {
	reps []*RepetitionBuilder
}

func BuildField(reps ...*RepetitionBuilder) *FieldBuilder {
	return &FieldBuilder{
		reps: reps,
	}
}

func SingleValueField(val hl7v2.Value) *FieldBuilder {
	return &FieldBuilder{
		reps: []*RepetitionBuilder{
			BuildRepetition(BuildComponent(BuildSubcomponent(val))),
		},
	}
}

func RepeatingField(vals ...hl7v2.Value) *FieldBuilder {
	reps := make([]*RepetitionBuilder, len(vals))

	for i, v := range vals {
		reps[i] = BuildRepetition(
			BuildComponent(
				BuildSubcomponent(v),
			),
		)

	}

	return BuildField(reps...)
}

func ComponentField(vals ...hl7v2.Value) *FieldBuilder {
	comps := make([]*ComponentBuilder, len(vals))

	for i, v := range vals {
		comps[i] = BuildComponent(
			BuildSubcomponent(v),
		)
	}

	return BuildField(BuildRepetition(comps...))
}

func EmptyField() *FieldBuilder {
	return &FieldBuilder{}
}

func (f *FieldBuilder) Repetition() *RepetitionBuilder {
	r := &RepetitionBuilder{}
	f.reps = append(f.reps, r)

	return r
}

func (f *FieldBuilder) AddRepetition(r *RepetitionBuilder) *FieldBuilder {
	f.reps = append(f.reps, r)

	return f
}

func (f *FieldBuilder) SetLocation(loc query.Location, val hl7v2.Value) *FieldBuilder {
	rep := 0
	if loc.FieldRep != nil {
		rep = *loc.FieldRep
	}

	size := len(f.reps)
	if size == 0 {
		if delta := rep - size; delta == 0 {
			return f.AddRepetition(EmptyRepetition()).SetLocation(loc, val)
		} else {
			for range delta {
				f.AddRepetition(EmptyRepetition())
			}
		}
	} else if rep > size {
		for range rep - size {
			f.AddRepetition(EmptyRepetition())
		}
	}

	f.reps[rep].SetLocation(loc, val)

	return f
}

func (f *FieldBuilder) Build() hl7v2.RawField {
	rf := make(hl7v2.RawField, len(f.reps))
	for i, r := range f.reps {
		rf[i] = r.Build()
	}

	return rf
}

func (f *FieldBuilder) BuildElement(parent hl7v2.Element, pos int) hl7v2.Element {
	return hl7v2.NewField(parent, pos, f.Build())
}

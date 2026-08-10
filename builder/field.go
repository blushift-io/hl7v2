package builder

import (
	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/query"
)

// FieldBuilder builds HL7 v2 message fields consisting of repetitions.
type FieldBuilder struct {
	reps []*RepetitionBuilder
}

// BuildField creates a FieldBuilder initialized with the given repetition builders.
func BuildField(reps ...*RepetitionBuilder) *FieldBuilder {
	return &FieldBuilder{
		reps: reps,
	}
}

// SingleValueField creates a FieldBuilder containing a single field value.
func SingleValueField(val hl7v2.Value) *FieldBuilder {
	return &FieldBuilder{
		reps: []*RepetitionBuilder{
			BuildRepetition(BuildComponent(BuildSubcomponent(val))),
		},
	}
}

// RepeatingField creates a FieldBuilder with multiple repeating values.
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

// ComponentField creates a FieldBuilder where each value is placed into a component.
func ComponentField(vals ...hl7v2.Value) *FieldBuilder {
	comps := make([]*ComponentBuilder, len(vals))

	for i, v := range vals {
		comps[i] = BuildComponent(
			BuildSubcomponent(v),
		)
	}

	return BuildField(BuildRepetition(comps...))
}

// EmptyField creates an empty FieldBuilder.
func EmptyField() *FieldBuilder {
	return &FieldBuilder{}
}

// Repetition appends a new empty RepetitionBuilder and returns it.
func (f *FieldBuilder) Repetition() *RepetitionBuilder {
	r := &RepetitionBuilder{}
	f.reps = append(f.reps, r)

	return r
}

// AddRepetition appends a RepetitionBuilder to the FieldBuilder.
func (f *FieldBuilder) AddRepetition(r *RepetitionBuilder) *FieldBuilder {
	f.reps = append(f.reps, r)

	return f
}

// SetLocation sets a value at the specified field repetition location.
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

// Build constructs and returns the RawField.
func (f *FieldBuilder) Build() hl7v2.RawField {
	rf := make(hl7v2.RawField, len(f.reps))
	for i, r := range f.reps {
		rf[i] = r.Build()
	}

	return rf
}

// BuildElement constructs an HL7 v2 Element representation of the field.
func (f *FieldBuilder) BuildElement(parent hl7v2.Element, pos int) hl7v2.Element {
	return hl7v2.NewField(parent, pos, f.Build())
}

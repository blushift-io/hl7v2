package schema

//go:generate enumer -type=Optionality,Repeatability -json -text -yaml -linecomment -output=schema_enums.go
type Optionality int

const (
	OptionalityOptional    Optionality = iota //O
	OptionalityRequired                       //R
	OptionalityConditional                    //C
)

type Repeatability int

const (
	RepeatabilitySingle   Repeatability = iota //-
	RepeatabilityFixed                         //F
	RepeatabilityInfinite                      //*
)

type ValueSchema struct {
	Name          string        `json:"name"`
	Length        int           `json:"length"`
	Optionality   Optionality   `json:"optionality"`
	Repeatability Repeatability `json:"repeatability"`
	MaxRepeats    int           `json:"max_repeats"`
}

type ValueSchemaOption func(*ValueSchema)

func NewValueSchema(opts ...ValueSchemaOption) ValueSchema {
	vs := ValueSchema{}

	for _, opt := range opts {
		opt(&vs)
	}

	return vs
}

func WithName(name string) ValueSchemaOption {
	return func(vs *ValueSchema) {
		vs.Name = name
	}
}

func WithLength(length int) ValueSchemaOption {
	return func(vs *ValueSchema) {
		vs.Length = length
	}
}

func WithOptionality(optionality Optionality) ValueSchemaOption {
	return func(vs *ValueSchema) {
		vs.Optionality = optionality
	}
}

func WithRepeatability(repeatability Repeatability) ValueSchemaOption {
	return func(vs *ValueSchema) {
		vs.Repeatability = repeatability
	}
}

func WithMaxRepeats(maxRepeats int) ValueSchemaOption {
	return func(vs *ValueSchema) {
		vs.MaxRepeats = maxRepeats
	}
}

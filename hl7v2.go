package hl7v2

//TODO: Escape/Unescape HL7 delimiters
//TODO: Implement MCF delayed acknowledgment
//TODO: Implement batch and file elements

type Marshaler interface {
	MarshalHL7() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalHL7(b []byte) error
}

func Marshal(v any) ([]byte, error) {
	return marshalHL7(v)
}

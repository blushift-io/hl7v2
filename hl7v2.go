// Package hl7v2 provides tools for parsing, building, manipulating, and marshaling HL7 v2 messages.
package hl7v2

//TODO: Escape/Unescape HL7 delimiters
//TODO: Implement MCF delayed acknowledgment
//TODO: Implement batch and file elements

// Marshaler is the interface implemented by types that can marshal themselves into valid HL7 v2 wire bytes.
type Marshaler interface {
	MarshalHL7() ([]byte, error)
}

// Unmarshaler is the interface implemented by types that can unmarshal an HL7 v2 wire byte representation of themselves.
type Unmarshaler interface {
	UnmarshalHL7(b []byte) error
}

// Marshal returns the HL7 v2 wire byte representation of v.
func Marshal(v any) ([]byte, error) {
	return marshalHL7(v)
}


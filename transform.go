package hl7v2

import "bytes"

// RawTransform is a function type for transforming raw HL7 byte slices.
type RawTransform func([]byte) ([]byte, error)

func fixLineEndings(b []byte) ([]byte, error) {
	return ReplaceLineEndings(b), nil
}

// ReplaceLineEndings replaces CRLF (\r\n) and LF (\n) with CR (\r) segment delimiters.
func ReplaceLineEndings(b []byte) []byte {
	b = bytes.ReplaceAll(b, []byte("\r\n"), []byte("\r"))
	return bytes.ReplaceAll(b, []byte{'\n'}, []byte{'\r'})
}

// ValueTransform is a function type for extracting or transforming a Value from an Element.
type ValueTransform func(el Element) (Value, error)

func noopLocationTransform(el Element) (Value, error) {
	return el.Value(), nil
}

// UnescapeElement extracts the unescaped Value from an Element using its delimiters.
func UnescapeElement(el Element) (Value, error) {
	return el.Value().Escape(el.Delimiters()), nil
}

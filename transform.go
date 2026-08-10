package hl7v2

import "bytes"

type RawTransform func([]byte) ([]byte, error)

func fixLineEndings(b []byte) ([]byte, error) {
	return ReplaceLineEndings(b), nil
}

func ReplaceLineEndings(b []byte) []byte {
	b = bytes.ReplaceAll(b, []byte("\r\n"), []byte("\r"))
	return bytes.ReplaceAll(b, []byte{'\n'}, []byte{'\r'})
}

type ValueTransform func(el Element) (Value, error)

func noopLocationTransform(el Element) (Value, error) {
	return el.Value(), nil
}

func UnescapeElement(el Element) (Value, error) {
	return el.Value().Escape(el.Delimiters()), nil
}

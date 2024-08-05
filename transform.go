package hl7v2

import "bytes"

type RawTransform func([]byte) ([]byte, error)

func fixLineEndings(b []byte) ([]byte, error) {
	return replaceLineEndings(b), nil
}

func replaceLineEndings(b []byte) []byte {
	b = bytes.ReplaceAll(b, []byte("\r\n"), []byte("\r"))
	return bytes.ReplaceAll(b, []byte{'\n'}, []byte{'\r'})
}

package mllp

import (
	"fmt"
	"io"
)

const (
	startByte      = byte(0x0b)
	endByte        = byte(0x1c)
	carriageReturn = byte(0x0d)
)

type ReadWriter struct {
	Reader *Reader
	Writer *Writer
}

func NewReadWriter(rw io.ReadWriter) *ReadWriter {
	return &ReadWriter{
		Reader: NewReader(rw),
		Writer: NewWriter(rw),
	}
}

func Encode(b []byte) []byte {
	return append(append([]byte{startByte}, b...), []byte{endByte, carriageReturn}...)
}

func Decode(b []byte) ([]byte, error) {
	if len(b) < 3 {
		return nil, fmt.Errorf("unwrap mllp: input too short")
	}

	if err := checkByte(b, 0, startByte); err != nil {
		return nil, fmt.Errorf("unwrap mllp: %w", err)
	}

	if err := checkByte(b, len(b)-2, endByte); err != nil {
		return nil, fmt.Errorf("unwrap mllp: %w", err)
	}

	if err := checkByte(b, len(b)-1, carriageReturn); err != nil {
		return nil, fmt.Errorf("unwrap mllp: %w", err)
	}

	return b[1 : len(b)-2], nil
}

func checkByte(b []byte, pos int, exp byte) error {
	if v := b[pos]; v != exp {
		return fmt.Errorf("mllp: expected byte %v, got %v", exp, v)
	}

	return nil
}

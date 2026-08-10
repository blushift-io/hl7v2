package mllp

import (
	"bufio"
	"fmt"
	"io"
)

// Read reads and decodes a single MLLP framed message from an io.Reader.
func Read(r io.Reader) ([]byte, error) {
	var br *bufio.Reader
	if v, ok := r.(*bufio.Reader); ok {
		br = v
	} else {
		br = bufio.NewReader(r)
	}

	b, err := br.ReadBytes(endByte)
	if err != nil {
		return nil, fmt.Errorf("mllp: %w", err)
	}

	lb, err := br.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("mllp: %w", err)
	}

	msg, err := Decode(append(b, lb))
	if err != nil {
		return nil, fmt.Errorf("mllp: %w", err)
	}

	return msg, nil
}

// Reader reads MLLP framed messages from a buffered input stream.
type Reader struct {
	b *bufio.Reader
}

// NewReader creates a Reader wrapping an io.Reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		b: bufio.NewReader(r),
	}
}

// ReadMessage reads the next MLLP message frame from the input stream.
func (r *Reader) ReadMessage() ([]byte, error) {
	var (
		err error
		b   []byte
	)

	for {
		b, err = r.b.ReadBytes(endByte)
		if err != nil {
			return nil, err
		}

		if len(b) > 1 && b[len(b)-2] == carriageReturn {
			break
		}
	}

	return b[:len(b)-1], nil
}

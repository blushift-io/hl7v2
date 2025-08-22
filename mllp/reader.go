package mllp

import (
	"bufio"
	"fmt"
	"io"
)

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

type Reader struct {
	b *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		b: bufio.NewReader(r),
	}
}

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

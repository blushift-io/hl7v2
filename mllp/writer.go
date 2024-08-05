package mllp

import (
	"fmt"
	"io"
)

func Write(w io.Writer, b []byte) error {
	if _, err := w.Write(Wrap(b)); err != nil {
		return fmt.Errorf("mllp: %w", err)
	}

	return nil
}

type Writer struct {
	w io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: w,
	}
}

func (w *Writer) Write(b []byte) error {
	return Write(w.w, b)
}

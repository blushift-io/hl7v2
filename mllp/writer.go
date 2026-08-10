package mllp

import (
	"fmt"
	"io"
)

// Write encodes payload bytes into MLLP framing and writes them to an io.Writer.
func Write(w io.Writer, b []byte) error {
	if _, err := w.Write(Encode(b)); err != nil {
		return fmt.Errorf("mllp: %w", err)
	}

	return nil
}

// Writer writes MLLP framed messages to an underlying io.Writer.
type Writer struct {
	w io.Writer
}

// NewWriter creates a Writer wrapping an io.Writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: w,
	}
}

// Write encodes payload bytes with MLLP framing and writes them to the underlying writer.
func (w *Writer) Write(b []byte) error {
	return Write(w.w, b)
}

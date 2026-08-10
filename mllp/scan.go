package mllp

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

// Scan is a bufio.SplitFunc for scanning MLLP framed messages from a byte stream.
func Scan(b []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(b) == 0 {
		return 0, nil, nil
	}

	sep := []byte{endByte, carriageReturn}
	if i := bytes.Index(b, sep); i > 0 {
		bb := b[:i+len(sep)]
		adv := len(bb)

		msg, err := Read(bytes.NewReader(bb))
		return adv, msg, err
	}

	if atEOF {
		return len(b), b, nil
	}

	return 0, nil, nil
}

// NewScanner creates a bufio.Scanner configured with the MLLP split function.
func NewScanner(r io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Split(Scan)

	return scanner
}

// NewBufferedScanner creates a bufio.Scanner with a custom buffer size and MLLP split function.
func NewBufferedScanner(r io.Reader, size int) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, size), size)
	scanner.Split(Scan)

	return scanner
}

// NewFileScanner creates a bufio.Scanner optimized for reading MLLP messages from an os.File.
func NewFileScanner(f *os.File) (*bufio.Scanner, error) {
	size := bufio.MaxScanTokenSize
	i, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if i.Size() < int64(size) {
		size = int(i.Size())
	}

	return NewBufferedScanner(f, size), nil
}

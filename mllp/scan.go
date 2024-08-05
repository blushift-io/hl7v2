package mllp

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

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

func NewScanner(r io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Split(Scan)

	return scanner
}

func NewBufferedScanner(r io.Reader, size int) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, size), size)
	scanner.Split(Scan)

	return scanner
}

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

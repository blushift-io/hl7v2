package hl7v2

import (
	"fmt"
	"os"
	"testing"
)

func TestParseRawMessage(t *testing.T) {
	b, err := os.ReadFile("./fixtures/ORU_A01_pdf.hl7")
	if err != nil {
		t.Fatal(err)
	}

	m, err := NewRawMessageFromBytes(b, FixLineEndings())
	if err != nil {
		t.Fatal(err)
	}

	j, err := m.JSON(true)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(j))
}

func TestQueryRawMessage(t *testing.T) {
	b, err := os.ReadFile("./fixtures/ADT_A01_2.hl7")
	if err != nil {
		t.Fatal(err)
	}

	m, err := NewRawMessageFromBytes(b, FixLineEndings())
	if err != nil {
		t.Fatal(err)
	}

	v, err := m.QueryValue("PV1.3[1]")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(v.String())
}

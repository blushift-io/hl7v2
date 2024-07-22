package hl7v2

import (
	"fmt"
	"os"
	"testing"

	"github.com/blushift-io/hl7v2/query"
	"github.com/davecgh/go-spew/spew"
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

	q, err := query.ParseLocation("PV1.3[1]")
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(q)

	v, err := m.Query(q)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(v.String())
}

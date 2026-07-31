package hl7v2

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestSegmentElement(t *testing.T) {
	msg, err := NewMessageFromFile("./test/fixtures/ORM_O01_1.hl7", FixLineEndings())
	if err != nil {
		t.Fatal(err)
	}

	sel, err := msg.Select("OBR")
	if err != nil {
		t.Fatal(err)
	}

	if len(sel) < 1 {
		t.Fatalf("expected at least one segment, got %d", len(sel))
	}

	for _, seg := range sel {
		obr, ok := seg.(*Segment)
		if !ok {
			t.Fatalf("expected *Segment, got %T", seg)
		}

		spew.Dump(obr.Index())
	}
}

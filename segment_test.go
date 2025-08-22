package hl7v2

import (
	"os"
	"testing"

	"github.com/blushift-io/hl7v2/query"
	"github.com/davecgh/go-spew/spew"
)

func TestSegmentElement(t *testing.T) {
	msg, err := NewMessageFromFile("./fixtures/ORM_O01_1.hl7", FixLineEndings())
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

func TestSegmentBuilder(t *testing.T) {
	msg, err := NewMessageFromFile("./fixtures/ORM_O01_1.hl7", FixLineEndings())
	if err != nil {
		t.Fatal(err)
	}

	zseg1 := BuildSegment("ZXX").
		AddField(SingleValueField(NewStringValue("1"))).
		AddField(SingleValueField(NewStringValue("2"))).
		AddField(SingleValueField(NewStringValue("3")))

	zseg2 := BuildSegment("ZXY").
		AddField(SingleValueField(NewStringValue("1"))).
		AddField(SingleValueField(NewStringValue("2"))).
		AddField(SingleValueField(NewStringValue("3")))

	if err := msg.Append(zseg1.BuildElement(msg, 0)); err != nil {
		t.Fatal(err)
	}

	if err := msg.Append(zseg2.BuildElement(msg, 0)); err != nil {
		t.Fatal(err)
	}

	sel, err := msg.Select("ZXX")
	if err != nil {
		t.Fatal(err)
	}

	if len(sel) < 1 {
		t.Fatalf("expected at least one segment, got %d", len(sel))
	}

	loc := query.Location{
		Segment: "ZXX",
		Field:   2,
	}

	if err := msg.SetLocation(loc, NewStringValue("4")); err != nil {
		t.Fatal(err)
	}

	enc, err := msg.Encode()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile("./tmp/appended_zxx.hl7", enc, 0644); err != nil {
		t.Fatal(err)
	}

}

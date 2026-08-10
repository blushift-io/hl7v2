package builder

import (
	"fmt"
	"os"
	"testing"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/query"
	_ "github.com/blushift-io/hl7v2/schema/spec/v251"
	"github.com/blushift-io/hl7v2/builder/field"
)

func TestBuilder(t *testing.T) {
	mt := hl7v2.MessageType{
		Code:      "ADT",
		Event:     "A01",
		Structure: "ADT_A01",
	}

	b := New()

	b.Header(mt, hl7v2.Version251,
		SendingApp("test_app"),
		SendingFacility("test_facility"),
		ReceivingApp("test_receiver"),
		ReceivingFacility("test_receiver_facility"),
		ProcessingID("abcd123"),
		func(hb *HeaderBuilder) {
			hb.SetSequenceNumber(1)
		},
	)

	b.Segment("PID", WithFields(
		SingleValueField(hl7v2.NewStringValue("1")),
		SingleValueField(hl7v2.NewStringValue("Doe")),
		SingleValueField(hl7v2.NewStringValue("John")),
	))

	msg, err := b.Build()
	if err != nil {
		t.Fatal(err)
	}

	if len(msg.Segments()) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(msg.Segments()))
	}

	if !msg.HasSegment("MSH") {
		t.Fatal("expected MSH segment to be present")
	}

	out, err := msg.Encode()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(out))
}

func TestSegmentBuilder(t *testing.T) {
	msg, err := hl7v2.NewMessageFromFile("../test/fixtures/ORM_O01_1.hl7", hl7v2.FixLineEndings())
	if err != nil {
		t.Fatal(err)
	}

	zseg1 := Segment("ZXX").
		AddField(SingleValueField(hl7v2.NewStringValue("1"))).
		AddField(SingleValueField(hl7v2.NewStringValue("2"))).
		AddField(SingleValueField(hl7v2.NewStringValue("3")))

	zseg2 := Segment("ZXY").
		AddField(SingleValueField(hl7v2.NewStringValue("1"))).
		AddField(SingleValueField(hl7v2.NewStringValue("2"))).
		AddField(SingleValueField(hl7v2.NewStringValue("3")))

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

	if err := msg.SetLocation(loc, hl7v2.NewStringValue("4")); err != nil {
		t.Fatal(err)
	}

	enc, err := msg.Encode()
	if err != nil {
		t.Fatal(err)
	}

	_ = os.MkdirAll("./tmp", 0755)
	if err := os.WriteFile("./tmp/appended_zxx.hl7", enc, 0644); err != nil {
		t.Fatal(err)
	}
}


func TestFullBuilderFlow(t *testing.T) {
	mt := hl7v2.MessageType{Code: "ADT", Event: "A01", Structure: "ADT_A01"}

	b := New().
		Header(mt, hl7v2.Version251,
			SendingApp("EPIC"),
			SendingFacility("HOSPITAL_A"),
		)
		
	b.Segment("PID",
		field.Int(1),
		field.String(""),
		field.Components("12345", "", "", "MRN"),
		field.Components("Doe", "John"),
	)

	b.Set("PV1-1", 1).
		Set("PV1-2", "I").
		Set("PV1-3.1", "ROOM1").
		Set("PV1-3.2", "BED2")

	msg, err := b.Build()
	if err != nil {
		t.Fatalf("build failed: %v", err)
	}

	if len(msg.Segments()) != 3 { // MSH, PID, PV1
		t.Fatalf("expected 3 segments, got %d", len(msg.Segments()))
	}
	pv1 := b.GetSegment("PV1").Build()
	if string(pv1[3][0][0][0]) != "ROOM1" || string(pv1[3][0][1][0]) != "BED2" {
		t.Fatalf("expected PV1-3 to be ROOM1^BED2, got %v", pv1[3][0])
	}
}

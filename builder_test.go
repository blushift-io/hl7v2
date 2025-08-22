package hl7v2

import (
	"fmt"
	"testing"

	_ "github.com/blushift-io/hl7v2/schema/spec/v251"
)

func TestBuilder(t *testing.T) {
	mt := MessageType{
		Code:      "ADT",
		Event:     "A01",
		Structure: "ADT_A01",
	}

	b := NewBuilder()

	b.Header(mt, Version251.String()).
		SetSendingApplication("test_app").
		SetSendingFacility("test_facility").
		SetReceivingApplication("test_receiver").
		SetReceivingFacility("test_receiver_facility").
		SetProcessingID("abcd123").
		SetSequenceNumber(1)

	b.Segment("PID",
		SingleValueField(NewStringValue("1")),
		SingleValueField(NewStringValue("Doe")),
		SingleValueField(NewStringValue("John")),
	)

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

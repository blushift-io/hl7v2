package hl7v2

import (
	"bytes"
	"testing"
)

type sampleComponent struct {
	ID   string `hl7:"1"`
	Name string `hl7:"2"`
}

type sampleSegment struct {
	Field1 string            `hl7:"1"`
	Field2 int               `hl7:"2"`
	Field3 float64           `hl7:"3"`
	Comp   sampleComponent   `hl7:"4"`
	Reps   []sampleComponent `hl7:"5"`
}

type sampleMessage struct {
	MSH sampleSegment `hl7:"MSH"`
	PID sampleSegment `hl7:"PID"`
}

func TestMarshalSegment(t *testing.T) {
	seg := sampleSegment{
		Field1: "test",
		Field2: 123,
		Field3: 456.78,
		Comp: sampleComponent{
			ID:   "ID1",
			Name: "Name1",
		},
		Reps: []sampleComponent{
			{ID: "R1", Name: "N1"},
			{ID: "R2", Name: "N2"},
		},
	}

	b, err := Marshal(seg)
	if err != nil {
		t.Fatalf("failed to marshal segment: %v", err)
	}

	expected := "sampleSegment|test|123|456.78|ID1^Name1|R1^N1~R2^N2\r"
	if string(b) != expected {
		t.Errorf("got %q, want %q", string(b), expected)
	}
}

func TestMarshalMessage(t *testing.T) {
	msg := sampleMessage{
		MSH: sampleSegment{
			Field1: "|",
			Field2: 2,
		},
		PID: sampleSegment{
			Field1: "PID1",
			Field2: 100,
		},
	}

	b, err := Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	if !bytes.Contains(b, []byte("MSH|")) || !bytes.Contains(b, []byte("PID|PID1|100")) {
		t.Errorf("unexpected marshaled message output: %q", string(b))
	}
}

func TestUnmarshal(t *testing.T) {
	er7 := []byte("MSH|^~\\&|SEND_APP|SEND_FAC\rPID|1|12345|Doe^John\r")

	type mshStruct struct {
		FieldSeparator string `hl7:"MSH.1"`
		SendingApp     string `hl7:"MSH.3"`
		SendingFac     string `hl7:"MSH.4"`
	}

	type pidName struct {
		Family string `hl7:"1"`
		Given  string `hl7:"2"`
	}

	type pidStruct struct {
		SetID string  `hl7:"PID.1"`
		ID    int     `hl7:"PID.2"`
		Name  pidName `hl7:"PID.3"`
	}

	type testMsg struct {
		MSH mshStruct `hl7:"MSH"`
		PID pidStruct `hl7:"PID"`
	}

	var msg testMsg
	if err := Unmarshal(er7, &msg); err != nil {
		t.Fatalf("failed to unmarshal ER7: %v", err)
	}

	if msg.MSH.SendingApp != "SEND_APP" {
		t.Errorf("got SendingApp %q, want %q", msg.MSH.SendingApp, "SEND_APP")
	}

	if msg.PID.SetID != "1" {
		t.Errorf("got SetID %q, want %q", msg.PID.SetID, "1")
	}

	if msg.PID.ID != 12345 {
		t.Errorf("got ID %d, want %d", msg.PID.ID, 12345)
	}

	if msg.PID.Name.Family != "Doe" || msg.PID.Name.Given != "John" {
		t.Errorf("got Name %s^%s, want Doe^John", msg.PID.Name.Family, msg.PID.Name.Given)
	}
}

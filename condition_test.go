package hl7v2

import "testing"

func TestHasCondition(t *testing.T) {
	msg, err := NewMessageFromFile("./test/fixtures/ADT_A01_1.hl7", FixLineEndings())
	if err != nil {
		t.Fatal(err)
	}

	cond, err := Has("MSH.4")
	if err != nil {
		t.Fatal(err)
	}

	if !cond.Check(msg) {
		t.Errorf("expected condition to be true")
	}
}

func TestValueContainsCondition(t *testing.T) {
	msg, err := NewMessageFromFile("./test/fixtures/ADT_A01_1.hl7", FixLineEndings())
	if err != nil {
		t.Fatal(err)
	}

	if cond, err := ValueContains("MSH.6", "GOOD HEALTH HOSPITAL"); err == nil {
		if !cond.Check(msg) {
			t.Errorf("expected condition to be true")
		}
	} else {
		t.Fatal(err)
	}

	if cond, err := ValueContains("PID.5.1", "EVERYMAN"); err == nil {
		if !cond.Check(msg) {
			t.Errorf("expected condition to be true")
		}
	} else {
		t.Fatal(err)
	}
}

func TestValueContainsAnyCondition(t *testing.T) {
	msg, err := NewMessageFromFile("./test/fixtures/ADT_A01_1.hl7", FixLineEndings())
	if err != nil {
		t.Fatal(err)
	}

	if cond, err := ValueContainsAny("MSH.6", "GOOD HEALTH HOSPITAL", "EKG"); err == nil {
		if !cond.Check(msg) {
			t.Errorf("expected condition to be true")
		}
	} else {
		t.Fatal(err)
	}

	if cond, err := ValueContainsAny("PID.5.1", "EVERYMAN", "SMITH"); err == nil {
		if !cond.Check(msg) {
			t.Errorf("expected condition to be true")
		}
	} else {
		t.Fatal(err)
	}
}

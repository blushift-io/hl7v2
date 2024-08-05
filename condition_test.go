package hl7v2

import "testing"

func TestHasCondition(t *testing.T) {
	msg, err := NewMessageFromFile("./fixtures/gcp_example.hl7", FixLineEndings())
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
	msg, err := NewMessageFromFile("./fixtures/gcp_example.hl7", FixLineEndings())
	if err != nil {
		t.Fatal(err)
	}

	if cond, err := ValueContains("MSH.6", "EHR"); err == nil {
		if !cond.Check(msg) {
			t.Errorf("expected condition to be true")
		}
	} else {
		t.Fatal(err)
	}

	if cond, err := ValueContains("OBX[3].3", "MDC_TEMP"); err == nil {
		if !cond.Check(msg) {
			t.Errorf("expected condition to be true")
		}
	} else {
		t.Fatal(err)
	}
}

func TestValueContainsAnyCondition(t *testing.T) {
	msg, err := NewMessageFromFile("./fixtures/gcp_example.hl7", FixLineEndings())
	if err != nil {
		t.Fatal(err)
	}

	if cond, err := ValueContainsAny("MSH.6", "EHR", "EKG"); err == nil {
		if !cond.Check(msg) {
			t.Errorf("expected condition to be true")
		}
	} else {
		t.Fatal(err)
	}

	if cond, err := ValueContainsAny("OBX[2].3", "MDC_TEMP", "MDC_PULS"); err == nil {
		if !cond.Check(msg) {
			t.Errorf("expected condition to be true")
		}
	} else {
		t.Fatal(err)
	}
}

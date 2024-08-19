package hl7v2

import (
	"testing"

	_ "github.com/blushift-io/hl7v2/schema/spec/v28"
)

func TestValidateMessageSchema(t *testing.T) {
	msg, err := NewMessageFromFile("./test/fixtures/ADT_A01_1.hl7", FixLineEndings())
	if err != nil {
		t.Fatal(err)
	}

	v := ValidateMessageSchema(msg)
	if !v.Valid {
		t.Errorf("expected valid message, got invalid: %v", v.Errors)
	}
}

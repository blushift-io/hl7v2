package hl7v2

import "testing"

func TestMarshal(t *testing.T) {
	type testStruct struct {
		A string  `hl7:"ZTT.1"`
		B int     `hl7:"ZTT.2"`
		C float64 `hl7:"ZTT.3"`
	}

	b, err := Marshal(&testStruct{
		A: "test",
		B: 123,
		C: 456.78,
	})
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	t.Logf("marshaled bytes: %v", string(b))
}

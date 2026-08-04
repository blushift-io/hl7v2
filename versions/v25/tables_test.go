package v25_test

import (
	"testing"

	"github.com/blushift-io/hl7v2/versions/v25"
)

func TestTables(t *testing.T) {
	if string(v25.Table0001Female) != "F" {
		t.Errorf("got %s, want F", v25.Table0001Female)
	}

	if v25.Table0001Female.String() != "F" {
		t.Errorf("got %s, want F", v25.Table0001Female.String())
	}

	if v25.Table0001Female.Description() != "Female" {
		t.Errorf("got Description %s, want Female", v25.Table0001Female.Description())
	}

	if !v25.Table0001Female.IsValid() {
		t.Errorf("expected Table0001Female to be valid")
	}

	// Test datatype conversions
	isVal := v25.Table0001Female.IS()
	if string(*isVal) != "F" {
		t.Errorf("got IS %s, want F", *isVal)
	}

	idVal := v25.Table0001Female.ID()
	if string(*idVal) != "F" {
		t.Errorf("got ID %s, want F", *idVal)
	}

	stVal := v25.Table0001Female.ST()
	if string(*stVal) != "F" {
		t.Errorf("got ST %s, want F", *stVal)
	}

	ceVal := v25.Table0001Female.CE()
	if string(ceVal.Identifier) != "F" || string(ceVal.Text) != "Female" {
		t.Errorf("got CE Identifier=%s Text=%s, want F/Female", ceVal.Identifier, ceVal.Text)
	}

	cweVal := v25.Table0001Female.CWE()
	if string(cweVal.Identifier) != "F" || string(cweVal.Text) != "Female" {
		t.Errorf("got CWE Identifier=%s Text=%s, want F/Female", cweVal.Identifier, cweVal.Text)
	}

	invalidSex := v25.Table0001("INVALID")
	if invalidSex.IsValid() {
		t.Errorf("expected invalidSex to be invalid")
	}

	if invalidSex.Description() != "" {
		t.Errorf("expected empty description for invalidSex")
	}
}

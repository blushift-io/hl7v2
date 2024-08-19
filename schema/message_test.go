package schema_test

import (
	"fmt"
	"testing"

	"github.com/blushift-io/hl7v2/schema"
	_ "github.com/blushift-io/hl7v2/schema/spec/v28"
)

func TestSchema(t *testing.T) {
	sch := schema.Open("2.8")
	if sch == nil {
		t.Fatal("schema not found")
	}

	adt := sch.Message("ADT_A01")
	if adt == nil {
		t.Fatal("message not found")
	}

	grp := adt.Segment("PROCEDURE")
	if grp == nil {
		t.Fatal("segment group not found")
	}

	g := adt.Grammar()
	if len(g) == 0 {
		t.Fatal("invalid message grammar")
	}

	fmt.Println(g.String())
}

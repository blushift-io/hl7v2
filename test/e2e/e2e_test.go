package e2e

import (
	"testing"

	"github.com/blushift-io/hl7v2"
	_ "github.com/blushift-io/hl7v2/schema/spec/v28"
)

func TestSchemaMessage(t *testing.T) {
	msg, err := hl7v2.NewMessageFromFile("../fixtures/ADT_A01_1.hl7", hl7v2.FixLineEndings())
	if err != nil {
		t.Fatal(err)
	}

	sch, err := msg.Schema()
	if err != nil {
		t.Fatal(err)
	}

	for _, seg := range sch.Segments {
		t.Logf("Segment: %s\n", seg.Name)
	}
}

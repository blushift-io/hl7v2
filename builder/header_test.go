package builder_test

import (
	"testing"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/builder"
)

func TestHeaderBuilderRefactored(t *testing.T) {
	mt := hl7v2.MessageType{Code: "ADT", Event: "A01", Structure: "ADT_A01"}
	hb := builder.BuildHeader(mt, hl7v2.Version251,
		builder.SendingApp("EPIC"),
		builder.SendingFacility("HOSPITAL_A"),
	)

	raw := hb.Build()
	if len(raw) < 5 {
		t.Fatalf("expected header segment fields, got %d", len(raw))
	}
	if string(raw[3][0][0][0]) != "EPIC" { // MSH-3
		t.Fatalf("expected EPIC sending app, got %s", raw[3][0][0][0])
	}
}

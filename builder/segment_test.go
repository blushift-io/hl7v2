package builder_test

import (
	"testing"

	"github.com/blushift-io/hl7v2/builder"
	"github.com/blushift-io/hl7v2/builder/field"
)

func TestSegmentBuilderRefactored(t *testing.T) {
	t.Run("Segment with Declarative Items", func(t *testing.T) {
		seg := builder.Segment("PID",
			field.Int(1),
			field.String(""),
			field.Components("12345", "", "", "MRN"),
			field.Components("Doe", "John", "A"),
		)

		raw := seg.Build()
		if len(raw) != 5 { // SEG_ID + 4 fields
			t.Fatalf("expected 5 raw fields in segment, got %d", len(raw))
		}
		if string(raw[0][0][0][0]) != "PID" {
			t.Fatalf("expected PID seg id, got %s", raw[0][0][0][0])
		}
	})

	t.Run("Segment with Fluent Set", func(t *testing.T) {
		seg := builder.Segment("PV1").
			Set(1, 1).
			Set(2, "I").
			Set(3, "ROOM1", "BED2")

		raw := seg.Build()
		if len(raw) != 4 { // SEG_ID + 3 fields
			t.Fatalf("expected 4 fields, got %d", len(raw))
		}
		if len(raw[3][0]) != 2 {
			t.Fatalf("expected 2 components in field 3, got %d", len(raw[3][0]))
		}
	})

	t.Run("Segment with Conditional Fields", func(t *testing.T) {
		seg := builder.Segment("PID", field.If(true, "ACTIVE"), field.If(false, "INACTIVE"))
		raw := seg.Build()
		if len(raw) != 3 {
			t.Fatalf("expected 3 fields, got %d", len(raw))
		}
		if string(raw[1][0][0][0]) != "ACTIVE" {
			t.Fatalf("expected ACTIVE, got %s", raw[1][0][0][0])
		}
		if len(raw[2]) != 0 {
			t.Fatalf("expected empty field for false condition")
		}
	})
}

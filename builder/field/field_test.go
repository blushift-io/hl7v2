package field_test

import (
	"testing"
	"time"

	"github.com/blushift-io/hl7v2/builder/field"
)

func TestFieldBuilders(t *testing.T) {
	t.Run("Primitive String Field", func(t *testing.T) {
		f := field.String("Doe")
		raw := f.Build()
		if len(raw) != 1 || string(raw[0][0][0]) != "Doe" {
			t.Fatalf("expected Doe, got %v", raw)
		}
	})

	t.Run("Primitive Int Field", func(t *testing.T) {
		f := field.Int(42)
		raw := f.Build()
		if string(raw[0][0][0]) != "42" {
			t.Fatalf("expected 42, got %v", raw)
		}
	})

	t.Run("Formatted Time Field", func(t *testing.T) {
		now := time.Date(2026, 8, 10, 15, 30, 0, 0, time.UTC)
		f := field.Time(now).Format("20060102150405")
		raw := f.Build()
		if string(raw[0][0][0]) != "20260810153000" {
			t.Fatalf("expected 20260810153000, got %v", raw)
		}
	})

	t.Run("Component Field", func(t *testing.T) {
		f := field.Components("Doe", "John", "A")
		raw := f.Build()
		if len(raw[0]) != 3 {
			t.Fatalf("expected 3 components, got %d", len(raw[0]))
		}
		if string(raw[0][0][0]) != "Doe" || string(raw[0][1][0]) != "John" {
			t.Fatalf("unexpected components: %v", raw)
		}
	})

	t.Run("Repeating Field", func(t *testing.T) {
		f := field.Repeated("555-1234", "555-5678")
		raw := f.Build()
		if len(raw) != 2 {
			t.Fatalf("expected 2 repetitions, got %d", len(raw))
		}
	})

	t.Run("Conditional Field When", func(t *testing.T) {
		f1 := field.When(true, "ACTIVE")
		if string(f1.Build()[0][0][0]) != "ACTIVE" {
			t.Fatalf("expected ACTIVE when true")
		}

		f2 := field.When(false, "ACTIVE")
		if len(f2.Build()) != 0 {
			t.Fatalf("expected empty field when false")
		}
	})

	t.Run("Repeated with Conditional Empty", func(t *testing.T) {
		f := field.Repeated(
			field.If(true, "First").Build(),
			field.If(false, "Second").Build(), // This produces empty RawField
			field.If(true, "Third").Build(),
		)
		raw := f.Build()
		if len(raw) != 2 {
			t.Fatalf("expected 2 repetitions, got %d", len(raw))
		}
		if string(raw[0][0][0]) != "First" || string(raw[1][0][0]) != "Third" {
			t.Fatalf("expected First and Third, got %v", raw)
		}
	})
}

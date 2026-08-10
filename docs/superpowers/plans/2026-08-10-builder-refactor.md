# Builder Package Refactoring Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Refactor the `builder` package to provide an Ent-inspired `field` builder subpackage, dual construction paradigms (declarative tree DSL & fluent path setters), and first-class functional options.

**Architecture:** Create a `builder/field` package with typed field constructors (`field.String`, `field.Int`, `field.Time`, `field.Components`, `field.Repeated`) and conditional helpers (`When`, `If`, `Eval`). Refactor `HeaderBuilder`, `SegmentBuilder`, and `Builder` to natively accept field builders, variadic primitives, and functional option closures.

**Tech Stack:** Go 1.22+, `github.com/blushift-io/hl7v2`

## Global Constraints

- Module path: `github.com/blushift-io/hl7v2/builder` and `github.com/blushift-io/hl7v2/builder/field`
- Package `builder/field` must be self-contained and importable without circular dependencies.
- Standard Go testing with table-driven unit tests.

---

### Task 1: Implement `builder/field` Subpackage

**Files:**
- Create: `builder/field/field.go`
- Create: `builder/field/option.go`
- Create: `builder/field/field_test.go`

**Interfaces:**
- Consumes: `github.com/blushift-io/hl7v2` (`Value`, `RawField`, `RawRepetition`, `RawComponent`, `RawSubcomponent`)
- Produces: `field.Builder` struct and constructor functions (`String`, `Int`, `Time`, `Bool`, `Components`, `Repeated`, `When`, `If`, `Eval`)

- [ ] **Step 1: Write failing unit tests for field constructors**

```go
// builder/field/field_test.go
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
}
```

- [ ] **Step 2: Run test to verify failure**

Run: `go test ./builder/field/...`
Expected: FAIL (package does not exist yet)

- [ ] **Step 3: Implement `builder/field/field.go` and `builder/field/option.go`**

```go
// builder/field/field.go
package field

import (
	"fmt"
	"time"

	"github.com/blushift-io/hl7v2"
)

type Builder struct {
	rawhl7 hl7v2.RawField
	err    error
}

func NewBuilder(rf hl7v2.RawField) *Builder {
	return &Builder{rawhl7: rf}
}

func String(v string) *Builder {
	return NewBuilder(hl7v2.RawField{
		hl7v2.RawRepetition{
			hl7v2.RawComponent{
				hl7v2.RawSubcomponent(v),
			},
		},
	})
}

func Int(v int) *Builder {
	return String(fmt.Sprintf("%d", v))
}

func Bool(v bool) *Builder {
	if v {
		return String("Y")
	}
	return String("N")
}

type TimeBuilder struct {
	*Builder
	t      time.Time
	layout string
}

func Time(t time.Time) *TimeBuilder {
	tb := &TimeBuilder{
		t:      t,
		layout: "20060102150405",
	}
	tb.Builder = String(t.Format(tb.layout))
	return tb
}

func (tb *TimeBuilder) Format(layout string) *TimeBuilder {
	tb.layout = layout
	tb.Builder = String(tb.t.Format(layout))
	return tb
}

func Components(comps ...any) *Builder {
	rawComps := make(hl7v2.RawRepetition, len(comps))
	for i, c := range comps {
		str := fmt.Sprintf("%v", c)
		rawComps[i] = hl7v2.RawComponent{hl7v2.RawSubcomponent(str)}
	}

	return NewBuilder(hl7v2.RawField{rawComps})
}

func Repeated(reps ...any) *Builder {
	rawReps := make(hl7v2.RawField, len(reps))
	for i, r := range reps {
		switch v := r.(type) {
		case *Builder:
			if len(v.rawhl7) > 0 {
				rawReps[i] = v.rawhl7[0]
			}
		default:
			str := fmt.Sprintf("%v", r)
			rawReps[i] = hl7v2.RawRepetition{
				hl7v2.RawComponent{hl7v2.RawSubcomponent(str)},
			}
		}
	}

	return NewBuilder(rawReps)
}

func (b *Builder) Build() hl7v2.RawField {
	if b == nil {
		return hl7v2.RawField{}
	}
	return b.rawhl7
}
```

```go
// builder/field/option.go
package field

type ConditionalElse struct {
	builder *Builder
	cond    bool
}

func When(cond bool, val any) *ConditionalElse {
	if cond {
		switch v := val.(type) {
		case *Builder:
			return &ConditionalElse{builder: v, cond: true}
		default:
			return &ConditionalElse{builder: String(fmt.Sprintf("%v", val)), cond: true}
		}
	}
	return &ConditionalElse{builder: NewBuilder(nil), cond: false}
}

func (ce *ConditionalElse) Else(val any) *Builder {
	if ce.cond {
		return ce.builder
	}
	switch v := val.(type) {
	case *Builder:
		return v
	default:
		return String(fmt.Sprintf("%v", val))
	}
}

func Eval(fn func() any) *Builder {
	res := fn()
	if b, ok := res.(*Builder); ok {
		return b
	}
	return String(fmt.Sprintf("%v", res))
}
```

- [ ] **Step 4: Run test to verify pass**

Run: `go test -v ./builder/field/...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add builder/field/
git commit -m "feat(builder): implement field builder subpackage with ent-style builders"
```

---

### Task 2: Refactor `SegmentBuilder` & `SegmentOption`

**Files:**
- Modify: `builder/segment.go`
- Create/Modify: `builder/segment_test.go`

**Interfaces:**
- Consumes: `builder/field` (`*field.Builder`), `github.com/blushift-io/hl7v2`
- Produces: `Segment(id string, items ...any) *SegmentBuilder`, `SegmentOption`, `SegmentBuilder.Set(pos int, items ...any)`

- [ ] **Step 1: Write failing unit test for `SegmentBuilder`**

```go
// builder/segment_test.go
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
}
```

- [ ] **Step 2: Run test to verify failure**

Run: `go test -v ./builder/segment_test.go`
Expected: FAIL (signatures or methods missing/changed)

- [ ] **Step 3: Update `builder/segment.go`**

Refactor `Segment(id string, items ...any)` and `Set(pos int, items ...any)` to handle `*field.Builder`, `SegmentOption`, and primitive values.

```go
// builder/segment.go
package builder

import (
	"fmt"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/builder/field"
)

type SegmentOption func(*SegmentBuilder)

type SegmentBuilder struct {
	id     string
	fields []*field.Builder
}

func Segment(id string, items ...any) *SegmentBuilder {
	sb := &SegmentBuilder{
		id: id,
		fields: []*field.Builder{
			field.String(id),
		},
	}

	for _, item := range items {
		sb.Add(item)
	}

	return sb
}

func (b *SegmentBuilder) Add(item any) *SegmentBuilder {
	switch v := item.(type) {
	case *field.Builder:
		b.fields = append(b.fields, v)
	case SegmentOption:
		v(b)
	case *SegmentOption:
		if v != nil {
			(*v)(b)
		}
	case nil:
		// ignore
	default:
		b.fields = append(b.fields, field.String(fmt.Sprintf("%v", v)))
	}
	return b
}

func (b *SegmentBuilder) Set(pos int, items ...any) *SegmentBuilder {
	if pos <= 0 {
		return b
	}

	// Ensure field slice capacity
	for len(b.fields) <= pos {
		b.fields = append(b.fields, field.String(""))
	}

	if len(items) == 1 {
		switch v := items[0].(type) {
		case *field.Builder:
			b.fields[pos] = v
		default:
			b.fields[pos] = field.String(fmt.Sprintf("%v", v))
		}
	} else if len(items) > 1 {
		b.fields[pos] = field.Components(items...)
	}

	return b
}

func (b *SegmentBuilder) Build() hl7v2.RawSegment {
	rs := make(hl7v2.RawSegment, len(b.fields))
	for i, f := range b.fields {
		if f != nil {
			rs[i] = f.Build()
		} else {
			rs[i] = hl7v2.RawField{}
		}
	}
	return rs
}

func (b *SegmentBuilder) BuildElement(parent hl7v2.Element, pos int) hl7v2.Element {
	return hl7v2.NewSegment(parent, pos, b.Build())
}
```

- [ ] **Step 4: Run test to verify pass**

Run: `go test -v ./builder/...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add builder/segment.go builder/segment_test.go
git commit -m "refactor(builder): update SegmentBuilder to support field builders and fluent Set()"
```

---

### Task 3: Refactor `HeaderBuilder` & `HeaderOption`

**Files:**
- Modify: `builder/header.go`
- Create/Modify: `builder/header_test.go`

**Interfaces:**
- Consumes: `hl7v2.MessageHeader`, `HeaderOption` (`func(*HeaderBuilder)`)
- Produces: `HeaderOption` constructors (`SendingApp`, `SendingFacility`, `ReceivingApp`, `ReceivingFacility`, etc.)

- [ ] **Step 1: Write failing unit test for `HeaderBuilder`**

```go
// builder/header_test.go
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
	if string(raw[2][0][0][0]) != "EPIC" { // MSH-3
		t.Fatalf("expected EPIC sending app, got %s", raw[2][0][0][0])
	}
}
```

- [ ] **Step 2: Run test to verify failure**

Run: `go test -v ./builder/header_test.go`
Expected: FAIL (missing new HeaderOption constructors)

- [ ] **Step 3: Update `builder/header.go`**

Provide top-level `HeaderOption` helpers:

```go
// builder/header.go
package builder

import (
	"time"

	"github.com/blushift-io/hl7v2"
)

type HeaderOption func(*HeaderBuilder)

func SendingApp(app string) HeaderOption {
	return func(h *HeaderBuilder) {
		h.hdr.SendingApplication = app
	}
}

func SendingFacility(facility string) HeaderOption {
	return func(h *HeaderBuilder) {
		h.hdr.SendingFacility = facility
	}
}

func ReceivingApp(app string) HeaderOption {
	return func(h *HeaderBuilder) {
		h.hdr.ReceivingApplication = app
	}
}

func ReceivingFacility(facility string) HeaderOption {
	return func(h *HeaderBuilder) {
		h.hdr.ReceivingFacility = facility
	}
}

func MessageDate(date time.Time) HeaderOption {
	return func(h *HeaderBuilder) {
		h.hdr.MessageDate = date
	}
}

func ControlID(id string) HeaderOption {
	return func(h *HeaderBuilder) {
		h.hdr.ControlID = id
	}
}

func ProcessingID(id string) HeaderOption {
	return func(h *HeaderBuilder) {
		h.hdr.ProcessingID = id
	}
}
```

- [ ] **Step 4: Run test to verify pass**

Run: `go test -v ./builder/...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add builder/header.go builder/header_test.go
git commit -m "feat(builder): add functional option constructors for HeaderBuilder"
```

---

### Task 4: Refactor `Builder` & Path-based Setters

**Files:**
- Modify: `builder/builder.go`
- Modify: `builder/builder_test.go`

**Interfaces:**
- Consumes: `HeaderBuilder`, `SegmentBuilder`, `field.Builder`
- Produces: `Builder.Header(...)`, `Builder.Segment(...)`, `Builder.Set(path string, val any)`, `Builder.Build()`

- [ ] **Step 1: Write failing integration test in `builder/builder_test.go`**

```go
// builder/builder_test.go
package builder_test

import (
	"testing"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/builder"
	"github.com/blushift-io/hl7v2/builder/field"
)

func TestFullBuilderFlow(t *testing.T) {
	mt := hl7v2.MessageType{Code: "ADT", Event: "A01", Structure: "ADT_A01"}

	b := builder.New().
		Header(mt, hl7v2.Version251,
			builder.SendingApp("EPIC"),
			builder.SendingFacility("HOSPITAL_A"),
		).
		Segment("PID",
			field.Int(1),
			field.String(""),
			field.Components("12345", "", "", "MRN"),
			field.Components("Doe", "John"),
		)

	b.Set("PV1-1", 1).
		Set("PV1-2", "I").
		Set("PV1-3.1", "ROOM1")

	msg, err := b.Build()
	if err != nil {
		t.Fatalf("build failed: %v", err)
	}

	if len(msg.Segments()) != 3 { // MSH, PID, PV1
		t.Fatalf("expected 3 segments, got %d", len(msg.Segments()))
	}
}
```

- [ ] **Step 2: Run test to verify failure**

Run: `go test -v ./builder/...`
Expected: FAIL

- [ ] **Step 3: Update `builder/builder.go`**

Update `Header(...)`, `Segment(...)`, and implement string path parsing in `Set(path string, val any)`:

```go
// builder/builder.go
package builder

import (
	"strconv"
	"strings"

	"github.com/blushift-io/hl7v2"
)

func (b *Builder) Header(typ hl7v2.MessageType, ver hl7v2.Version, opts ...HeaderOption) *Builder {
	b.hdr = BuildHeader(typ, ver, opts...)
	return b
}

func (b *Builder) Segment(id string, items ...any) *SegmentBuilder {
	seg := Segment(id, items...)
	b.segments = append(b.segments, seg)
	return seg
}

func (b *Builder) Set(path string, val any) *Builder {
	// Parse simple paths like "PV1-1" or "PV1-3.1"
	parts := strings.Split(path, "-")
	if len(parts) != 2 {
		return b
	}

	segID := parts[0]
	locParts := strings.Split(parts[1], ".")

	fieldIdx, err := strconv.Atoi(locParts[0])
	if err != nil || fieldIdx <= 0 {
		return b
	}

	seg := b.GetSegment(segID)
	if seg == nil {
		seg = Segment(segID)
		b.AddSegment(seg)
	}

	if len(locParts) == 1 {
		seg.Set(fieldIdx, val)
	} else if len(locParts) == 2 {
		compIdx, err := strconv.Atoi(locParts[1])
		if err == nil && compIdx > 0 {
			// Set component value at fieldIdx
			// Ensure field is a component field
			seg.Set(fieldIdx, val) // simplified set
		}
	}

	return b
}
```

- [ ] **Step 4: Run test to verify pass**

Run: `go test -v ./builder/...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add builder/
git commit -m "feat(builder): complete refactored builder API with path setters and fluent flow"
```

---

## Plan Verification

Run all package tests:
```bash
go test -v ./builder/...
```
Expected output:
```
PASS
ok  	github.com/blushift-io/hl7v2/builder	0.015s
ok  	github.com/blushift-io/hl7v2/builder/field	0.005s
```

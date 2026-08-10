# Message Builder Refactoring Design

**Date**: 2026-08-10  
**Status**: Proposed / Approved for Implementation  
**Target Package**: `github.com/blushift-io/hl7v2/builder` and `github.com/blushift-io/hl7v2/builder/field`  

---

## 1. Overview & Goals

The `builder` package in `hl7v2` provides a fluent API for constructing HL7 v2 messages. However, message construction currently requires verbose value wrapping (e.g., `SingleValueField(hl7v2.NewStringValue("1"))`), deeply nested component builders, and manual `query.Location` struct instantiations.

This refactoring redesigns the `builder` package to allow clean, expressive, readable, and type-safe message construction. Since breaking changes are acceptable prior to `v1.0.0`, the new design introduces:

1. **Ent-Inspired Field Builders (`field` package)**: Type-safe, fluid constructors for primitive fields (`field.String`, `field.Int`, `field.Time`, `field.Bool`), component fields (`field.Components`), and repeating fields (`field.Repeated`).
2. **Dual Message Construction Paradigms**:
   - **Declarative Tree DSL**: A clean hierarchical structure (`Segment("PID", field.Int(1), ...)`).
   - **Imperative Fluent Chaining & Path Setters**: Direct index setters (`seg.Set(1, "1")`) and path-based setters (`b.Set("PV1-3.1", "ROOM1")`).
3. **First-Class Functional Option Integration**: Seamless support for runtime-context conditional options (`HeaderOption`, `SegmentOption`) and inline field evaluation (`field.When`, `field.If`, `field.Eval`).

---

## 2. Package Architecture

```
builder/
├── builder.go        # Core Builder & path-based Setters
├── builder_test.go   # Integration tests for Builder
├── header.go         # HeaderBuilder & HeaderOption functional options
├── segment.go        # SegmentBuilder & SegmentOption functional options
└── field/            # Ent-inspired field builder subpackage
    ├── field.go      # String, Int, Time, Bool, Components, Repeated constructors
    ├── field_test.go # Unit tests for field builders
    └── option.go     # Conditional field helpers (When, If, Eval)
```

---

## 3. Detailed API Specifications

### 3.1 Field Builder Package (`builder/field`)

The `field` package provides constructors returning a `*field.Builder` that implements field construction, component/repetition assembly, and formatting modifiers.

#### Primitive Field Constructors
```go
field.String("Doe")                       // Simple string value
field.Int(1)                              // Integer value "1"
field.Time(time.Now()).Format("20060102") // Formatted time string
field.Bool(true)                          // "Y" / "N"
```

#### Composite & Repeating Field Constructors
```go
field.Components("Doe", "John", "A")       // Component field: Doe^John^A
field.Components("12345", "", "", "MRN")   // 12345^^^MRN
field.Repeated("555-1234", "555-5678")     // Repeating field: 555-1234~555-5678
field.Repeated(                            // Repeating components:
    field.Components("555-1234", "PRN"),   // 555-1234^PRN~555-5678^WPN
    field.Components("555-5678", "WPN"),
)
```

#### Modifiers & Formatting
- `.Format(layout string)`: Sets date/time format string.
- `.Optional()`: Omits field if value is empty/zero-value.
- `.Default(fallback string)`: Uses fallback if value is empty.
- `.Pad(length int, char byte)`: Pads field value.

#### Inline Conditional Logic
```go
field.When(cond, "VAL")                    // Includes "VAL" if cond is true
field.When(cond, "VAL_A").Else("VAL_B")    // Conditional ternary field
field.Eval(func() any { ... })             // Dynamic closure evaluation
```

---

### 3.2 Header Builder & Functional Options (`HeaderOption`)

`HeaderBuilder` accepts functional options (`type HeaderOption func(*HeaderBuilder)`), enabling domain-specific runtime configuration.

```go
type HeaderOption func(*HeaderBuilder)

// Built-in header options
func SendingApp(app string) HeaderOption
func SendingFacility(facility string) HeaderOption
func ReceivingApp(app string) HeaderOption
func ReceivingFacility(facility string) HeaderOption
func MessageDate(date time.Time) HeaderOption
func ProcessingID(id string) HeaderOption
func SequenceNumber(seq int) HeaderOption

// Custom domain options
func SetSendingFacilityFromReport(r Report, defaultFacility string) HeaderOption {
    return func(h *HeaderBuilder) {
        fac := defaultFacility
        if strings.TrimSpace(r.Location) == "Outside" {
            fac = strings.TrimSuffix(r.Facility, "_OUT")
        }
        h.SetSendingFacility(fac)
    }
}
```

---

### 3.3 Segment Builder & Functional Options (`SegmentOption`)

`SegmentBuilder` constructs HL7 v2 segments and accepts variadic field values, field builders, or `SegmentOption` functions (`type SegmentOption func(*SegmentBuilder)`).

```go
type SegmentOption func(*SegmentBuilder)

// Creating a segment with variadic fields or options
func Segment(id string, items ...any) *SegmentBuilder
```

#### Variadic Handling in `Segment(id, items...)` and `seg.Set(index, items...)`:
- If `item` is a `*field.Builder`, it is added as a field.
- If `item` is a `SegmentOption` (`func(*SegmentBuilder)`), it is executed on the builder.
- If `items` are strings/primitives passed to `Set(index, val1, val2, ...)`:
  - 1 value: added as a single value field.
  - $>1$ values: auto-converted into a `field.Components(val1, val2, ...)`.

---

### 3.4 Message Builder (`Builder`) & Path Setters

`Builder` aggregates the header and segment builders into a finalized `*hl7v2.Message`.

```go
// Path-based setters on Builder
b.Set("PV1-1", 1)
b.Set("PV1-2", "I")
b.Set("PV1-3.1", "ROOM1")
b.Set("PV1-3.2", "BED2")
```

---

## 4. Usage Examples

### 4.1 Declarative DSL Style
```go
msg, err := builder.New().
    Header(hl7v2.MessageType{Code: "ADT", Event: "A01"}, hl7v2.Version251,
        builder.SendingApp("EPIC"),
        SetSendingFacilityFromReport(report, "MAIN_FACILITY"),
    ).
    Segment("PID",
        field.Int(1),
        field.String(""),
        field.Repeated(
            field.Components("12345", "", "", "MRN"),
            field.Components("67890", "", "", "SSN"),
        ),
        field.Components("Doe", "John", "A"),
        field.When(patient.HasDOB(), field.Time(patient.DOB).Format("20060102")),
        field.String("M"),
    ).
    Segment("PV1",
        field.Int(1),
        field.String("I"),
        field.Components("ROOM1", "BED2", "NURSE3"),
    ).
    Build()
```

### 4.2 Fluent Chaining Style
```go
b := builder.New().
    Header(hl7v2.MessageType{Code: "ADT", Event: "A01"}, hl7v2.Version251).
    SendingApp("EPIC")

b.Segment("PID").
    Set(1, 1).
    Set(3, field.Components("12345", "", "MRN")).
    Set(5, "Doe", "John", "A").
    Set(7, field.Time(dob).Format("20060102")).
    Set(8, "M")

b.Set("PV1-1", 1).
  Set("PV1-2", "I").
  Set("PV1-3.1", "ROOM1")

msg, err := b.Build()
```

---

## 5. Verification Plan

1. **Unit Tests (`builder/field/field_test.go`)**:
   - Verify string, int, time, bool field construction.
   - Verify component fields (`field.Components`) encoding with delimiters.
   - Verify repeating fields (`field.Repeated`) encoding.
   - Verify conditional field helpers (`When`, `If`, `Eval`).
2. **Integration Tests (`builder/builder_test.go`)**:
   - Test full ADT_A01 message generation with `Header()`, `Segment()`, and `Build()`.
   - Verify MSH segment output formatting and delimiters.
   - Test path-based setting (`b.Set("PID-5.1", "Doe")`).
   - Run `go test ./builder/...` to confirm 100% pass rate.

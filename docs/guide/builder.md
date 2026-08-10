# Message Builder API (`builder` package)

The `builder` package provides a fluent, type-safe API for programmatically constructing HL7 v2 messages, headers, segments, fields, components, and subcomponents from scratch.

## Overview

Rather than manually manipulating raw strings or nested slices, the `builder` package allows constructing messages through chainable method calls and helper functions in `builder/field`.

```go
import (
	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/builder"
	"github.com/blushift-io/hl7v2/builder/field"
)
```

---

## Constructing a Message

Instantiate a builder with `builder.New()`, configure its header, append segments, and call `.Build()` to compile an `*hl7v2.Message`.

### 1. Header Options
Configure the `MSH` segment header with `b.Header(...)`:

```go
mt := hl7v2.MessageType{
	Code:      "ADT",
	Event:     "A01",
	Structure: "ADT_A01",
}

b := builder.New().
	Header(mt, hl7v2.Version251,
		builder.SendingApp("EPIC"),
		builder.SendingFacility("HOSPITAL_A"),
		builder.ReceivingApp("LAB_SYS"),
		builder.ReceivingFacility("MAIN_LAB"),
		builder.ProcessingID("P"),
	)
```

### 2. Adding Segments & Typed Fields

Use `b.Segment(id, items...)` alongside helpers from `builder/field`:

- `field.Int(1)`: Integer field value
- `field.String("text")`: String field value
- `field.Components("comp1", "comp2", "comp3")`: Multi-component field

```go
b.Segment("PID",
	field.Int(1),                             // PID-1 (Set ID)
	field.String(""),                         // PID-2
	field.Components("12345", "", "", "MRN"), // PID-3 (Patient Identifier List)
	field.Components("Doe", "John", "A"),     // PID-5 (Patient Name)
)
```

### 3. Path-Based Value Assignment (`b.Set`)

Assign field or component values using path location notation:

```go
b.Set("PV1-1", 1).
  Set("PV1-2", "I").
  Set("PV1-3.1", "ROOM101").
  Set("PV1-3.2", "BED2")
```

### 4. Compiling & Encoding

Compile the builder into a structured `*hl7v2.Message` and encode to ER7 bytes:

```go
msg, err := b.Build()
if err != nil {
	log.Fatalf("Build error: %v", err)
}

// Encode to ER7 byte slice
er7Bytes, err := msg.Encode()
if err != nil {
	log.Fatalf("Encoding error: %v", err)
}

fmt.Println(string(er7Bytes))
```

---

## Constructing Standalone Segments

You can also construct independent segment builders to modify or append to existing messages:

```go
package main

import (
	"fmt"
	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/builder"
	"github.com/blushift-io/hl7v2/builder/field"
)

func main() {
	// Existing message
	msg, _ := hl7v2.NewMessage([]byte("MSH|^~\\&|SEND|FAC|REC|FAC|20260810||ADT^A01|101|P|2.5\rPID|1||12345"))

	// Create custom Z-segment
	zseg := builder.Segment("ZXX").
		AddField(field.SingleValueField(hl7v2.NewStringValue("1"))).
		AddField(field.SingleValueField(hl7v2.NewStringValue("CUSTOM_VAL")))

	// Append element to message
	_ = msg.Append(zseg.BuildElement(msg, 0))

	encoded, _ := msg.Encode()
	fmt.Println(string(encoded))
}
```

---

## Complete Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/builder"
	"github.com/blushift-io/hl7v2/builder/field"
)

func main() {
	mt := hl7v2.MessageType{
		Code:      "ADT",
		Event:     "A01",
		Structure: "ADT_A01",
	}

	b := builder.New().
		Header(mt, hl7v2.Version251,
			builder.SendingApp("HIS_APP"),
			builder.SendingFacility("CLINIC_EAST"),
			builder.ReceivingApp("LIS_APP"),
			builder.ReceivingFacility("LAB_MAIN"),
		)

	b.Segment("PID",
		field.Int(1),
		field.String(""),
		field.Components("987654", "", "", "MR"),
		field.Components("SMITH", "ALICE", "M"),
	)

	b.Set("PV1-1", 1).
		Set("PV1-2", "O").
		Set("PV1-3.1", "OUTPATIENT_BUILDING")

	msg, err := b.Build()
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}

	output, err := msg.Encode()
	if err != nil {
		log.Fatalf("Encoding failed: %v", err)
	}

	fmt.Println("Generated Message:")
	fmt.Println(string(output))
}
```

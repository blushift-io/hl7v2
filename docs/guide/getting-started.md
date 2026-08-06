# Getting Started with hl7v2

`hl7v2` is a high-performance Go library for parsing, constructing, modifying, and transmitting Health Level Seven (HL7) v2 messages.

## Installation

Install `hl7v2` using standard Go module tooling:

```bash
go get github.com/blushift-io/hl7v2
```

## RawMessage vs Message: Choosing the Right API

`hl7v2` offers two distinct message representations tailored to different healthcare integration scenarios:

| Feature | `RawMessage` | `Message` |
| :--- | :--- | :--- |
| **Primary Goal** | Ultra-fast parsing, zero/low allocation | Rich object model & mutation |
| **Use Case** | Message routing, filtering, headers check | Building messages, field updates, unmarshaling |
| **Parsing Cost** | Minimal (slices over input bytes) | Full tree construction |
| **Mutation** | Read-only / append-only | Full Get/Set on segments, fields & components |
| **Schema Validation** | N/A | Supported against version schemas |

### Quick Example: RawMessage (Fast Read-Only Lookup)

Use `RawMessage` when you need to inspect fields or route messages quickly without building full segment trees:

```go
package main

import (
	"fmt"
	"log"

	"github.com/blushift-io/hl7v2"
)

func main() {
	rawHL7 := []byte("MSH|^~\\&|SENDing_APP|SENDING_FAC|REC_APP|REC_FAC|20260806120000||ADT^A01|MSG00001|P|2.5\rPID|1||12345^^^HOSP^MR||DOE^JOHN||19800101|M")

	raw, err := hl7v2.ParseRaw(rawHL7)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	patientID, _ := raw.Get("PID-3.1")
	fmt.Printf("Patient ID: %s\n", patientID)
}
```

### Quick Example: Message (Object Tree & Mutation)

Use `Message` when you need to update field values, add new segments, validate schemas, or generate ACK responses:

```go
package main

import (
	"fmt"
	"log"

	"github.com/blushift-io/hl7v2"
)

func main() {
	msg, err := hl7v2.NewMessage([]byte("MSH|^~\\&|SEND|FAC|REC|FAC|20260806||ADT^A01|101|P|2.5\rPID|1||999||SMITH^JANE"))
	if err != nil {
		log.Fatalf("Message error: %v", err)
	}

	pid, err := msg.Segment("PID")
	if err == nil {
		fmt.Printf("Segment Name: %s\n", pid.Name())
	}
}
```

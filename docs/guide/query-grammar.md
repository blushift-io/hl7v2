# Query Language & Grammar Selection

`hl7v2` includes a flexible query mechanism and a structural grammar parser for validating and selecting message segment sequences.

## Location Queries

Location queries use standard HL7 path notations:

- `PID-3` -> Field 3 of PID segment
- `PID-3.1` -> Component 1 of Field 3 in PID
- `MSH-9.1` -> Message Type (e.g. `ADT`)
- `MSH-9.2` -> Trigger Event (e.g. `A01`)

```go
package main

import (
	"fmt"
	"github.com/blushift-io/hl7v2"
)

func main() {
	raw, _ := hl7v2.ParseRaw([]byte("MSH|^~\\&|SEND|FAC|REC|FAC|20260806||ADT^A01|101|P|2.5\rPID|1||12345^^^HOSP^MR"))

	patientID, _ := raw.Get("PID-3.1")
	fmt.Printf("Patient ID: %s\n", patientID)
}
```

## Structural Grammar Selection

Grammar expressions allow validating and selecting message segments matching specific structural rules:
- `[...]` -> Optional segment
- `{...}` -> Repeating segment
- `(...)` -> Grouped segments

```go
package main

import (
	"fmt"
	"log"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/query"
)

func main() {
	msg, _ := hl7v2.NewMessage([]byte("MSH|^~\\&|SEND|FAC|REC|FAC|20260806||ADT^A01|101|P|2.5\rEVN|A01|20260806\rPID|1||12345"))

	// Parse a grammar expression
	grammars, err := query.ParseGrammars("MSH EVN PID [PD1] {NK1}")
	if err != nil {
		log.Fatalf("Grammar error: %v", err)
	}

	_ = grammars

	// Select matching elements from a Message
	elements, err := msg.Select("MSH EVN PID [PD1] {NK1}")
	if err == nil {
		fmt.Printf("Matched %d elements\n", len(elements))
	}
}
```

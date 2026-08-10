# Schema Inspection & Validation (`schema` package)

`hl7v2` provides embedded HL7 v2 specification schemas (v2.1 through v2.8) in the `schema` package. You can inspect message definitions, segment specifications, data types, and value tables, as well as validate messages against version-specific schemas.

## Loading a Schema

Register and open a schema by importing its spec package (e.g. `schema/spec/v21`) or version package (`versions/v21`):

```go
package main

import (
	"fmt"
	"log"

	"github.com/blushift-io/hl7v2/schema"
	_ "github.com/blushift-io/hl7v2/schema/spec/v21" // Registers HL7 v2.1 schema
)

func main() {
	sch := schema.Open("2.1")
	if sch == nil {
		log.Fatal("Schema v2.1 not found")
	}

	fmt.Printf("Loaded Schema Version: %s\n", sch.Version())
}
```

Alternatively, load it directly from the `v21` version package:

```go
import "github.com/blushift-io/hl7v2/versions/v21"

sch := v21.Schema()
```

## Validating Messages Against Schema

Validate a parsed `Message` against its schema using `ValidateMessageSchema`:

```go
package main

import (
	"fmt"
	"log"

	"github.com/blushift-io/hl7v2"
	_ "github.com/blushift-io/hl7v2/schema/spec/v21"
)

func main() {
	hl7Data := []byte("MSH|^~\\&|SEND|FAC|REC|FAC|20260806||ADT^A01|101|P|2.1\rEVN|A01|20260806\rPID|1||12345^^^HOSP^MR||DOE^JOHN")

	msg, err := hl7v2.NewMessage(hl7Data)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	result := hl7v2.ValidateMessageSchema(msg)
	if result.Valid {
		fmt.Println("Message is valid according to HL7 v2.1 schema!")
	} else {
		fmt.Println("Validation errors:")
		for _, err := range result.Errors {
			fmt.Printf(" - %v\n", err)
		}
	}
}
```

## Querying Message Types

Inspect message specifications, including required segments and repeat limits:

```go
msgSpec := sch.Message("ADT_A01")
fmt.Printf("Message: %s - %s\n", msgSpec.ID, msgSpec.Name)

for _, ms := range msgSpec.GetSegments() {
	fmt.Printf("Segment %s (Required: %t, MaxRepetitions: %d)\n",
		ms.ID, ms.Required(), ms.MaxRepetitions)
}

// List all messages in schema
for _, m := range sch.Messages() {
	fmt.Printf("Message ID: %s (%s)\n", m.ID, m.Name)
}
```

## Querying Segments & Fields

Look up segment field positions, component definitions, and max lengths:

```go
pidSpec := sch.Segment("PID")
fmt.Printf("Segment %s: %s\n", pidSpec.ID, pidSpec.Name)

for _, field := range pidSpec.Fields {
	fmt.Printf("Field %d: %s (%s) [MaxLen: %d, Required: %t]\n",
		field.Sequence, field.ID, field.Name, field.Length, field.Required())
}
```

## Querying Data Types

Inspect data type structures and component sequences:

```go
dtSpec := sch.DataType("CE") // Coded Element
fmt.Printf("DataType %s: %s\n", dtSpec.ID, dtSpec.Name)

for _, comp := range dtSpec.Components {
	fmt.Printf("Component %d: %s (%s)\n",
		comp.Sequence, comp.ID, comp.Name)
}

// List all data types in the schema
for _, dt := range sch.DataTypes() {
	fmt.Printf("DataType: %s - %s\n", dt.ID, dt.Name)
}
```

## Querying Value Tables

Explore HL7 standard tables and entry value codelists:

```go
tableSpec := sch.Table("0003") // Event Type Code
fmt.Printf("Table %s: %s (Type: %s)\n", tableSpec.ID, tableSpec.Name, tableSpec.Type)

for _, entry := range tableSpec.Entries {
	fmt.Printf("  Code: %s -> %s\n", entry.Value, entry.Description)
}

// List all tables in the schema
for _, tbl := range sch.Tables() {
	fmt.Printf("Table %s: %s\n", tbl.ID, tbl.Name)
}
```

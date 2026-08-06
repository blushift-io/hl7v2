# Raw Messages (`RawMessage`)

The `RawMessage` type provides a lightweight, slice-based parser designed for maximum speed and minimal memory allocations.

## Creating a RawMessage

Parse raw HL7 ER7 byte slices or streams using `ParseRaw` or `ReadRaw`:

```go
package main

import (
	"bytes"
	"fmt"
	"github.com/blushift-io/hl7v2"
)

func main() {
	data := []byte("MSH|^~\\&|SEND|FAC|REC|FAC|20260806||ADT^A01|101|P|2.5\rPID|1||12345")

	// From a byte slice
	raw, err := hl7v2.ParseRaw(data)
	if err != nil {
		panic(err)
	}

	// From an io.Reader
	reader := bytes.NewReader(data)
	rawFromReader, err := hl7v2.ReadRaw(reader)
	if err != nil {
		panic(err)
	}

	_ = rawFromReader
	fmt.Println("RawMessage parsed successfully")
}
```

## Extracting Fields & Values

You can retrieve values by location query strings:

```go
val, err := raw.Get("PID-3.1") // Patient ID Component 1
fmt.Printf("Patient ID: %s\n", val)

msgType, err := raw.Get("MSH-9") // Returns compound field "ADT^A01"
fmt.Printf("Message Type: %s\n", msgType)
```

## Iterating Raw Segments

```go
for _, seg := range raw.Segments() {
	fmt.Printf("Segment Name: %s\n", seg.Name())
}
```

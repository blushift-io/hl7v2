# Structured Messages (`Message`)

The `Message` type represents a full HL7 v2 object model hierarchy (`Message` -> `Segment` -> `Field` -> `Component` -> `Subcomponent`).

## Constructing and Reading Messages

```go
package main

import (
	"fmt"
	"github.com/blushift-io/hl7v2"
)

func main() {
	data := []byte("MSH|^~\\&|SEND|FAC|REC|FAC|20260806||ADT^A01|101|P|2.5\rPID|1||12345^^^HOSP^MR||DOE^JOHN")

	// From byte slice
	msg, err := hl7v2.NewMessage(data)
	if err != nil {
		panic(err)
	}

	// From file path
	// msgFromFile, err := hl7v2.NewMessageFromFile("message.hl7")

	fmt.Printf("Message delimiter: %s\n", string(msg.Delimiters().Field))
}
```

## Traversing Segments & Fields

```go
pid, err := msg.Segment("PID")
if err == nil {
	field, err := pid.Field(3) // Patient Identifier List
	if err == nil {
		fmt.Printf("Field Value: %s\n", field.Value())
	}
}
```

## Modifying Values & Building Segments

```go
// Set field value using path syntax
msg.Set("PID-5.1", "SMITH")
msg.Set("PID-5.2", "ALICE")

// Generate ACK Response
ack, err := msg.NewACK("AA")
if err == nil {
	ackBytes := ack.Value().Bytes()
	fmt.Println(string(ackBytes))
}
```

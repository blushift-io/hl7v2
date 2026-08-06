# Real-World Examples

## 1. High-Throughput HL7 Router

Using `RawMessage` to route messages based on `MSH-9` without full parsing overhead:

```go
package main

import (
	"fmt"

	"github.com/blushift-io/hl7v2"
)

func RouteHL7(data []byte) string {
	raw, err := hl7v2.ParseRaw(data)
	if err != nil {
		return "REJECT"
	}

	msgType, _ := raw.Get("MSH-9.1")
	switch msgType {
	case "ADT":
		return "PATIENT_QUEUE"
	case "ORU":
		return "LAB_QUEUE"
	default:
		return "DEFAULT_QUEUE"
	}
}

func main() {
	sample := []byte("MSH|^~\\&|SEND|FAC|REC|FAC|20260806||ADT^A01|101|P|2.5\rPID|1||12345")
	target := RouteHL7(sample)
	fmt.Printf("Routed to: %s\n", target)
}
```

## 2. ADT_A01 Acknowledgement Generator

```go
package main

import (
	"fmt"

	"github.com/blushift-io/hl7v2"
)

func ProcessAndAck(hl7Data []byte) ([]byte, error) {
	msg, err := hl7v2.NewMessage(hl7Data)
	if err != nil {
		return nil, err
	}

	ack, err := msg.NewACK("AA")
	if err != nil {
		return nil, err
	}

	return ack.Value().Bytes(), nil
}

func main() {
	sample := []byte("MSH|^~\\&|SEND|FAC|REC|FAC|20260806||ADT^A01|101|P|2.5\rPID|1||12345")
	ackBytes, err := ProcessAndAck(sample)
	if err == nil {
		fmt.Println("Generated ACK:")
		fmt.Println(string(ackBytes))
	}
}
```

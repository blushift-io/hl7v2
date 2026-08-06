# MLLP & Networking

Minimum Lower Layer Protocol (MLLP) is the standard framing protocol for transporting HL7 v2 messages over TCP sockets.

## MLLP Framing Details

MLLP wraps HL7 messages with framing bytes:
- **Start Block**: `0x0B` (`<VT>`, Vertical Tab)
- **End Block**: `0x1C 0x0D` (`<FS><CR>`, File Separator followed by Carriage Return)

## Using the MLLP Framer

The `mllp` package provides utilities to wrap read and write operations over net connections:

```go
package main

import (
	"fmt"
	"net"

	"github.com/blushift-io/hl7v2/mllp"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := mllp.NewReader(conn)
	for {
		msg, err := reader.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("Received MLLP message (%d bytes)\n", len(msg))
	}
}
```

# MLLP & Networking

Minimum Lower Layer Protocol (MLLP) is the standard framing protocol for transporting HL7 v2 messages over TCP sockets. `hl7v2` provides both low-level framing primitives (`mllp` package) and a high-level, concurrent TCP MLLP server framework (`net/tcp` package).

## MLLP Framing Protocol

MLLP wraps HL7 messages with framing bytes to mark message boundaries over raw TCP streams:
- **Start Block**: `0x0B` (`<VT>`, Vertical Tab)
- **End Block**: `0x1C 0x0D` (`<FS><CR>`, File Separator followed by Carriage Return)

## Low-Level MLLP Reader & Writer (`mllp` Package)

For custom socket handling, `mllp.NewReader` and `mllp.NewWriter` handle framing byte enforcement automatically:

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
	writer := mllp.NewWriter(conn)

	for {
		msgBytes, err := reader.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("Received raw MLLP message (%d bytes)\n", len(msgBytes))

		// Echo framing or send acknowledgment
		_ = writer
	}
}
```

## Setting Up an MLLP TCP Server (`net/tcp` Package)

The `net/tcp` package provides a full-featured MLLP TCP server framework. To handle incoming connections and messages, implement the `tcp.Handler` interface:

```go
type Handler interface {
	OnConnect(ctx *Context) error
	OnMessage(ctx *Context) error
	OnError(ctx *Context) error
	OnClose(ctx *Context) error
}
```

### Complete TCP Server Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/net/tcp"
)

type HL7Handler struct{}

func (h *HL7Handler) OnConnect(ctx *tcp.Context) error {
	log.Printf("Client connected: %s", ctx.Conn.RemoteAddr())
	return nil
}

func (h *HL7Handler) OnMessage(ctx *tcp.Context) error {
	msg, ok := ctx.Data.(*hl7v2.RawMessage)
	if !ok {
		return fmt.Errorf("unexpected message data type: %T", ctx.Data)
	}

	// Query message control ID
	cid, err := msg.QueryValue("MSH.10")
	if err != nil {
		return fmt.Errorf("failed to query MSH.10: %w", err)
	}

	log.Printf("Received message %s from %s", cid.String(), ctx.Conn.RemoteAddr())

	// Automatically construct and send ACK back to the client over MLLP
	_, err = ctx.Conn.AckMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send ACK: %w", err)
	}

	return nil
}

func (h *HL7Handler) OnError(ctx *tcp.Context) error {
	log.Printf("Connection error: %v", ctx.Error)
	return nil
}

func (h *HL7Handler) OnClose(ctx *tcp.Context) error {
	log.Printf("Client disconnected: %s", ctx.Conn.RemoteAddr())
	return nil
}

func main() {
	handler := &HL7Handler{}

	server, err := tcp.NewServer(handler, tcp.WithAddresses(":2525"))
	if err != nil {
		log.Fatalf("Failed to create MLLP server: %v", err)
	}

	log.Println("Starting MLLP server on :2525...")
	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
```

## Reference Server (`cmd/mllp_server`)

The repository includes a standalone production-ready MLLP server in [`cmd/mllp_server`](file:///Users/tnt/Projects/oss/hl7v2/cmd/mllp_server/main.go) that:
- Listens on TCP port `:2525`.
- Parses incoming messages into `*hl7v2.RawMessage`.
- Saves received HL7 messages and generated ACK responses to disk under `./tmp`.
- Uses `oklog/run` for graceful signal handling (`SIGINT`, `SIGTERM`).

### Running `cmd/mllp_server`

```bash
go run ./cmd/mllp_server
```

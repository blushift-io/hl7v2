# hl7v2

[![Go Version](https://img.shields.io/github/go-mod/go-version/blushift-io/hl7v2)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/blushift-io/hl7v2.svg)](https://pkg.go.dev/github.com/blushift-io/hl7v2)
[![Documentation](https://img.shields.io/badge/docs-blushift--io.github.io%2Fhl7v2-blue)](https://blushift-io.github.io/hl7v2)

> [!WARNING]
> **Pre-v1 API Notice**: This repository is under active development. The API is not yet stable and may undergo breaking changes until a `v1.0.0` release.

`hl7v2` is a high-performance, feature-complete Go library for Health Level Seven (HL7) v2 message parsing, object model manipulation, schema validation, transformation, and MLLP TCP network transmission.

Designed for modern healthcare integration engines and microservices, `hl7v2` balances ultra-fast zero-allocation routing with a rich, ergonomic document object model for message mutation and validation.

📖 **Full Documentation**: [https://blushift-io.github.io/hl7v2](https://blushift-io.github.io/hl7v2)

---

## Key Features

- **Dual Parsing Architecture**:
  - `RawMessage`: High-speed, zero/low-allocation parser ideal for message routing, filtering, and header inspection.
  - `Message`: Rich document tree for field mutation, segment insertion, unmarshaling, and schema validation.
- **Location Querying & Grammar Matching**:
  - Intuitive location syntax (`PID-3.1`, `MSH-9.2`).
  - Structural grammar parsing and selection for complex segment patterns (`MSH EVN PID [PD1] {NK1}`).
- **Built-in Schemas (v2.1 – v2.8+)**:
  - Embedded specification definitions for messages, segments, data types, and value tables across HL7 versions v2.1 through v2.8.
- **Automatic Acknowledgments**:
  - Effortless ACK/NACK message generation (`msg.Ack()`, `msg.Nack()`).
- **MLLP & TCP Networking**:
  - Low-level MLLP framing primitives (`mllp` package).
  - High-level concurrent TCP MLLP server framework (`net/tcp` package) with context-driven handlers.
  - Reference MLLP server implementation (`cmd/mllp_server`).
- **Transformation Engine**:
  - JavaScript-powered transformation scripts using Goja for dynamic HL7 message mapping.

---

## Installation

Install `hl7v2` using standard Go tooling:

```bash
go get github.com/blushift-io/hl7v2
```

Requires **Go 1.25** or higher.

---

## Quick Start

### 1. Fast Read-Only Query (`RawMessage`)

For message routing or header inspection without constructing a full object tree:

```go
package main

import (
	"fmt"
	"log"

	"github.com/blushift-io/hl7v2"
)

func main() {
	rawHL7 := []byte("MSH|^~\\&|SENDING_APP|SENDING_FAC|REC_APP|REC_FAC|20260810120000||ADT^A01|MSG00001|P|2.5\rPID|1||12345^^^HOSP^MR||DOE^JOHN||19800101|M")

	raw, err := hl7v2.ParseRaw(rawHL7)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	patientID, _ := raw.Get("PID-3.1")
	fmt.Printf("Patient ID: %s\n", patientID)
}
```

### 2. Message Tree & Mutation (`Message`)

For updating fields, appending segments, or validating against schemas:

```go
package main

import (
	"fmt"
	"log"

	"github.com/blushift-io/hl7v2"
)

func main() {
	hl7Data := []byte("MSH|^~\\&|SEND|FAC|REC|FAC|20260810||ADT^A01|101|P|2.5\rPID|1||999||SMITH^JANE")

	msg, err := hl7v2.NewMessage(hl7Data)
	if err != nil {
		log.Fatalf("Message error: %v", err)
	}

	pid, err := msg.Segment("PID")
	if err == nil {
		fmt.Printf("Segment Name: %s\n", pid.Name())
	}

	// Generate an ACK response
	ack, err := msg.Ack(hl7v2.AckCodeAA)
	if err == nil {
		fmt.Println(string(ack.Encode()))
	}
}
```

### 3. MLLP TCP Server (`net/tcp`)

Set up a concurrent TCP server that accepts HL7 messages over MLLP:

```go
package main

import (
	"log"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/net/tcp"
)

type ServerHandler struct{}

func (h *ServerHandler) OnConnect(ctx *tcp.Context) error { return nil }
func (h *ServerHandler) OnError(ctx *tcp.Context) error   { return nil }
func (h *ServerHandler) OnClose(ctx *tcp.Context) error   { return nil }

func (h *ServerHandler) OnMessage(ctx *tcp.Context) error {
	msg, ok := ctx.Data.(*hl7v2.RawMessage)
	if !ok {
		return nil
	}

	// Auto-send ACK over MLLP socket
	_, err := ctx.Conn.AckMessage(msg)
	return err
}

func main() {
	server, err := tcp.NewServer(&ServerHandler{}, tcp.WithAddresses(":2525"))
	if err != nil {
		log.Fatalf("Server create error: %v", err)
	}

	log.Println("Starting MLLP server on :2525...")
	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
```

---

## Documentation

Comprehensive guides, architecture overviews, and API references are available on our documentation site:

🌐 **[https://blushift-io.github.io/hl7v2](https://blushift-io.github.io/hl7v2)**

Topics covered:
- [Getting Started & API Comparison](https://blushift-io.github.io/hl7v2/guide/getting-started)
- [RawMessage Deep-Dive](https://blushift-io.github.io/hl7v2/guide/raw-message)
- [Message Object Model](https://blushift-io.github.io/hl7v2/guide/message)
- [Query Language & Grammar Selection](https://blushift-io.github.io/hl7v2/guide/query-grammar)
- [Schema Inspection & Validation](https://blushift-io.github.io/hl7v2/guide/schema)
- [Transformers & Scripting](https://blushift-io.github.io/hl7v2/guide/transformers)
- [MLLP Protocol & TCP Server](https://blushift-io.github.io/hl7v2/guide/mllp)

---

## Contributing

We welcome contributions from the healthcare and open-source community!

### Local Development Setup

1. **Prerequisites**:
   - [Go 1.25+](https://go.dev/dl/)
   - [pnpm](https://pnpm.io/) (for building/running documentation site)
   - [Mage](https://magefile.org/) (optional helper runner)

2. **Clone the repository**:
   ```bash
   git clone https://github.com/blushift-io/hl7v2.git
   cd hl7v2
   ```

3. **Run Unit Tests**:
   ```bash
   go test ./...
   ```

4. **Code Formatting & Verification**:
   ```bash
   go fmt ./...
   go vet ./...
   ```

5. **Documentation Development**:
   To preview the VitePress documentation site locally:
   ```bash
   # Using pnpm
   pnpm --prefix docs dev

   # Or using mage
   mage docsdev
   ```

For detailed pull request guidelines, commit standards, and architecture details, please see [CONTRIBUTING.md](./CONTRIBUTING.md).

---

## License

This project is licensed under the [MIT License](./LICENSE).
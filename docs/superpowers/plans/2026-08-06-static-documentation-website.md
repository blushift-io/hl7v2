# Static Documentation Website Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create a statically generated documentation website for `hl7v2` inside `docs/` using VitePress, `pnpm`, and GitHub Actions for GitHub Pages deployment.

**Architecture:** VitePress static site generator configured with site sidebar, search, dark mode, code snippets for Go and HL7 ER7, explicit guides contrasting `RawMessage` vs `Message`, and a GitHub Action workflow to build and publish docs on release.

**Tech Stack:** Node.js, `pnpm`, VitePress, Markdown, Vue 3, GitHub Actions, Mage.

## Global Constraints

- Source code lives under `docs/`.
- Package manager must be `pnpm`.
- Explicit distinction between `RawMessage` and `Message`.
- Query & Grammar selection syntax documented under `docs/guide/query-grammar.md`.
- Deploy target is GitHub Pages via `.github/workflows/docs.yml`.

---

### Task 1: Scaffold VitePress project in `docs/`

**Files:**
- Create: `docs/package.json`
- Create: `docs/.vitepress/config.mts`
- Create: `docs/index.md`

**Interfaces:**
- Produces: Base VitePress layout and site configuration runnable via `pnpm --prefix docs dev` and `pnpm --prefix docs build`.

- [ ] **Step 1: Create `docs/package.json`**

```json
{
  "name": "hl7v2-docs",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vitepress dev",
    "build": "vitepress build",
    "preview": "vitepress preview"
  },
  "devDependencies": {
    "vitepress": "^1.5.0",
    "vue": "^3.5.0"
  }
}
```

- [ ] **Step 2: Install dependencies using `pnpm`**

Run: `pnpm --prefix docs install`  
Expected: `pnpm-lock.yaml` created and dependencies installed in `docs/node_modules`.

- [ ] **Step 3: Create `docs/.vitepress/config.mts`**

```typescript
import { defineConfig } from 'vitepress'

export default defineConfig({
  title: "hl7v2",
  description: "High-Performance Go HL7 v2 Parser, Encoder, Query Engine, & MLLP Toolkit",
  base: "/hl7v2/",
  themeConfig: {
    logo: '/logo.svg',
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Guide', link: '/guide/getting-started' },
      { text: 'Examples', link: '/guide/examples' },
      {
        text: 'v0.1.0',
        items: [
          { text: 'v0.1.0 (Latest)', link: '/guide/getting-started' },
          { text: 'Release Notes', link: '/guide/versioning' }
        ]
      }
    ],

    sidebar: [
      {
        text: 'Getting Started',
        items: [
          { text: 'Introduction', link: '/guide/getting-started' }
        ]
      },
      {
        text: 'Core Guides',
        items: [
          { text: 'Raw Messages (Fast Parser)', link: '/guide/raw-message' },
          { text: 'Structured Messages (Object Model)', link: '/guide/message' },
          { text: 'Query & Grammar Selection', link: '/guide/query-grammar' },
          { text: 'MLLP & Networking', link: '/guide/mllp' }
        ]
      },
      {
        text: 'Resources',
        items: [
          { text: 'Real-World Examples', link: '/guide/examples' },
          { text: 'Versioning & Releases', link: '/guide/versioning' }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/blushift-io/hl7v2' }
    ],

    search: {
      provider: 'local'
    },

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2026 blushift-io'
    }
  }
})
```

- [ ] **Step 4: Create `docs/index.md`**

```markdown
---
layout: home

hero:
  name: "hl7v2"
  text: "Go HL7 v2 Parser & Encoder"
  tagline: High-performance, zero-allocation raw parser and complete HL7 v2 object model for Go.
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/blushift-io/hl7v2

features:
  - title: Dual Parsing Architecture
    details: Choose between high-speed RawMessage for read-only routing or Message for rich object mutation and unmarshaling.
  - title: Powerful Query Language
    details: Query fields and components using intuitive syntax (e.g. PID-3.1) and validate message structures with Grammars.
  - title: Production-Ready MLLP
    details: Built-in MLLP framer, TCP client, and server wrappers for reliable healthcare integration.
---
```

- [ ] **Step 5: Verify site build**

Run: `pnpm --prefix docs run build`  
Expected: Site successfully builds to `docs/.vitepress/dist`.

- [ ] **Step 6: Commit**

```bash
git add docs/package.json docs/pnpm-lock.yaml docs/.vitepress/config.mts docs/index.md
git commit -m "docs: scaffold VitePress documentation site with pnpm"
```

---

### Task 2: Create Primary Guide Pages (`getting-started.md`, `raw-message.md`, `message.md`)

**Files:**
- Create: `docs/guide/getting-started.md`
- Create: `docs/guide/raw-message.md`
- Create: `docs/guide/message.md`

**Interfaces:**
- Consumes: VitePress setup from Task 1.
- Produces: Comprehensive documentation explaining installation, `RawMessage` vs `Message`, and code examples.

- [ ] **Step 1: Create `docs/guide/getting-started.md`**

```markdown
# Getting Started with hl7v2

`hl7v2` is a high-performance Go library for parsing, constructing, modifying, and transmitting Health Level Seven (HL7) v2 messages.

## Installation

Install `hl7v2` using standard Go tools:

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
```

- [ ] **Step 2: Create `docs/guide/raw-message.md`**

```markdown
# Raw Messages (`RawMessage`)

The `RawMessage` type provides a lightweight, slice-based parser designed for maximum speed and minimal memory allocations.

## Creating a RawMessage

Parse raw HL7 ER7 byte slices or streams:

```go
// From a byte slice
raw, err := hl7v2.ParseRaw(data)

// From an io.Reader
raw, err := hl7v2.ReadRaw(reader)
```

## Extracting Fields & Values

You can retrieve values by location query strings:

```go
val, err := raw.Get("PID-5.1") // Family Name
fmt.Printf("Last Name: %s\n", val)

msgType, err := raw.Get("MSH-9") // Returns compound field "ADT^A01"
```

## Iterating Raw Segments

```go
for _, seg := range raw.Segments() {
    fmt.Printf("Segment: %s\n", seg.Name())
}
```
```

- [ ] **Step 3: Create `docs/guide/message.md`**

```markdown
# Structured Messages (`Message`)

The `Message` type represents a full HL7 v2 object model hierarchy (`Message` -> `Segment` -> `Field` -> `Component` -> `Subcomponent`).

## Constructing and Reading Messages

```go
// From bytes
msg, err := hl7v2.NewMessage(data)

// From file
msg, err := hl7v2.NewMessageFromFile("message.hl7")
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
// Set field value
msg.Set("PID-5.1", "SMITH")
msg.Set("PID-5.2", "ALICE")

// Generate ACK Response
ack, err := msg.NewACK("AA")
if err == nil {
    ackBytes := ack.Value().Bytes()
    fmt.Println(string(ackBytes))
}
```
```

- [ ] **Step 4: Verify site build**

Run: `pnpm --prefix docs run build`  
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add docs/guide/getting-started.md docs/guide/raw-message.md docs/guide/message.md
git commit -m "docs: add getting started, raw-message, and message guide pages"
```

---

### Task 3: Create Query/Grammar, MLLP, Examples, and Versioning Pages

**Files:**
- Create: `docs/guide/query-grammar.md`
- Create: `docs/guide/mllp.md`
- Create: `docs/guide/examples.md`
- Create: `docs/guide/versioning.md`

**Interfaces:**
- Consumes: VitePress setup from Task 1.
- Produces: Deep-dive documentation for location query expressions, Grammar selection, MLLP framing/networking, practical examples, and release versioning.

- [ ] **Step 1: Create `docs/guide/query-grammar.md`**

```markdown
# Query Language & Grammar Selection

`hl7v2` includes a flexible query mechanism and a structural grammar parser for validating message sequence rules.

## Location Queries

Location queries use standard HL7 path notations:

- `PID-3` -> Field 3 of PID segment
- `PID-3.1` -> Component 1 of Field 3 in PID
- `MSH-9.1` -> Message Type (e.g. `ADT`)
- `MSH-9.2` -> Trigger Event (e.g. `A01`)

```go
id, _ := msg.Get("PID-3.1")
```

## Structural Grammar Selection

Grammar expressions allow validating and selecting message segments matching specific structural rules (optional `[...]`, repeating `{...}`, or grouped `(...)` segments).

```go
import "github.com/blushift-io/hl7v2/query"

// Parse a grammar expression
grammars, err := query.ParseGrammars("MSH EVN PID [PD1] {NK1}")
if err != nil {
    log.Fatalf("Grammar error: %v", err)
}

// Select matching elements from a Message
elements, err := msg.Select("MSH EVN PID [PD1] {NK1}")
if err == nil {
    fmt.Printf("Matched %d elements\n", len(elements))
}
```
```

- [ ] **Step 2: Create `docs/guide/mllp.md`**

```markdown
# MLLP & Networking

Minimum Lower Layer Protocol (MLLP) is the standard framing protocol for transporting HL7 v2 messages over TCP sockets.

## MLLP Framing Details

MLLP wraps HL7 messages with framing bytes:
- **Start Block**: `0x0B` (`<VT>`)
- **End Block**: `0x1C 0x0D` (`<FS><CR>`)

## Using the MLLP Framer

```go
import "github.com/blushift-io/hl7v2/mllp"

// Wrapping an io.Reader or io.Writer with MLLP framing
framer := mllp.NewFramer(conn)
```
```

- [ ] **Step 3: Create `docs/guide/examples.md`**

```markdown
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
```
```

- [ ] **Step 4: Create `docs/guide/versioning.md`**

```markdown
# Versioning & Release Documentation

## Documentation Versioning

Documentation for `hl7v2` is versioned alongside project releases. 

- **`latest`**: Represents the current `main` development branch.
- **Tagged Releases**: Each version tag (e.g. `v0.1.0`) generates a snapshot available via the navigation header version selector.

## Building Docs for Releases

The documentation site is automatically built and published to GitHub Pages upon publishing a GitHub Release tag (`v*`) via GitHub Actions.
```

- [ ] **Step 5: Verify site build**

Run: `pnpm --prefix docs run build`  
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add docs/guide/query-grammar.md docs/guide/mllp.md docs/guide/examples.md docs/guide/versioning.md
git commit -m "docs: add query-grammar, mllp, examples, and versioning guide pages"
```

---

### Task 4: Add GitHub Actions CI Workflow and Magefile Targets

**Files:**
- Create: `.github/workflows/docs.yml`
- Modify: `magefiles/magefile.go`

**Interfaces:**
- Produces: GitHub Action workflow for building and deploying VitePress to GitHub Pages on release/push, and local `mage` targets for docs development.

- [ ] **Step 1: Create `.github/workflows/docs.yml`**

```yaml
name: Deploy Documentation

on:
  push:
    branches: [ main ]
  release:
    types: [ published ]
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: false

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Setup pnpm
        uses: pnpm/action-setup@v3
        with:
          version: 9

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: 20
          cache: 'pnpm'
          cache-dependency-path: docs/pnpm-lock.yaml

      - name: Install documentation dependencies
        run: pnpm --prefix docs install --frozen-lockfile

      - name: Build documentation site
        run: pnpm --prefix docs run build

      - name: Upload Pages artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: docs/.vitepress/dist

  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

- [ ] **Step 2: Add `DocsDev` and `DocsBuild` targets to `magefiles/magefile.go`**

```go
package main

import (
	"github.com/magefile/mage/sh"
)

var goDeps = map[string]string{
	"github.com/alvaroloes/enumer": "latest",
}

func Build() {}

func InstallDeps() error {
	for pkg, ver := range goDeps {
		if err := sh.Run("go", "install", pkg+"@"+ver); err != nil {
			return err
		}
	}

	return nil
}

// DocsDev starts the local VitePress documentation development server.
func DocsDev() error {
	return sh.RunV("pnpm", "--prefix", "docs", "run", "dev")
}

// DocsBuild builds the static documentation site.
func DocsBuild() error {
	return sh.RunV("pnpm", "--prefix", "docs", "run", "build")
}
```

- [ ] **Step 3: Test `pnpm --prefix docs build`**

Run: `pnpm --prefix docs build`  
Expected: Build passes clean and generates static output.

- [ ] **Step 4: Commit**

```bash
git add .github/workflows/docs.yml magefiles/magefile.go
git commit -m "ci: add GitHub Actions workflow for docs deployment and mage targets"
```

---

### Task 5: Final Verification & Build Check

**Files:**
- Test all created docs and build artifacts.

- [ ] **Step 1: Run full production build**

Run: `pnpm --prefix docs run build`  
Expected: `docs/.vitepress/dist` created with `index.html` and assets.

- [ ] **Step 2: Verify git status is clean**

Run: `git status`  
Expected: Clean working directory (excluding built assets/node_modules in `.gitignore`).

- [ ] **Step 3: Update `.gitignore` if needed**

Ensure `docs/.vitepress/dist`, `docs/.vitepress/cache`, and `docs/node_modules` are in `.gitignore`.

Check/Add to `.gitignore`:
```
docs/node_modules/
docs/.vitepress/dist/
docs/.vitepress/cache/
```

- [ ] **Step 4: Commit `.gitignore` changes if modified**

```bash
git add .gitignore
git commit -m "chore: add docs build artifacts to .gitignore"
```

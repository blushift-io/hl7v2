# Design Specification: Static Documentation Website for hl7v2

**Date**: 2026-08-06  
**Status**: Approved  
**Target Directory**: `docs/`  
**Framework**: VitePress (Vue-based static site generator)  
**Package Manager**: `pnpm`  
**Deployment Target**: GitHub Pages via GitHub Actions  

---

## 1. Executive Summary

This specification defines the architecture, content structure, tooling, and CI/CD workflow for the `hl7v2` project documentation website. The site will be statically generated from Markdown files inside the `docs/` directory, hosted on GitHub Pages, and support versioning for release documentation.

---

## 2. Directory Structure & Technology Stack

### 2.1 Tech Stack
- **Static Site Generator**: VitePress
- **Runtime / Package Manager**: Node.js & `pnpm`
- **Deployment Platform**: GitHub Pages (`gh-pages`)
- **CI Tooling**: GitHub Actions

### 2.2 Repository Layout
```
hl7v2/
├── docs/
│   ├── .vitepress/
│   │   ├── config.mts           # VitePress configuration (nav, sidebar, theme, search)
│   │   └── theme/
│   │       └── index.ts         # Theme customizations & CSS
│   ├── package.json             # Node dependencies (vitepress, vue)
│   ├── pnpm-lock.yaml           # pnpm lockfile
│   ├── index.md                 # Documentation Homepage
│   ├── guide/
│   │   ├── getting-started.md   # Overview, Installation, RawMessage vs Message
│   │   ├── raw-message.md       # High-performance read-only raw message parsing
│   │   ├── message.md           # Full message object model, get/set values, building segments
│   │   ├── query-grammar.md     # Location queries & Grammar selection syntax
│   │   ├── mllp.md              # MLLP client & server networking
│   │   ├── examples.md          # End-to-end workflow examples (ADT, ORU, ACK generation)
│   │   └── versioning.md        # Release versioning documentation
│   └── versions/                # Versioned documentation snapshots (e.g. v0.1/...)
├── magefiles/
│   └── magefile.go              # Added mage targets for docs:dev and docs:build
└── .github/
    └── workflows/
        └── docs.yml             # GitHub Actions deployment workflow
```

---

## 3. Documentation Content Architecture

### 3.1 Core Distinction: `RawMessage` vs `Message`
The documentation prominently distinguishes between the two parsing and representation models in `hl7v2`:

1. **`RawMessage` (`docs/guide/raw-message.md`)**
   - **Purpose**: High-speed, slice/byte-based parsing with zero or minimal allocations.
   - **Use Cases**: High-throughput message routing, filtering, headers inspection, read-only extraction.
   - **Key APIs**: `ParseRaw(b)`, `ReadRaw(r)`, `msg.Find(...)`, `msg.Get(...)`.

2. **`Message` (`docs/guide/message.md`)**
   - **Purpose**: Full object model representation with rich DOM-like navigation and mutation capabilities.
   - **Use Cases**: Constructing new HL7 messages, modifying existing fields/components, appending segments, struct unmarshaling, schema validation, generating ACKs.
   - **Key APIs**: `NewMessage(b)`, `ReadMessage(r)`, `msg.Segment(...)`, `msg.Set(...)`, `msg.Unmarshal(...)`, `NewACK(...)`.

### 3.2 Query & Grammar Selection (`docs/guide/query-grammar.md`)
- **Location Querying**: Explaining segment/field/component specifiers (e.g., `PID-3.1`, `MSH-9.2`, `OBX-5`).
- **Grammar Expression Selection**: Using structural grammar syntax (e.g., `m.Select("MSH EVN PID [PD1] {NK1}")`) via `query.ParseGrammars()` to validate expected segment order and select structured message elements.

### 3.3 MLLP & Networking (`docs/guide/mllp.md`)
- MLLP framing & unframing (`mllp.Framer`, `mllp.Reader`, `mllp.Writer`).
- TCP Client and Server implementations for transmitting raw or parsed HL7 streams over MLLP connections.

### 3.4 Real-World Examples (`docs/guide/examples.md`)
- **Example 1**: High-throughput HL7 message router using `RawMessage`.
- **Example 2**: ADT_A01 message builder & modifier using `Message`.
- **Example 3**: MLLP ACK response server using `net` / `mllp` packages.

---

## 4. Release & Versioning Strategy

- **Nav Selector**: The VitePress header will include a version dropdown menu (`latest`, release version tags).
- **Archive Snapshots**: When a new release tag (e.g. `v0.1.0`) is published, a snapshot of the documentation can be stored under `docs/versions/v0.1/` or updated in VitePress navigation configuration.

---

## 5. GitHub Actions CI/CD Workflow (`.github/workflows/docs.yml`)

The GitHub Action triggers on:
- Pushes to the `main` branch
- Published releases (`release` event)
- Manual workflow dispatch (`workflow_dispatch`)

### Key Workflow Steps
1. Checkout repository (`actions/checkout@v4`).
2. Setup `pnpm` (`pnpm/action-setup@v3`).
3. Setup Node.js 20 with `pnpm` caching (`actions/setup-node@v4`).
4. Install dependencies: `pnpm --prefix docs install --frozen-lockfile`.
5. Build static VitePress site: `pnpm --prefix docs run build` (outputs to `docs/.vitepress/dist`).
6. Deploy to GitHub Pages via `actions/upload-pages-artifact@v3` and `actions/deploy-pages@v4`.

---

## 6. Developer Experience & Local Tooling

To ensure seamless local execution for developers:
- `docs/package.json` contains:
  - `"dev"`: `"vitepress dev"`
  - `"build"`: `"vitepress build"`
  - `"preview"`: `"vitepress preview"`
- Magefile targets added to `magefiles/magefile.go`:
  - `DocsDev()`: Runs `pnpm --prefix docs dev`
  - `DocsBuild()`: Runs `pnpm --prefix docs build`

---

## 7. Spec Self-Review Checklist
- [x] **Placeholder Scan**: No TBD/TODO or vague statements.
- [x] **Internal Consistency**: Directory layout matches VitePress, pnpm scripts, and GitHub Action paths.
- [x] **Scope Check**: Focused strictly on documentation site scaffold, markdown content, release versioning, and CI workflow.
- [x] **Ambiguity Check**: Clear distinction between `RawMessage` and `Message`, as well as location query vs grammar selection.

# Contributing to hl7v2

Thank you for your interest in contributing to `hl7v2`! We welcome bug reports, feature requests, documentation improvements, and code contributions from the community.

---

## Code of Conduct & Guidelines

- **Respectful & Constructive**: Please keep discussions and code reviews constructive and welcoming.
- **Backwards Compatibility**: `hl7v2` prioritizes performance and API stability. Non-breaking changes and performance-conscious designs are preferred.
- **Test-Driven Development**: All new features and bug fixes must include unit tests.

---

## Local Development Setup

### Prerequisites

- **Go**: Version `1.25` or higher ([go.dev](https://go.dev/dl/)).
- **pnpm**: Package manager for documentation tooling ([pnpm.io](https://pnpm.io/)).
- **Mage**: Task runner for convenient task execution (`go install github.com/magefile/mage@latest`).

### Step-by-Step Setup

1. **Fork and Clone**:
   ```bash
   git clone https://github.com/<your-username>/hl7v2.git
   cd hl7v2
   ```

2. **Install Go & Doc Dependencies**:
   ```bash
   # Install documentation dependencies
   pnpm --prefix docs install
   ```

3. **Run Unit Tests**:
   Ensure all existing tests pass before making changes:
   ```bash
   go test ./...
   ```

4. **Run Benchmarks (Optional)**:
   ```bash
   go test -bench=. -benchmem ./...
   ```

---

## Code Style & Standards

- **Formatting**: Always format your Go code using `gofmt` or `go fmt ./...`.
- **Static Analysis**: Run `go vet ./...` to check for common mistakes and correctness issues.
- **Error Handling**: Follow idiomatic Go error handling. Errors should be descriptive and wrapped where appropriate using `%w`.
- **Comments**: Exported types, functions, and methods must have doc comments starting with the symbol name.
- **Zero Allocations for Hot Paths**: When contributing to `RawMessage` or low-level parsing routines, avoid unnecessary heap allocations.

---

## Documentation

Documentation is hosted using VitePress in the `docs/` directory and deployed to [https://blushift-io.github.io/hl7v2](https://blushift-io.github.io/hl7v2).

### Previewing Documentation Locally

```bash
# Using pnpm
pnpm --prefix docs dev

# Using mage
mage docsdev
```

Open [http://localhost:5173](http://localhost:5173) in your browser to view live updates.

### Building Documentation Static Site

```bash
# Using pnpm
pnpm --prefix docs run build

# Using mage
mage docsbuild
```

---

## Submitting Pull Requests

1. **Create a Feature Branch**:
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. **Commit Changes**:
   Write clear, concise commit messages that describe what changed and why.

3. **Verify locally**:
   Make sure all tests and linters pass:
   ```bash
   go test ./...
   go vet ./...
   ```

4. **Push and Open PR**:
   Push your branch to GitHub and open a Pull Request against the `main` branch of `blushift-io/hl7v2`. Provide a clear description of the problem solved or feature added in the PR description.

---

## Questions & Discussions

If you have questions about usage, design, or architecture, feel free to open a [GitHub Issue](https://github.com/blushift-io/hl7v2/issues) or start a discussion.

# JRX-001-001: Define the parser I/O contract and version detection

- **Status:** Not Started
- **Feature:** [JRX-001](../../FEATURE_MATRIX.md#jrx--jrxml-syntaxversion-compatibility-and-compilation)
- **Milestone:** M0
- **Priority:** P0
- **Dependencies:** FND-003-001
- **Blocked by:** none
- **Primary owner/package:** `internal/jrxml`
- **Likely touched files:** `internal/jrxml/doc.go`, `internal/jrxml/contract.go`, `internal/jrxml/contract_test.go`
- **Collision risk:** Medium — `internal/jrxml` will host many M1 tasks; this task defines shared types first, so others must follow its shape.
- **Parallel-safe with:** `DAT-002-001`, `PRT-001-001`, `EXP-001-001` (different packages)
- **Estimated size:** S

## Objective

Define the parser's input/output contract: the entry point, the template model
top-level types (minimal), dialect/version representation, and how diagnostics
are returned.

## Current state / prerequisites

`internal/jrxml/doc.go` is a stub. Diagnostics contract required.

## Scope (in)

- Entry point signature accepting a reader and returning a template model +
  diagnostics.
- A `Dialect`/`Version` representation (e.g. `6.x`, `unknown`) recorded from the
  root element.
- Top-level `Template` struct skeleton (no element details; those come in later
  tasks).
- Tests for version parsing of a couple of root declarations.

## Out of scope

- Parsing bands/elements/fields (M1 tasks).
- Normalization (`JRX-020`).

## Proposed contracts and integration points

```go
// Proposed — non-binding.
type Dialect struct { Major, Minor int; Raw string }
type Template struct { Dialect Dialect /* + metadata added later */ }
func Parse(r io.Reader) (*Template, diag.List, error)
```

Integration: `internal/compiler` consumes `*Template`.

## Implementation checklist (future agent only)

- [ ] Define `Dialect`, `Template`, `Parse`.
- [ ] Detect version from the root element if present; else `unknown`.
- [ ] Return diagnostics (not panic) on malformed XML.
- [ ] Tests for `6.20.0`, `6.0`, and absent version.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- Parsing a minimal `<jasperReport version="6.x">` yields the expected
  `Dialect`.
- Malformed XML yields a diagnostic with code `parse`, no panic.
- Unknown version yields a warning diagnostic, not an error.

## Tests and fixtures to add (future agent only)

- `internal/jrxml/contract_test.go` with inline XML strings (no `testdata` yet).

## Failure cases / edge cases

- Namespace prefixes; missing `version` attribute; BOM.

## Definition of done

Contract compiles and tests pass; documented in `internal/jrxml/doc.go`.

## References and downstream dependencies

- Feature: `JRX-001`; compatibility: [`../../COMPATIBILITY.md`](../../COMPATIBILITY.md).
- Unblocks: `JRX-001-002` and all other `internal/jrxml` M1 tasks.

# FND-001-001: Define the public SDK surface contract

- **Status:** Not Started
- **Feature:** [FND-001](../../FEATURE_MATRIX.md#fnd--foundation-public-api-cli-infrastructure)
- **Milestone:** M0
- **Priority:** P0
- **Dependencies:** FND-003-001, FND-009-001
- **Blocked by:** none
- **Primary owner/package:** `report` (proposed, module root)
- **Likely touched files:** `report/doc.go`, `report/api.go` (types/contracts only, no engine), `report/api_test.go`, `docs/ARCHITECTURE.md`
- **Collision risk:** Medium — establishes the public package name other tasks import. Coordinate the package name decision first.
- **Parallel-safe with:** `PRT-001-001`, `DAT-002-001`, `EXP-001-001` (different packages)
- **Estimated size:** M

## Objective

Decide and document the public API shape (package name, entry points, options,
format selection) so later tasks implement against a stable facade. This task
adds **types and documentation only**; it must not implement report behavior.

## Current state / prerequisites

No public API package. `cmd/gosper` is standalone. Requires the diagnostic and
limits contracts.

## Scope (in)

- Choose the public package path: `report` at module root (proposed) vs a
  subpackage. Record the decision (and rationale) in this task's PR and, if
  durable, an ADR.
- Define exported types/interfaces for `CompiledReport`, `FilledReport`,
  `Format`, `Option`, `Params`, and the function signatures for
  `Compile`/`Fill`/`Export`.
- Documentation stating these are **planned**; no implementation.

## Out of scope

- Any actual compile/fill/export behavior (later tasks).
- CLI wiring (`FND-005`).

## Proposed contracts and integration points

```go
// Proposed — non-binding.
func Compile(ctx context.Context, r io.Reader, opts ...Option) (*CompiledReport, error)
func Fill(ctx context.Context, c *CompiledReport, params Params, ds DataSource, opts ...Option) (*FilledReport, error)
func Export(ctx context.Context, f *FilledReport, w io.Writer, format Format) error
```

Integration: facade forwards to `internal/exec` (proposed) or orchestrates
directly — decide here (OQ-6).

## Implementation checklist (future agent only)

- [ ] Add `report` package with documented types and stubbed signatures
      returning `errors.New("not implemented")` **only if** the PR is explicitly
      a contract PR; otherwise leave behavior for `FND-001-002`.
- [ ] Ensure no `internal/*` package imports `report`.
- [ ] Update `docs/ARCHITECTURE.md` with the chosen package path.
- [ ] Add compile-only tests asserting the package builds and types exist.

## Acceptance criteria (observable)

- `go build ./...` succeeds.
- The chosen package path is documented and referenced from `ARCHITECTURE.md`.
- A test asserts `report` does not create import cycles (build is sufficient).

## Tests and fixtures to add (future agent only)

- `report/api_test.go`: minimal compile-time usage of the exported types.

## Failure cases / edge cases

- Naming clash with `internal/*`; avoid by keeping all public names in one
  package.

## Definition of done

Public contract documented, builds, no behavior claimed, decision recorded.

## References and downstream dependencies

- Feature: `FND-001`; OQ-6 in [`../../FEATURE_MATRIX.md`](../../FEATURE_MATRIX.md#open-questions).
- Unblocks: `FND-001-002`, `FIL-001-001`.

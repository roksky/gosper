# FND-006-001: Build and version metadata

- **Status:** Done
- **Feature:** [FND-006](../../FEATURE_MATRIX.md#fnd--foundation-public-api-cli-infrastructure)
- **Milestone:** M1
- **Priority:** P2
- **Dependencies:** none
- **Blocked by:** none
- **Primary owner/package:** `report` + `cmd/gosper`
- **Likely touched files:** `report/version.go`, `report/version_test.go`, `cmd/gosper/main.go` (use shared version)
- **Collision risk:** Low — small, isolated.
- **Parallel-safe with:** `LAY-005-001`, `internal/jrxml` tasks
- **Estimated size:** XS

## Objective

Expose a single build/version string from the library (overridable with
`-ldflags -X`), and have the CLI report the same value.

## Current state / prerequisites

CLI already has a `version` var (`dev`) and tests.

## Scope (in)

- `report.Version()` returning the injected value (default `dev`).
- CLI uses the shared value.
- Tests: default non-empty; override works.

## Out of scope

- Semantic versioning policy, changelog (`OPS-010`).

## Proposed contracts and integration points

```go
// Proposed — non-binding.
var version = "dev"
func Version() string
```

## Implementation checklist (future agent only)

- [x] Add shared version; wire CLI; tests.
- [x] Run `make check`.

## Acceptance criteria (observable)

- CLI `--version` matches the library value.
- Overriding via ldflags changes both.

## Tests and fixtures to add (future agent only)

- `report/version_test.go`.

## Failure cases / edge cases

- Empty override ⇒ fallback to `dev`.

## Definition of done

Implemented/tested; documented.

## References and downstream dependencies

- Feature: `FND-006`; `OPS-010`.

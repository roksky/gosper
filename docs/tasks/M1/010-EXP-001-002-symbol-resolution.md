# EXP-001-002: Implement expression symbol resolution

- **Status:** Done
- **Feature:** [EXP-001](../../FEATURE_MATRIX.md#exp--expression-language-and-type-semantics)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** EXP-001-001
- **Blocked by:** none
- **Primary owner/package:** `internal/expression`
- **Likely touched files:** `internal/expression/resolve.go`, `internal/expression/resolve_test.go`
- **Collision risk:** Medium — same package as `EXP-003-001`/`EXP-006-001`; separate files, agree on shared `Value` first.
- **Parallel-safe with:** `DAT-001-001`, `LAY-005-001`
- **Estimated size:** S

## Objective

Implement resolution of `$F{name}`, `$P{name}`, `$V{name}` against an injected
`Symbols` implementation, including scope precedence and missing-symbol
behavior.

## Current state / prerequisites

`EXP-001-001` defined `Symbols`/`Value`.

## Scope (in)

- Parse the three reference forms and dispatch to `Symbols`.
- Precedence: variable > parameter > field (documented), or per Jasper — decide
  and test.
- Missing symbol returns null/`(nil,false)` path used by `EXP-006-001`.

## Out of scope

- Full grammar/operators (`EXP-003-001`).
- Type coercion (`EXP-007`).

## Proposed contracts and integration points

Integrates with the evaluator built by `EXP-003-001`.

## Implementation checklist (future agent only)

- [x] Recognize `$F{}`/`$P{}`/`$V{}` and resolve.
- [x] Tests for each kind, precedence, and missing names.
- [x] Run `make check`.

## Acceptance criteria (observable)

- Each reference resolves to the injected value.
- Missing symbol follows the documented rule (no panic).

## Tests and fixtures to add (future agent only)

- `internal/expression/resolve_test.go`.

## Failure cases / edge cases

- Whitespace inside braces; empty name; nested braces.

## Definition of done

Resolution implemented/tested; documented.

## References and downstream dependencies

- Feature: `EXP-001`.
- Unblocks: `EXP-003-001`, `FIL-001-001`.

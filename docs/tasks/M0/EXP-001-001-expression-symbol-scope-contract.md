# EXP-001-001: Define the expression symbol and scope contract

- **Status:** Done
- **Feature:** [EXP-001](../../FEATURE_MATRIX.md#exp--expression-language-and-type-semantics)
- **Milestone:** M0
- **Priority:** P0
- **Dependencies:** none
- **Blocked by:** none
- **Primary owner/package:** `internal/expression`
- **Likely touched files:** `internal/expression/symbols.go`, `internal/expression/symbols_test.go`, `internal/expression/doc.go`
- **Collision risk:** Medium — `internal/expression` hosts multiple M1 tasks; define shared value/symbol types here first.
- **Parallel-safe with:** `DAT-002-001`, `PRT-001-001` (different packages)
- **Estimated size:** S

## Objective

Define how expressions receive symbols (`$F`, `$P`, `$V`) and the minimal
`Value` representation, so the parser and evaluator can be built without a
symbol store inside the package.

## Current state / prerequisites

`internal/expression/doc.go` is a stub.

## Scope (in)

- A `Value` abstraction covering null, bool, int64, decimal (placeholder type or
  `interface{}` with documented invariants), string, time.
- A `Symbols` interface for field/parameter/variable lookup.
- A `Scope` concept for lexical precedence (e.g. variables shadow parameters).
- Tests for lookup precedence and missing symbols.

## Out of scope

- Full type coercion (`EXP-007`), decimal backend (`EXP-008`), grammar.

## Proposed contracts and integration points

```go
// Proposed — non-binding.
type Symbols interface {
    Field(name string) (Value, bool)
    Parameter(name string) (Value, bool)
    Variable(name string) (Value, bool)
}
type Value interface{ isValue() }
```

Integration: `internal/expression` compile/eval take `Symbols`; `internal/fill`
implements it over run state.

## Implementation checklist (future agent only)

- [x] Define `Value` and `Symbols`; document null representation.
- [x] Leave decimal concrete type as a documented placeholder tied to OQ-1.
- [x] Tests for precedence and missing-symbol behavior.
- [x] Run `make check`.

## Acceptance criteria (observable)

- A test fake `Symbols` resolves according to the documented precedence.
- Missing symbol returns `(nil, false)` without panic.

## Tests and fixtures to add (future agent only)

- `internal/expression/symbols_test.go`.

## Failure cases / edge cases

- Case sensitivity of names; empty names.

## Definition of done

Contract compiles/tested; documented in package comment.

## References and downstream dependencies

- Feature: `EXP-001`; architecture:
  [`../../architecture/expressions-and-types.md`](../../architecture/expressions-and-types.md).
- Unblocks: `EXP-001-002`, `EXP-003-001`.

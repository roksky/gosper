# EXP-003-001: Expression parser and evaluator (arithmetic/comparison/boolean)

- **Status:** Not Started
- **Feature:** [EXP-003](../../FEATURE_MATRIX.md#exp--expression-language-and-type-semantics)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** EXP-001-002, EXP-019-001
- **Blocked by:** none
- **Primary owner/package:** `internal/expression`
- **Likely touched files:** `internal/expression/lexer.go`, `internal/expression/parser.go`, `internal/expression/eval.go`, tests
- **Collision risk:** High — core evaluator; keep null behavior in `EXP-006-001` and resolve `Value` types with `EXP-001-001`.
- **Parallel-safe with:** `DAT-001-001`, `LAY-005-001`, `JRX-002-001`
- **Estimated size:** L (split if needed: lexer/parser first, then evaluator)

## Objective

Implement a bounded lexer, parser and evaluator for literals, arithmetic,
comparison, boolean operators and parentheses, over the sandbox policy.

## Current state / prerequisites

Symbol resolution and sandbox policy in place.

## Scope (in)

- Tokens: number (int/decimal placeholder), string, bool, null, identifiers,
  operators `+ - * / % == != < <= > >= && || !` (and Jasper `and/or/not`).
- Precedence-correct parser producing an AST.
- Evaluator over `Value` with no side effects.
- Tests: precedence, comparisons, boolean truth tables, division-by-zero
  diagnostic.

## Out of scope

- Ternary (`EXP-005`), functions (`EXP-009`, `EXP-015`), Java methods
  (`EXP-014`), decimal backend finalization (`EXP-008`).

## Proposed contracts and integration points

```go
// Proposed — non-binding.
type Program struct { /* AST */ }
func Compile(src string, syms Symbols) (*Program, error)
func (p *Program) Eval(ctx context.Context, env Env) (Value, error)
```

## Implementation checklist (future agent only)

- [ ] Lexer with positions for diagnostics.
- [ ] Recursive-descent parser with documented precedence.
- [ ] Evaluator honoring `EXP-006` null rules once available.
- [ ] Fuzz the lexer/parser (no panics).
- [ ] Run `make check`.

## Acceptance criteria (observable)

- `1 + 2 * 3` evaluates to `7`.
- `a < b && b < c` yields expected boolean.
- Parse errors carry positions; evaluation never panics.

## Tests and fixtures to add (future agent only)

- `internal/expression/parser_test.go`, `eval_test.go`, and a fuzz test.

## Failure cases / edge cases

- Integer overflow; division by zero; deeply nested expressions (bound depth).

## Definition of done

Implemented/tested including fuzz; documented; compatibility rows for
arithmetic/comparison updated.

## References and downstream dependencies

- Features: `EXP-003`, `EXP-004`; architecture:
  [`../../architecture/expressions-and-types.md`](../../architecture/expressions-and-types.md).
- Unblocks: `EXP-005`, `EXP-006-001`, `FIL-008`.

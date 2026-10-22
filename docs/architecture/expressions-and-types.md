# Expressions and types

## Purpose and scope

Define the expression language scope, symbol resolution, type semantics, and
the decimal/null/date rules needed to be JasperReports-compatible for a
**documented subset**. This document constrains `internal/expression`.

## Scope

- **In scope:** `$F{field}`, `$P{parameter}`, `$V{variable}` references;
  literals; arithmetic, comparison and boolean operators; ternary conditionals;
  a curated function set; compile-time type checking.
- **Out of scope:** arbitrary Java/Groovy/JavaScript, reflection, class
  loading, method invocation on arbitrary objects, side effects.

## Symbol resolution

`internal/expression` does **not** store symbols. Callers pass a resolver:

```go
// Proposed — non-binding.
type Symbols interface {
    Field(name string) (Value, bool)
    Parameter(name string) (Value, bool)
    Variable(name string) (Value, bool)
}
```

Unresolved symbols are a **compile-time** error when statically known
(`EXP-016`), otherwise a runtime diagnostic with the expression span.

## Type semantics

- A small closed type set: null, boolean, integer (`int64`), decimal
  (arbitrary precision), string, date/time, and opaque/unknown.
- **Decimals are mandatory for financial expressions.** Binary floating point
  must not be the default (`EXP-008`); the concrete decimal representation is
  an open decision (see below).
- **Null propagation** follows Jasper-like rules; comparisons and arithmetic
  with null produce null unless a function defines otherwise (`EXP-006`).
- **Coercion** between numeric types is explicit and documented (`EXP-007`).

## Proposed contract — non-binding

```go
// Proposed — non-binding.
type Program struct { /* compiled AST; immutable */ }

func Compile(src string, syms Symbols) (*Program, error) // EXP-016
func (p *Program) Eval(ctx context.Context, env Env) (Value, error)
```

## Evaluation timing

Evaluation time (`Now`, `Report`, `Page`, `Column`, `Group`, `Band`) is part of
the compiled program metadata, not the grammar. It is honored by the fill and
layout stages, not by the evaluator, which only sees the current environment
(`FIL-005`, `FIL-010`).

## Failure modes

- Parse error → diagnostic with column/offset.
- Unknown symbol at compile time → error.
- Type mismatch → compile-time error where inferable, runtime otherwise.
- Division by zero / invalid date → runtime diagnostic, run aborts unless
  inside a guarded expression.

## Traps and tests

- Deterministic date/time formatting requires a pinned locale/timezone
  (`EXP-011`, `TST-001`).
- Decimal rounding must be banker's/half-up per Jasper behavior — decide and
  test explicitly (`EXP-010`, `EXP-013`).
- Fuzzing the lexer/parser must never panic (`TST-007`).

## Linked features and tasks

- `EXP-001` references — task `EXP-001-001` (contract), `EXP-001-002` (impl).
- `EXP-003` grammar — task `EXP-003-001`.
- `EXP-006` null semantics — task `EXP-006-001`.
- `EXP-019` sandbox — task `EXP-019-001`.

## Assumptions and unresolved alternatives

- **Unresolved:** decimal backend (`math/big.Rat`, a third-party decimal, or
  scaled integers). Tracked by task `EXP-008` (M2) and open question OQ-1 in
  [`STATUS.md`](../STATUS.md).
- **Unresolved:** depth of Java method translation (`EXP-014`).

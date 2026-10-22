# EXP — Expression language and type semantics

**Architecture:** [`../architecture/expressions-and-types.md`](../architecture/expressions-and-types.md).

## Purpose

Provide a bounded, sandboxed expression language covering `$F`/`$P`/`$V`
references, operators, conditionals, a curated function set, and deterministic
type semantics (null, decimal, date/time).

## Technical expectations

- **No arbitrary code execution.** A closed interpreter over a fixed AST
  (`EXP-019`).
- **Decimals are first-class** for financial math; binary float is not the
  default (`EXP-008`).
- **Null propagates** per Jasper-like rules (`EXP-006`).
- **Symbols are injected**, not stored (`EXP-001`).
- **Compile-time validation** where types are statically known (`EXP-016`).

### Proposed contract — non-binding

```go
type Value interface{ /* null | bool | int64 | decimal | string | time | unknown */ }
type Symbols interface {
    Field(string) (Value, bool)
    Parameter(string) (Value, bool)
    Variable(string) (Value, bool)
}
```

## Features

| ID | Expectation | Acceptance test (planned) |
| --- | --- | --- |
| EXP-001 | `$F`/`$P`/`$V` resolve with scope. | Scope precedence test. |
| EXP-002 | Literals (string/number/bool/null). | Literal parse test. |
| EXP-003 | Operator grammar + precedence. | Precedence table test. |
| EXP-004 | Boolean/comparison semantics. | Truth-table test. |
| EXP-005 | Ternary expressions. | `cond ? a : b` test. |
| EXP-006 | Null semantics. | Null propagation table. |
| EXP-007 | Type coercion rules. | Coercion matrix test. |
| EXP-008 | Decimal arithmetic. | Precision/rounding test. |
| EXP-009 | String functions. | Function behavior test. |
| EXP-010 | Numeric conversion/rounding. | Rounding-mode test. |
| EXP-011 | Date/time/locale/timezone. | Pinned-clock test. |
| EXP-012 | Date formatting. | Pattern test. |
| EXP-013 | Number/currency formatting. | Pattern test. |
| EXP-014 | Java-style method translation subset. | Supported-method test. |
| EXP-015 | Safe function registry. | Unknown function ⇒ error. |
| EXP-016 | Compile-time type validation. | Type-error test. |
| EXP-017 | Expression caching. | Cache hit test. |
| EXP-018 | Alt languages (Deferred). | n/a |
| EXP-019 | Sandbox: no arbitrary execution. | Malicious expression rejected. |

## Tasks

`EXP-001-001`, `EXP-019-001` (M0); `EXP-001-002`, `EXP-003-001`, `EXP-006-001`
(M1).

## Risks

- Decimal backend choice (`OQ-1`) affects everything downstream; decide early
  (`EXP-008`) and record an ADR.

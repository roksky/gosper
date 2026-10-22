# Concurrency and limits

## Purpose and scope

Define concurrency safety, cache isolation, cancellation and resource budgets
(`internal/limits`) so runs are predictable and safe.

## Concurrency model

- **Compiled reports are immutable** → any number of goroutines may fill them
  concurrently.
- **Fill states, cursors and print documents are single-goroutine** values.
- **Compilation cache** is concurrency-safe; entries are immutable once stored
  (`FND-004`).
- No package may hold mutable global state.

## Cancellation

- Every long-running entry point accepts `context.Context`.
- Cancellation is checked: before each stage, per record, per band/page, and in
  the exporter's write loop.
- `context.Canceled`/`DeadlineExceeded` are returned unchanged.

## Resource budgets

`internal/limits` provides a `Budget` and a single-goroutine `Guard` for:

- maximum records/rows,
- maximum pages,
- maximum expression evaluation steps,
- maximum bytes for resources and output.

A limit that is zero or negative means "unlimited". Exceeding a budget yields a
diagnostic with code `diag.CodeLimit` (`OPS-001`, `OPS-007`), never a panic or
unbounded allocation. Each guard call checks the context first and returns
`ctx.Err()` unchanged.

`report.WithBudget` applies these per-run limits. The M1 fill loop counts one
step per record, selected field, and parameter declaration; byte accounting
covers retained scalar field and parameter payloads plus their names. Layout
counts pages before appending them. The limits are checked at exact boundaries,
and accounting saturates rather than overflowing. Full expression-operation
and exporter-byte accounting remain for later budget work.

```go
type Budget struct {
    MaxRecords int
    MaxPages   int
    MaxSteps   int64
    MaxBytes   int64
}
```

Recursion depth (subreports) is bounded by a separate mechanism (`FIL-013`).

## Determinism under concurrency

Deterministic output must not depend on goroutine scheduling. Parallel fills
must produce identical bytes to serial fills for the same input.

## Traps and tests

- Race detector on all shared types (`TST-009`).
- Leak tests: after cancellation, cursors closed and goroutines drained
  (`TST-009`).
- Budget tests: each limit is exercised with a fixture that trips it
  (`OPS-001-001`).

## Linked features and tasks

- `FND-009` context/budget API — task `FND-009-001` (**Done**).
- `OPS-001` enforcement — task `OPS-001-001`.
- `OPS-007`, `DAT-008`, `FIL-013`, `TST-009`.

## Assumptions and unresolved alternatives

- **Unresolved:** whether budgets are per-run defaults or globally configured;
  default is per-run (`FND-009-001`).

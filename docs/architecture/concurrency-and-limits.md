# Concurrency and limits

## Purpose and scope

Define concurrency safety, cache isolation, cancellation and resource budgets
(`internal/limits`, proposed) so runs are predictable and safe.

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

`internal/limits` (proposed) provides counters for:

- maximum records/rows,
- maximum pages,
- maximum expression evaluation steps,
- maximum recursion depth (subreports),
- maximum bytes for resources and output.

Exceeding a budget yields a diagnostic (`OPS-001`, `OPS-007`), never a panic or
unbounded allocation.

```go
// Proposed — non-binding.
type Budget struct {
    MaxRecords int
    MaxPages   int
    MaxSteps   int64
    MaxBytes   int64
}
```

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

- `FND-009` context/budget API — task `FND-009-001`.
- `OPS-001` enforcement — task `OPS-001-001`.
- `OPS-007`, `DAT-008`, `FIL-013`, `TST-009`.

## Assumptions and unresolved alternatives

- **Unresolved:** whether budgets are per-run defaults or globally configured;
  default is per-run (`FND-009-001`).

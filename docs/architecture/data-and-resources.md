# Data and resources

## Purpose and scope

Define the data source model, query parameter binding, and the resource
resolution model (images, fonts, templates, URIs) in
`internal/datasource` and `internal/resources` (proposed), with explicit
security boundaries.

## Data source model

A data source opens a cursor; a cursor yields records; records expose typed
fields. `internal/fill` consumes records but does not know the source kind.

```go
// Proposed — non-binding.
type Record interface {
    Value(field string) (Value, bool)
}
```

Planned concrete sources: slices/structs/maps (`DAT-001`), JSON (`DAT-003`),
CSV (`DAT-004`), SQL via `database/sql` (`DAT-005`).

## Query and parameter binding

- SQL parameters are bound positionally/named through the driver's parameter
  API; **string concatenation of parameter values is forbidden** (`OPS-003`).
- Parameters may come from report parameters or runtime overrides.
- Field-to-Value conversion is centralized (`DAT-007`) to keep type semantics
  consistent with `internal/expression`.

## Resource resolution

```go
// Proposed interface — non-binding.
type Resolver interface {
    Open(ctx context.Context, ref string) (io.ReadCloser, error)
}
```

- Resolution supports injected filesystems and explicit allowlists
  (`OPS-002`).
- URIs (http/https/file) are **denied by default**; enabling them is explicit
  and audited.
- Images and fonts are referenced by logical name in compiled IR and resolved
  through the resolver at fill/layout time.

## Security boundaries

| Threat | Mitigation | Feature |
| --- | --- | --- |
| SSRF via template/URL | Deny-by-default URI allowlist | `OPS-002` |
| Path traversal | Rooted filesystem resolver | `OPS-002` |
| SQL injection | Parameterized queries only | `OPS-003` |
| Zip/decompression bombs | Size/ratio limits | `OPS-004` |
| Runaway recursion | Depth limits | `FIL-013` |

## Cancellation

`Open` and `Next` take `context.Context`; long-running queries must be
cancellable (`DAT-008`). Resources are closed even on error.

## Traps and tests

- Tests must assert that `file://` and `http://` refs are rejected by default
  (`OPS-002`).
- Deterministic SQL fixtures must not depend on ambient database state
  (`TEST_STRATEGY.md`).

## Linked features and tasks

- `DAT-001` — `DAT-001-001`; `DAT-002` — `DAT-002-001`; `DAT-003`/`DAT-004`
  (M2 tasks); `DAT-005` (M2 feature).
- `OPS-002`, `OPS-003`, `OPS-004`.

## Assumptions and unresolved alternatives

- **Assumption:** resource resolution is a separate injected interface.
- **Unresolved:** whether remote fonts/images are ever enabled, and under what
  policy (`OPS-002`, OQ-4).

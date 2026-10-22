# FND-002-001: Template loading and resource resolution

- **Status:** Not Started
- **Feature:** [FND-002](../../FEATURE_MATRIX.md#fnd--foundation-public-api-cli-infrastructure)
- **Milestone:** M1
- **Priority:** P1
- **Dependencies:** none
- **Blocked by:** none
- **Primary owner/package:** `internal/resources` (proposed)
- **Likely touched files:** `internal/resources/resolver.go`, `internal/resources/resolver_test.go`, `internal/resources/doc.go`
- **Collision risk:** Low — new package.
- **Parallel-safe with:** `LAY-005-001`, `FND-006-001`, `EXP-001-002`
- **Estimated size:** S

## Objective

Provide a deny-by-default resource resolver that loads templates (and later
images/fonts) from injected sources, blocking `file://`/`http://` unless
explicitly allowed.

## Current state / prerequisites

No resolver. `FND-002` feature is M1.

## Scope (in)

- A `Resolver` interface and a rooted-filesystem + in-memory implementation.
- Deny remote/absolute schemes by default; explicit allowlist.
- Context-aware `Open`.
- Tests: allowed local read, denied remote, denied traversal.

## Out of scope

- Image/font decoding (`LAY-005`, `JRX-012`).
- HTTP sources (`DAT-011`, `OPS-002`).

## Proposed contracts and integration points

```go
// Proposed — non-binding.
type Resolver interface{ Open(ctx context.Context, ref string) (io.ReadCloser, error) }
```

Integration: `report.Compile` uses it to read the template.

## Implementation checklist (future agent only)

- [ ] Implement resolver + options for allowlists.
- [ ] Reject `..` traversal outside root.
- [ ] Tests for allow/deny cases.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- `file://`/`http://` refs error with code `unsupported` (or `limit`) by
  default.
- A relative ref inside the root opens; `../` outside root is rejected.

## Tests and fixtures to add (future agent only)

- `internal/resources/resolver_test.go` using `fstest.MapFS`.

## Failure cases / edge cases

- Symlinks escaping root; very large files (defer size limits to `OPS-004`).

## Definition of done

Resolver tested; deny-by-default documented; linked from
[`../../architecture/data-and-resources.md`](../../architecture/data-and-resources.md).

## References and downstream dependencies

- Feature: `FND-002`; security: `OPS-002`.
- Unblocks: `FND-001-002`, later image/font tasks.

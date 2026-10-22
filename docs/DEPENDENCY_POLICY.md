# Dependency policy

> Establishes how third-party dependencies are justified, license-reviewed and
> tracked. Owned by `OPS-008` (task `OPS-008-001`). The project license itself
> is still undecided — see [`LICENSING.md`](LICENSING.md).

## Principles

1. **Standard library first.** Prefer the standard library. A dependency must
   earn its place.
2. **Justify every addition.** A new dependency requires a written rationale in
   the PR: what it provides, why the standard library is insufficient, and its
   maintenance status.
3. **License review before merge.** The license of every new dependency must be
   reviewed for compatibility with the project license (pending) and with the
   licenses of existing dependencies. This is a **hard gate**.
4. **No speculative dependencies.** Do not add a dependency solely for a future
   feature. Notable expected candidates are a PDF library (`OUT-001`, OQ-5) and
   a text-shaping backend (`LAY-007`, OQ-3); each needs a review when it is
   actually introduced.
5. **Small, pinned, reproducible.** Dependencies are pinned via `go.mod` /
   `go.sum`. Prefer small, well-maintained modules with a stable release history.
6. **No runtime JVM.** No dependency may introduce a JVM requirement
   (`ADR-0001`).

## Review checklist for a new dependency

- [ ] Rationale recorded in the PR description.
- [ ] License identified and compatible with the project license.
- [ ] Transitive dependency licenses reviewed (`go mod graph`).
- [ ] Maintenance signals checked (recent releases, issue activity).
- [ ] Size/scope proportionate to the value provided.
- [ ] No network access at build or runtime implied by the dependency.
- [ ] `go.mod`/`go.sum` updated via `go mod tidy`; no unrelated modules pulled in.

## SBOM and dependency inventory

**Decision: deferred, with a documented manual process.**

Rationale: the module currently has **zero** third-party dependencies, so SBOM
tooling would add cost and CI surface without value. When the first non-standard
dependency is added, the SBOM approach will be chosen and pinned as part of that
same PR.

Until then, the dependency inventory is produced manually with the Go toolchain
(no extra tooling, no network):

```sh
go list -m all          # module graph
go version -m ./bin/gosper  # modules embedded in a built binary (requires make build)
```

When SBOM tooling is adopted it MUST be:

- pinned to a specific version,
- run read-only in CI with `permissions: contents: read` and no secrets,
- free of network access beyond module fetching, or marked optional.

No SBOM tooling is added to CI at this time.

## Relationship to other work

- **`OPS-009`** covers the separate JasperReports interoperability/licensing
  review for any reference test tooling; that concern is distinct from
  production Go dependencies.
- **[`../CONTRIBUTING.md`](../CONTRIBUTING.md)** restates the "justify every new
  dependency" rule for code contributors.

## References

- [`LICENSING.md`](LICENSING.md)
- [`adr/0001-native-go-library-first-runtime.md`](adr/0001-native-go-library-first-runtime.md)
- `OPS-008`, `OPS-009` in [`FEATURE_MATRIX.md`](FEATURE_MATRIX.md)

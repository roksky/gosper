# Test strategy

> Records testing **strategy** and clearly separates what runs **today** from
> future plans. This supersedes the earlier bootstrap document previously at
> `docs/testing.md` (merged here).

## What runs today

- `go test ./...` — unit tests for CLI metadata/help/version behavior.
- `go build ./cmd/gosper` — a compile/smoke check exercised in CI.
- `gofmt` and `go vet` — enforced in CI and via `make check`.
- `go test -race ./...` — race detector on Linux in CI.

Run locally:

```sh
make check
make test-race
make coverage
```

No parser/layout/exporter test exists yet because no engine code exists.

## Fixture conventions (planned)

- **Location:** `testdata/` at the repository root; per-suite subdirectories.
- **Inputs:** `.jrxml` templates and data fixtures (JSON/CSV) side by side.
- **Expected:** golden outputs under `testdata/<suite>/golden/`.
- **Naming:** `<feature-id>-<NNN>-<slug>` matching the compatibility registry
  (`JRX-022`).
- **Updates:** goldens are updated only via an explicit `-update` flag and must
  be reviewed in the diff; never auto-overwritten silently.

See task `TST-001-002` and `TST-002-001`.

## Deterministic inputs

All output-affecting environment must be pinned in tests:

| Factor | Control |
| --- | --- |
| Time/timezone | injected clock + fixed `time.Location` (`TST-001-001`) |
| Locale | fixed locale for number/date formatting |
| Fonts | bundled test fonts with stable metrics (`LAY-005`) |
| Database | ephemeral test DB or in-memory fixtures, no ambient state |
| Map iteration | never rely on Go map order in output |

Determinism is a requirement (`OUT-010`, `LAY-015`).

## Test layers

| Layer | Purpose | Feature/task |
| --- | --- | --- |
| Unit | per-package contracts, edge cases | each module |
| Contract | interfaces and boundaries honored | `FND-001`, `DAT-002` |
| Golden | stable serialized output | `TST-002` |
| Semantic | structure/values, not bytes | `TST-003` |
| Pagination | page counts and break positions | `TST-004` |
| Visual | rendered-image comparison with tolerance | `TST-005` |
| Differential | compare to JasperReports (test-only, isolated) | `TST-006` |
| Fuzz/property | robustness, no panics | `TST-007`, `TST-008` |
| Race/leak/cancel | concurrency and resource safety | `TST-009` |
| Benchmark | declared hardware + datasets | `TST-010` |
| Cross-platform smoke | Linux/macOS/Windows parity | `TST-011` |

## Assertions and tolerances

- Prefer **semantic assertions** over byte comparisons where layout may shift.
- Golden byte comparisons are acceptable where determinism is guaranteed
  (e.g. CSV, and PDF once `OUT-010` lands).
- Visual comparisons define an explicit tolerance and diff artifact policy.
- Pagination assertions check page count and the set of elements per page.

## Differential testing (future)

- Runs in an **isolated, test-only** environment with JasperReports tooling
  obtained under appropriate license terms (`OPS-009`).
- Never part of the production path or default CI unless licensing permits.
- Outputs compared with documented tolerances; differences are recorded as
  known deviations.

## Security testing

- Fuzz parse/evaluate paths for panics and resource exhaustion (`TST-007`).
- Test allowlists and traversal defenses (`OPS-002`, `OPS-003`).
- Test budget enforcement (`OPS-001`).

## Performance measurement methodology

- Benchmarks declare hardware, Go version, dataset size and font set.
- Report allocations and bytes, not just time.
- No benchmark numbers are claimed until measured (none measured today).

## Fixture ideas (not yet created)

Template ideas for future fixtures: minimal text report; tabular multi-page
report; grouped report with variables; subreport composition; shapes and
borders; Unicode text; conditional styles. Data ideas: slice, JSON, CSV, SQL.
These are **ideas**, not existing files; `testdata/` remains empty of engine
fixtures until the features that use them land.

## References

- [`COMPATIBILITY.md`](COMPATIBILITY.md)
- [`architecture/concurrency-and-limits.md`](architecture/concurrency-and-limits.md)
- [`tasks/INDEX.md`](tasks/INDEX.md)

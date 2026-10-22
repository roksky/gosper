# Status

> **Derived snapshot.** This is a summary of the evidenced baseline, not the
> source of truth. Canonical status lives in
> [`FEATURE_MATRIX.md`](FEATURE_MATRIX.md) (features) and
> [`tasks/INDEX.md`](tasks/INDEX.md) (tasks). Snapshot date: **2026-10-22**.

## Repository facts (observed)

| Fact | Value | Evidence |
| --- | --- | --- |
| Module path | `github.com/roksky/gosper` | `go.mod` |
| Go version | `1.26` | `go.mod` (`go 1.26`) |
| Branch / commit | `dev` @ `7a787f4` "Setup Repo" | `git log`, `git status` |
| Working tree | clean (no uncommitted changes) | `git status --short` |
| Remote | `git@github.com:roksky/gosper.git` | `git remote -v` |
| License | **pending** | [`LICENSING.md`](LICENSING.md) (no `LICENSE` file) |

## Implemented capabilities (evidenced)

| Capability | State | Evidence |
| --- | --- | --- |
| Go module + builds | Present | `go.mod`; `go build ./cmd/gosper` |
| CLI metadata: `--help`/`help`, `--version`/`version`, unknown-command error | Present | `cmd/gosper/main.go`, `cmd/gosper/main_test.go` |
| Unit tests for CLI metadata | Present | `cmd/gosper/main_test.go` |
| Formatting/vet/test/race automation | Present | `Makefile` (`check`, `test-race`) |
| CI (gofmt, vet/test/build, race) | Present | `.github/workflows/ci.yml` (Linux; macOS/Windows temporarily disabled) |
| Diagnostics model (`internal/diag`) | Present | `internal/diag/diag_test.go` |
| Resource budgets (`internal/limits`) | Present | `internal/limits/limits_test.go` |
| Expression symbols/value (`internal/expression`) | Present | `internal/expression/symbols_test.go` |
| Print model + unit (`internal/printdoc`) | Present | `internal/printdoc/model_test.go` |
| Data-source cursor contract (`internal/datasource`) | Present | `internal/datasource/datasource_test.go` |
| Deterministic test helpers (`internal/testutil`) | Present | `internal/testutil/env_test.go` |
| Public SDK contract (`report`) | Contract stubs | `report/api_test.go` (entry points return `ErrNotImplemented`) |
| JRXML parser contract (`internal/jrxml`) | Version detection | `internal/jrxml/contract_test.go` |
| Editor/build hygiene | Present | `.editorconfig`, `.gitignore` |
| Documentation (bootstrap) | Present | `README.md`, `CONTRIBUTING.md`, `docs/` |

> Note: the CLI is genuinely implemented and tested (`FND-011` = `Done`). No
> report functionality exists.

## Scaffold-only (not implementation evidence)

The following are **compile-only `doc.go` stubs** with no types or logic. Their
presence does **not** indicate implemented features:

`internal/jrxml`, `internal/compiler`, `internal/fill`, `internal/layout`,
`internal/exporter`.

Packages with only contract-level code (interfaces/model, no engine behavior):
`internal/expression`, `internal/datasource`, `internal/printdoc`.

## Not present

- JRXML parsing, compilation, expression evaluation.
- Data sources, fill/iteration, layout/pagination, print IR.
- Any exporter (PDF/XLSX/HTML/CSV).
- Public report API (`Compile`/`Fill`/`Export`).
- Network service, plugin system, report designer.
- A `LICENSE` file.

## Aggregate progress (derived)

From [`FEATURE_MATRIX.md`](FEATURE_MATRIX.md):

| Status | Count |
| --- | --- |
| Done | 9 (`FND-002`, `FND-003`, `FND-009`, `FND-011`, `DAT-002`, `EXP-019`, `JRX-022`, `OPS-008`, `TST-001`) |
| In Progress | 4 (`EXP-001`, `PRT-001`, `FND-001`, `JRX-001`) |
| Deferred | 3 (`JRX-021`, `EXP-018`, `PRT-006`) |
| Not Started | 130 |
| **Total** | **146** |

Conditional verification note: `go test ./...`, `go vet ./...`, `go build`, and
`go test -race ./...` were observed passing in the prior bootstrap task
(changeset that produced commit `7a787f4`). They were **not** re-run as part of
this documentation-only pass.

## Documentation reconciliation (this pass)

The required planning hierarchy uses uppercase filenames that would **case-collide**
with existing lowercase bootstrap docs on macOS/Windows. To avoid divergent
sources of truth, the earlier `docs/architecture.md`, `docs/compatibility.md`
and `docs/testing.md` were **merged** into `docs/ARCHITECTURE.md`,
`docs/COMPATIBILITY.md` and `docs/TEST_STRATEGY.md`, and the lowercase files
were removed. Existing decisions were preserved. `docs/development.md` and
`docs/LICENSING.md` have no required counterpart and were kept. Inbound links
in `README.md` and `CONTRIBUTING.md` were updated.

## Curated blockers / open questions

| # | Question | Affects | Process |
| --- | --- | --- | --- |
| OQ-1 | Decimal backend | `EXP-008` | ADR when `EXP-008` starts |
| OQ-2 | Print IR unit system | `PRT-001`, `LAY-*` | `PRT-001-001` |
| OQ-3 | Text shaping backend | `LAY-006`, `LAY-007` | ADR when `LAY-007` starts |
| OQ-4 | Remote resource policy | `OPS-002`, `DAT-011` | `OPS-002` |
| OQ-5 | PDF library choice | `OUT-001` | `OUT-001-001` / ADR |
| OQ-6 | `internal/exec` vs `report` orchestration | `FND-001` | `FND-001-001` |

See [`FEATURE_MATRIX.md#open-questions`](FEATURE_MATRIX.md#open-questions).

## Next milestone

**M0 — Baseline and contracts.** Start with
[`tasks/INDEX.md`](tasks/INDEX.md); all M0 tasks are parallel-safe.

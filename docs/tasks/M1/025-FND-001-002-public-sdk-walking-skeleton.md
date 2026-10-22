# FND-001-002: Public SDK walking skeleton

- **Status:** Done
- **Feature:** [FND-001](../../FEATURE_MATRIX.md#fnd--foundation-public-api-cli-infrastructure)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** PRT-001-002, FIL-001-001, OUT-001-001
- **Blocked by:** none
- **Primary owner/package:** `report` (proposed, module root)
- **Likely touched files:** `report/api.go`, `report/api_test.go`
- **Collision risk:** High — the public facade wiring all stages; land after the stage implementations exist.
- **Parallel-safe with:** none on the critical path (it integrates others)
- **Estimated size:** M

## Objective

Wire the M1 stages behind the public API so a caller can compile a JRXML
template, fill it from a slice data source, and export a PDF end-to-end.

## Current state / prerequisites

Parser, fill, layout, print model and PDF exporter are implemented
(`JRX-*`, `FIL-001-001`, `LAY-003-001`, `OUT-001-001`).

## Scope (in)

- Implement `Compile`/`Fill`/`Export` (and an optional `Run` convenience) as
  thin orchestration with `context` and error propagation.
- Keep no engine logic in the facade; it forwards to `internal/*`.
- End-to-end test: template + slice ⇒ PDF bytes non-empty with expected page
  count.

## Out of scope

- CLI commands (`FND-005`); caching (`FND-004`).

## Proposed contracts and integration points

Per `FND-001-001`.

## Implementation checklist (future agent only)

- [x] Implement orchestration with cancellation at boundaries.
- [x] End-to-end test.
- [x] Run `make check`.

## Acceptance criteria (observable)

- A caller obtains a PDF from a small template + slice in one test.
- Cancelled context aborts without panic and closes resources.

## Tests and fixtures to add (future agent only)

- `report/api_test.go` + `testdata/m1/hello.jrxml`.

## Failure cases / edge cases

- Fill without export; export with empty document.

## Definition of done

End-to-end test passes; M1 walking skeleton complete.

## References and downstream dependencies

- Feature: `FND-001`; [`../../ROADMAP.md`](../../ROADMAP.md) M1.
- Unblocks: `FND-005`, `OPS-011`.

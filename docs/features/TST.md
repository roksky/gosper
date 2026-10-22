# TST — Test infrastructure and correctness suites

**Plan:** [`../TEST_STRATEGY.md`](../TEST_STRATEGY.md).

## Purpose

Provide reproducible fixtures, golden/semantic/pagination assertions, fuzzing,
race/leak tests, benchmarks and an isolated differential harness.

## Technical expectations

- Deterministic environment (clock, locale, timezone, fonts) via test helpers.
- Fixtures named by feature ID; goldens updated only via explicit flag.
- Differential tests are isolated and test-only.
- Fuzz targets never panic.

## Features

| ID | Expectation | Acceptance test (planned) |
| --- | --- | --- |
| TST-001 | Reproducible baseline/fixtures. | `make check` / test helpers. |
| TST-002 | Golden harness. | Golden compare + update flag. |
| TST-003 | Semantic assertions. | Structure/value assertions. |
| TST-004 | Pagination assertions. | Page-count/break checks. |
| TST-005 | Visual regression. | Image diff with tolerance. |
| TST-006 | Differential harness. | Isolated reference compare. |
| TST-007 | Fuzzing. | Fuzz corpus, no panics. |
| TST-008 | Property tests. | Property invariants. |
| TST-009 | Race/cancel/leak. | `-race` + leak checks. |
| TST-010 | Benchmarks. | Declared-hardware benchmarks. |
| TST-011 | Cross-platform smoke. | OS matrix build/test. |
| TST-012 | Compatibility corpus. | Corpus linked to `COMPATIBILITY.md`. |

## Tasks

`TST-001-001`, `TST-001-002` (M0); `TST-002-001`, `TST-009-001` (M1).

## Risks

- Flaky visual/differential tests; pin environment and define tolerances.

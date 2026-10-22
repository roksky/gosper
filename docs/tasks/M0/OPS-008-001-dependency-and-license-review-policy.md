# OPS-008-001: Dependency and license review policy with SBOM

- **Status:** Not Started
- **Feature:** [OPS-008](../../FEATURE_MATRIX.md#ops--security-performance-diagnostics-integration)
- **Milestone:** M0
- **Priority:** P1
- **Dependencies:** none
- **Blocked by:** none
- **Primary owner/package:** `docs/` + `.github/`
- **Likely touched files:** `docs/LICENSING.md`, `docs/development.md`, optional `docs/adr/NNNN-dependency-policy.md`, optional CI dependency-check step
- **Collision risk:** Low/Medium — may touch `.github/workflows/ci.yml`; coordinate with any CI task.
- **Parallel-safe with:** `FND-003-001`, `TST-001-002`, `JRX-022-001`
- **Estimated size:** S

## Objective

Establish how new dependencies are justified, reviewed for license
compatibility, and tracked (SBOM), before the engine starts pulling libraries
(notably for PDF/text).

## Current state / prerequisites

No third-party dependencies exist; `go.sum` is absent. License decision is
pending ([`../../LICENSING.md`](../../LICENSING.md)).

## Scope (in)

- A written dependency policy: standard library first; justify each addition;
  license compatibility check with the (pending) project license.
- A decision on SBOM tooling (e.g. `go version -m`, CycloneDX generators) —
  choose or explicitly defer, without adding a dependency prematurely.
- Note interaction with `OPS-009` (JasperReports interoperability review).

## Out of scope

- Choosing the project license itself (`LICENSING.md` remains pending).
- Adding dependencies.

## Proposed contracts and integration points

Policy doc + optional CI step. If a tool is adopted, it must be pinned. If
uncertain, document the manual review process instead.

## Implementation checklist (future agent only)

- [ ] Write the policy (in `docs/LICENSING.md` or a dedicated section).
- [ ] Record the SBOM tooling decision (adopt + pin, or defer with rationale).
- [ ] If adding CI, keep `permissions: contents: read` and no secrets.

## Acceptance criteria (observable)

- Policy explicitly requires license review before merging any new dependency.
- SBOM approach is either pinned/documented or explicitly deferred with a
  reason.

## Tests and fixtures to add (future agent only)

- None (process/documentation). CI step (if any) must pass.

## Failure cases / edge cases

- Tools that reach the network in CI: avoid or make optional.

## Definition of done

Policy documented and discoverable from `docs/README.md`.

## References and downstream dependencies

- Features: `OPS-008`, `OPS-009`; [`../CONTRIBUTING.md`](../../../CONTRIBUTING.md).

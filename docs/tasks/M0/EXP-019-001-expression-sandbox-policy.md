# EXP-019-001: Define the expression sandbox policy

- **Status:** Done
- **Feature:** [EXP-019](../../FEATURE_MATRIX.md#exp--expression-language-and-type-semantics)
- **Milestone:** M0
- **Priority:** P0
- **Dependencies:** none
- **Blocked by:** none
- **Primary owner/package:** `internal/expression` + `docs/`
- **Likely touched files:** `docs/architecture/expressions-and-types.md`, `docs/COMPATIBILITY.md`, `internal/expression/doc.go`
- **Collision risk:** Low — mostly documentation plus a package comment; avoid editing shared `internal/expression` code files owned by `EXP-001-001`.
- **Parallel-safe with:** `FND-003-001`, `PRT-001-001`
- **Estimated size:** S

## Objective

Record the non-negotiable sandbox policy: expressions are evaluated by a closed
interpreter; no arbitrary Go/Java execution, reflection, class loading, method
invocation on arbitrary objects, or I/O side effects.

## Current state / prerequisites

No expression code. Policy only partially stated in docs.

## Scope (in)

- Document the allowed evaluation surface and forbidden operations.
- Define that the function registry is an allowlist (`EXP-015`).
- State that Java-style method translation is a bounded subset (`EXP-014`).
- Add enforcement expectations for tests (malicious expressions rejected).

## Out of scope

- Implementing the evaluator (`EXP-003-001`).
- The function registry contents (`EXP-015`).

## Proposed contracts and integration points

Policy links to future parse-time rejection of unknown constructs. No code
symbols introduced here.

## Implementation checklist (future agent only)

- [x] Expand [`../../architecture/expressions-and-types.md`](../../architecture/expressions-and-types.md)
      with an explicit "Sandbox policy" section.
- [x] Add a note in [`../../COMPATIBILITY.md`](../../COMPATIBILITY.md) that custom
      classes and Groovy/JS are `Unsupported`.
- [x] Reference the policy from `internal/expression/doc.go`.

## Acceptance criteria (observable)

- Docs state the forbidden operations and the allowlist principle.
- The compatibility matrix marks custom classes/Groovy/JS `Unsupported`.

## Tests and fixtures to add (future agent only)

- None at M0 (documentation task). Future `EXP-003-001` adds rejection tests.

## Failure cases / edge cases

- Ambiguity about "safe functions": resolved by `EXP-015`.

## Definition of done

Policy documented and cross-linked; no contradictory statements elsewhere.

## References and downstream dependencies

- Features: `EXP-019`, `EXP-014`, `EXP-015`, `EXP-018`.
- Downstream: `EXP-003-001`, `EXP-015`.

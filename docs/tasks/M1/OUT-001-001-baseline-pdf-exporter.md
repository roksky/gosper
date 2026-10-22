# OUT-001-001: Baseline PDF exporter

- **Status:** Not Started
- **Feature:** [OUT-001](../../FEATURE_MATRIX.md#out--format-exporters)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** PRT-001-002, PRT-002-001
- **Blocked by:** PDF library choice (OQ-5)
- **Primary owner/package:** `internal/exporter`
- **Likely touched files:** `internal/exporter/pdf.go`, `internal/exporter/pdf_test.go`, possibly `go.mod` (if a lib is chosen)
- **Collision risk:** High — the M1 end-to-end exporter; also may touch `go.mod`. Isolate the dependency decision and record it.
- **Parallel-safe with:** `internal/text`, `internal/jrxml` tasks (not with other `go.mod` edits)
- **Estimated size:** L (split if needed: page shell first, text next)

## Objective

Emit a valid multi-page PDF from the print IR: pages with correct sizes, text
runs at absolute positions with resolved fonts/sizes, streaming to an
`io.Writer`.

## Current state / prerequisites

Print model + styles exist. OQ-5 (pure-Go vs dependency) must be decided; if a
dependency is added, follow `OPS-008-001` (license/SBOM) and record an ADR.

## Scope (in)

- PDF page tree, content streams, basic fonts (from `internal/text`).
- Text drawing at `Bounds` with `Style`.
- Streaming write to `io.Writer` (`OUT-009`).
- Tests: output is a parseable PDF with expected page count; contains expected
  text (extract or structural check).

## Out of scope

- Images, vector shapes, embedding, links (`OUT-003`..`OUT-005`).
- Determinism guarantees (separate `OUT-010-001`).

## Proposed contracts and integration points

```go
// Proposed — non-binding.
func PDF(ctx context.Context, doc *printdoc.Document, w io.Writer) error
```

## Implementation checklist (future agent only)

- [ ] Decide PDF approach; record ADR/dependency review if needed.
- [ ] Emit pages + text; stream to writer.
- [ ] Tests validating structure and page count.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- Exported file opens as a valid PDF with N pages.
- Text from the fixture appears in the output.
- No full-document buffering beyond the streaming requirement.

## Tests and fixtures to add (future agent only)

- `testdata/m1/expected.pdf`-style golden plus a structural test in
  `pdf_test.go`.

## Failure cases / edge cases

- Zero-page document; writer error mid-stream; missing font.

## Definition of done

Valid PDF produced/tested; M1 walking skeleton renders end-to-end.

## References and downstream dependencies

- Features: `OUT-001`, `OUT-002`, `OUT-009`; OQ-5.
- Unblocks: `FND-001-002`, `OUT-010-001`, `TST-011`.

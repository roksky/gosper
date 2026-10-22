# PRT — Print document and storage

**Architecture:** [`../architecture/intermediate-representations.md`](../architecture/intermediate-representations.md).

## Purpose

Define the engine-neutral print IR consumed by exporters: pages, positioned
elements, resolved styles, resource references and page metadata.

## Technical expectations

- No dependency on parser/compiler/fill (`package-boundaries.md`).
- Resolved style **values**, not style references (`PRT-002`).
- Resource refs (fonts/images) rather than raw bytes where possible
  (`PRT-003`).
- Streaming/lazy storage to avoid holding huge documents (`PRT-004`); spill to
  disk later (`PRT-005`).
- No public serialization format in early milestones (`PRT-006`, Deferred).

### Proposed contract — non-binding

```go
type Document struct { /* Pages []Page; Metadata Metadata */ }
type Page struct {
    Width, Height Unit
    Elements []Element // text | image | shape, with absolute coords
}
type Element interface{ Bounds() Rect }
```

## Features

| ID | Expectation | Acceptance test (planned) |
| --- | --- | --- |
| PRT-001 | Print document model. | Construct + inspect a page. |
| PRT-002 | Resolved style data. | Style values present. |
| PRT-003 | Resource references. | Font/image ref resolved. |
| PRT-004 | Lazy/streaming storage. | Streamed page iteration. |
| PRT-005 | Spill-to-disk. | Large-doc memory bound. |
| PRT-006 | Serialization (Deferred). | n/a |
| PRT-007 | Page metadata. | Page number/size/count. |

## Tasks

`PRT-001-001` (M0); `PRT-001-002`, `PRT-002-001` (M1).

## Risks

- Unit-system decision (`OQ-2`) must be settled before implementation.

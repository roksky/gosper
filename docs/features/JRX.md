# JRX — JRXML syntax/version compatibility and compilation

**Architecture:** [`../architecture/intermediate-representations.md`](../architecture/intermediate-representations.md).
**Compatibility policy:** [`../COMPATIBILITY.md`](../COMPATIBILITY.md).

## Purpose

Parse JRXML source into a template model, detect dialect/version, normalize, and
compile to an immutable report IR with precise diagnostics.

## Technical expectations

- **Decode** with `encoding/xml` initially; tolerate namespace variants.
- **Version detection** records the declared dialect; unknown versions warn,
  not silently pass (`JRX-001`, `JRX-020`).
- **Diagnostics** carry source spans and codes for unsupported constructs
  (`JRX-019`); unsupported ≠ silently ignored.
- **Immutable compiled IR**: no per-run state; safe to cache/share.
- **Bands/elements** are modeled structurally, not as raw XML.

### Element scope for M1

`staticText`, `textField`, and the core bands; fields/parameters declarations;
page properties. Everything else produces a diagnostic or is deferred.

## Features

| ID | Expectation | Acceptance test (planned) |
| --- | --- | --- |
| JRX-001 | Parse root, detect version. | Fixture parses; version recorded. |
| JRX-002 | Page properties/margins/columns mapped. | Metadata equals fixture values. |
| JRX-003 | Parameters parsed with types/defaults. | Parameter list matches fixture. |
| JRX-004 | Fields parsed with types. | Field list matches fixture. |
| JRX-005 | Variables parsed with expressions/reset. | Variable list matches fixture. |
| JRX-006 | Query element parsed to a query model. | Query string + language captured. |
| JRX-007 | Properties + basic styles parsed. | Style attributes captured. |
| JRX-008 | Core bands assembled in order. | Band order/content matches fixture. |
| JRX-009 | Advanced bands parsed. | background/lastPageFooter/noData present. |
| JRX-010 | staticText/textField parsed. | Element tree matches fixture. |
| JRX-011 | line/rectangle/ellipse parsed. | Shapes match fixture. |
| JRX-012 | Image element parsed. | Image ref captured. |
| JRX-013 | Style inheritance resolved. | Inherited attrs resolve correctly. |
| JRX-014 | Position/stretch/printWhen/evalTime parsed. | Settings match fixture. |
| JRX-015 | Groups + reset parsed. | Group tree matches fixture. |
| JRX-016 | Datasets/subdatasets parsed. | Dataset tree matches fixture. |
| JRX-017 | Subreport element parsed. | Subreport ref captured. |
| JRX-018 | Conditional styles/hyperlinks parsed. | Conditions captured. |
| JRX-019 | Unsupported constructs diagnosed. | Unsupported element ⇒ coded diagnostic. |
| JRX-020 | Version-aware normalization. | 6.x variants normalize equivalently. |
| JRX-021 | `.jasper` loading decision (Deferred). | n/a |
| JRX-022 | Compatibility registry maintained. | Each claim links a fixture. |
| JRX-023 | Font declarations parsed. | Font refs captured. |
| JRX-024 | Components (Deferred to M6). | n/a |

## Tasks

`JRX-001-001`, `JRX-022-001` (M0); `JRX-001-002`, `JRX-002-001`, `JRX-003-001`,
`JRX-004-001`, `JRX-008-001`, `JRX-010-001`, `JRX-014-001`, `JRX-019-001` (M1);
`JRX-020-001` (M2).

## Risks

- Parsing XML into an over-general model. Keep the M1 model minimal.

`JRX-020` normalizes the supported 6.x model in `report.Compile`: absent Java language, canonical dialect, and legacy/newer text stretch attributes. `internal/jrxml/normalize_test.go` compares version fixtures and checks warnings for uncertain cases.

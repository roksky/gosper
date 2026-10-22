# Compatibility

> **These are targets and a tracking register, not claims of current support.**
> Parsing support for a small M1 subset is now implemented and evidenced below;
> all other dimensions (compile/fill/export) remain unimplemented. This
> supersedes the earlier bootstrap document previously at `docs/compatibility.md`
> (merged here). See [`STATUS.md`](STATUS.md) for the evidenced baseline.

## Compatibility policy

1. **Subset first.** Target an explicitly documented subset of JasperReports
   **6.x-style JRXML**. See [`adr/0003`](adr/0003-jrxml-6x-subset-first.md).
2. **No blanket claims.** Compatibility is asserted per construct, per stage,
   with evidence. Percentages are not used.
3. **Four distinct dimensions** are tracked separately:
   - **Syntax acceptance** — the parser accepts the construct.
   - **Semantic correctness** — fill/compile produce correct values.
   - **Visual fidelity** — layout matches a reference within tolerance.
   - **Export-specific** — a construct survives a given exporter.
4. **7.x is separate.** JasperReports 7.x is a distinct later milestone
   (`EXT-008`), not part of the initial target.
5. **`.jasper` binaries are out of scope.** JRXML source only (`JRX-021`).
6. **No JVM in production.** Reference tooling is test-only and isolated
   (`TST-006`), licensed appropriately (`OPS-009`).
7. **Sandboxed expressions.** Expressions are interpreted by a closed
   implementation. Custom classes, Groovy/JavaScript and arbitrary Java
   evaluation are `Unsupported` (`EXP-019`, `EXP-018`); see the sandbox policy
   in [`architecture/expressions-and-types.md`](architecture/expressions-and-types.md).

## Status vocabulary

| Status | Meaning | Evidence required |
| --- | --- | --- |
| `Not Tested` | Not yet exercised. | none |
| `Planned` | Scheduled; no implementation. | feature/task link |
| `Partial` | Works for a documented sub-case. | fixture + note of limits |
| `Verified` | Implementation + test passes against a fixture. | fixture + test path |
| `Unsupported` | Deliberately not supported. | policy note |

**Rule:** a construct may only be `Verified` when a test path exists and
passes. Planned Go features do **not** imply compatibility parity.

## Compatibility matrix

`JRXML construct | Jasper version/dialect | Parse | Compile | Fill | PDF | XLSX | Notes / fixture`

Rows below may be `Verified` for syntax parsing where a fixture and passing
test exist (`JRX-001`/`JRX-002`/`JRX-003`/`JRX-004`/`JRX-008`/`JRX-010`/`JRX-014`);
compile/fill/export columns remain `Planned`.

| Construct | Dialect | Parse | Compile | Fill | PDF | XLSX | Notes / fixture |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `<jasperReport>` root + namespaces | 6.x | Verified | Planned | — | — | — | `JRX-001`; fixture `testdata/jrxml/JRX-001-002-minimal-report.jrxml` |
| Page properties/margins/columns | 6.x | Verified | Planned | — | — | — | `JRX-002`; fixture `testdata/jrxml/JRX-002-001-page-properties.jrxml` |
| `<parameter>` | 6.x | Verified | Planned | Planned | — | — | `JRX-003`; fixture `testdata/jrxml/JRX-003-001-parameters.jrxml` |
| `<field>` | 6.x | Verified | Planned | Planned | — | — | `JRX-004`; fixture `testdata/jrxml/JRX-004-001-fields.jrxml` |
| `<variable>` | 6.x | Planned | Planned | Planned | — | — | `JRX-005` |
| `<queryString>` | 6.x | Planned | Planned | Planned | — | — | `JRX-006`, `DAT-005` |
| `<style>` inheritance | 6.x | Planned | Planned | — | Planned | Planned | `JRX-007`, `JRX-013` |
| Title band | 6.x | Verified | Planned | Planned | Planned | — | `JRX-008`; fixture `testdata/jrxml/JRX-008-001-core-bands.jrxml` |
| Page header/footer | 6.x | Verified | Planned | Planned | Planned | — | `JRX-008`; fixture `testdata/jrxml/JRX-008-001-core-bands.jrxml` |
| Column header/footer | 6.x | Verified | Planned | Planned | Planned | — | `JRX-008`; fixture `testdata/jrxml/JRX-008-001-core-bands.jrxml` |
| Detail band | 6.x | Verified | Planned | Planned | Planned | — | `JRX-008`; fixture `testdata/jrxml/JRX-008-001-core-bands.jrxml` |
| Summary band | 6.x | Verified | Planned | Planned | Planned | — | `JRX-008`; fixture `testdata/jrxml/JRX-008-001-core-bands.jrxml` |
| Background/lastPageFooter/noData | 6.x | Planned | Planned | Planned | Planned | — | `JRX-009` |
| `<staticText>` | 6.x | Verified | Planned | — | Planned | Planned | `JRX-010`; fixture `testdata/jrxml/JRX-010-001-basic-elements.jrxml` |
| `<textField>` | 6.x | Verified | Planned | Planned | Planned | Planned | `JRX-010`; fixture `testdata/jrxml/JRX-010-001-basic-elements.jrxml` |
| 6.x text stretch attribute normalization | 6.0 / 6.20.0 | Verified | Partial | — | — | — | `JRX-020`; fixtures `testdata/jrxml/versions/JRX-020-001-{legacy,modern}.jrxml`; `internal/jrxml/normalize_test.go`; stretch rendering remains deferred |
| `<line>`, `<rectangle>`, `<ellipse>` | 6.x | Planned | Planned | — | Planned | Planned | `JRX-011` |
| `<image>` | 6.x | Planned | Planned | Planned | Planned | — | `JRX-012` |
| Element `printWhenExpression` | 6.x | Verified | Planned | Planned | Planned | Planned | `JRX-014`; fixture `testdata/jrxml/JRX-014-001-element-settings.jrxml` |
| Element `evaluationTime` | 6.x | Verified | Planned | Planned | Planned | — | `JRX-014`; fixture `testdata/jrxml/JRX-014-001-element-settings.jrxml` |
| Groups (headers/footers/reset) | 6.x | Planned | Planned | Planned | Planned | — | `JRX-015` |
| Subdatasets | 6.x | Planned | Planned | Planned | Planned | — | `JRX-016` |
| Subreports | 6.x | Planned | Planned | Planned | Planned | — | `JRX-017` |
| Conditional styles | 6.x | Planned | Planned | — | Planned | Planned | `JRX-018` |
| Hyperlinks | 6.x | Planned | Planned | — | Planned | — | `JRX-018`, `LAY-018` |
| Component elements | 6.x | Unsupported | — | — | — | — | `JRX-024`, `EXT-001..003` |
| `$F{}`/`$P{}`/`$V{}` expressions | 6.x | Planned | Planned | Planned | — | — | `EXP-001` |
| Arithmetic/comparison/boolean | 6.x | Planned | Planned | Planned | — | — | `EXP-003`, `EXP-004` |
| Ternary expressions | 6.x | Planned | Planned | Planned | — | — | `EXP-005` |
| Java method calls | 6.x | Partial | Partial | Partial | — | — | Subset only; `EXP-014` |
| Groovy/JavaScript expressions | 6.x | Unsupported | — | — | — | — | `EXP-018` (Deferred) |
| Custom classes | 6.x | Unsupported | — | — | — | — | `EXP-019` sandbox |
| Compiled `.jasper` binary | any | Unsupported | — | — | — | — | `JRX-021` (Deferred) |
| JRXML 7.x dialect | 7.x | Unsupported | — | — | — | — | `EXT-008` (future) |

## Dimension notes

- **Syntax vs semantics vs visual.** A template can parse successfully
  (`Parse: Verified`) yet render differently (`PDF: Partial`). Track each
  independently.
- **Export format behavior differs.** PDF is strict/visual; XLSX/CSV are
  data/grid-oriented and do not preserve visual layout the same way. See
  [`architecture/exporters.md`](architecture/exporters.md).
- **Environment sensitivity.** Fonts, locale and timezone can change visual
  output; fixtures must pin them ([`TEST_STRATEGY.md`](TEST_STRATEGY.md)).

## Differential testing (future, test-only)

Future compatibility confidence may be established by comparing output against
JasperReports tooling, run in an isolated test environment only. Any such
tooling must be obtained under appropriate license terms. No harness exists
today (`TST-006`, `OPS-009`).

## Updating this matrix

When a construct gains a fixture and passing test, update its status and record
the fixture path. Changes to this file go in the **same PR** as the test that
justifies them. See [`CONTRIBUTING_TASKS.md`](CONTRIBUTING_TASKS.md).

## Fixture registry and naming

Every `Verified` or `Partial` claim must point at a fixture and its test. Full
conventions are in [`../testdata/README.md`](../testdata/README.md); summary:

- **Location:** `testdata/<suite>/`, where `<suite>` groups by area (for example
  `jrxml`, `layout`, `pdf`).
- **Naming:** `<feature-id>-<NNN>-<slug>` (for example
  `jrxml/JRX-001-001-minimal-report.jrxml`), matching the feature/task ID that
  introduces it.
- **Goldens:** expected outputs live under `testdata/<suite>/golden/` and are
  regenerated only through the explicit `-update` test flag, then reviewed
  before committing; see [`TEST_STRATEGY.md`](TEST_STRATEGY.md) and
  [`../testdata/README.md`](../testdata/README.md).
- **`Unsupported` rows** cite the policy and need no fixture.
- **Rule:** a row may only move to `Verified` in the same PR that adds the
  fixture and a passing test; the Notes/fixture column must link that path.
  `Planned`/`Not Tested` rows may name the intended feature/task but not a
  nonexistent fixture.

This registry is maintained by `JRX-022` (task `JRX-022-001`, **Done**).

# Compatibility

> **These are targets and a tracking register, not claims of current support.**
> No compatibility is implemented in the bootstrap/scaffold state. This
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

All rows start `Planned`/`Not Tested`; none are `Verified` today.

| Construct | Dialect | Parse | Compile | Fill | PDF | XLSX | Notes / fixture |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `<jasperReport>` root + namespaces | 6.x | Planned | Planned | — | — | — | `JRX-001`; fixture `testdata/` |
| Page properties/margins/columns | 6.x | Planned | Planned | — | — | — | `JRX-002` |
| `<parameter>` | 6.x | Planned | Planned | Planned | — | — | `JRX-003` |
| `<field>` | 6.x | Planned | Planned | Planned | — | — | `JRX-004` |
| `<variable>` | 6.x | Planned | Planned | Planned | — | — | `JRX-005` |
| `<queryString>` | 6.x | Planned | Planned | Planned | — | — | `JRX-006`, `DAT-005` |
| `<style>` inheritance | 6.x | Planned | Planned | — | Planned | Planned | `JRX-007`, `JRX-013` |
| Title band | 6.x | Planned | Planned | Planned | Planned | — | `JRX-008` |
| Page header/footer | 6.x | Planned | Planned | Planned | Planned | — | `JRX-008` |
| Column header/footer | 6.x | Planned | Planned | Planned | Planned | — | `JRX-008` |
| Detail band | 6.x | Planned | Planned | Planned | Planned | — | `JRX-008` |
| Summary band | 6.x | Planned | Planned | Planned | Planned | — | `JRX-008` |
| Background/lastPageFooter/noData | 6.x | Planned | Planned | Planned | Planned | — | `JRX-009` |
| `<staticText>` | 6.x | Planned | Planned | — | Planned | Planned | `JRX-010` |
| `<textField>` | 6.x | Planned | Planned | Planned | Planned | Planned | `JRX-010` |
| `<line>`, `<rectangle>`, `<ellipse>` | 6.x | Planned | Planned | — | Planned | Planned | `JRX-011` |
| `<image>` | 6.x | Planned | Planned | Planned | Planned | — | `JRX-012` |
| Element `printWhenExpression` | 6.x | Planned | Planned | Planned | Planned | Planned | `JRX-014`, `FIL-008` |
| Element `evaluationTime` | 6.x | Planned | Planned | Planned | Planned | — | `JRX-014` |
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

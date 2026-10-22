# Compatibility

> These are **targets**, not claims of current support. No compatibility is
> implemented in the bootstrap milestone.

## Initial target

- **JRXML 6.x-oriented compatibility.** The first compatibility effort focuses
  on template constructs and expression behavior found in JasperReports 6.x
  templates.

## Explicitly separate or deferred

- **JasperReports 7.x** compatibility is a separate future effort and is not
  covered by the initial target.
- **Java expression behavior** (including JVM type coercion, method invocation
  and library semantics) is a distinct compatibility area. Matching it exactly
  may require deliberate, documented divergences.
- **Compiled `.jasper` binaries** are not part of the initial target. The
  project reads JRXML source templates, not JasperReports' serialized compiled
  form.
- **No JVM runtime** is planned for production use.

## Differential testing (future)

Future compatibility confidence may be established through differential testing
that compares engine output against JasperReports tooling. Any such tooling must
be obtained under appropriate license terms, and any test harness is future work;
none is implemented in this repository today. See `docs/testing.md`.

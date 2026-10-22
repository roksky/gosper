# Test data

Test fixtures live here: input templates and data, plus the expected ("golden")
outputs used by tests. Go's toolchain ignores any directory named `testdata`, so
these files never affect package builds.

This directory currently contains **no fixtures**. The conventions below are
fixed by task `TST-001-002` so future tests are consistent; the comparison
harness itself is `TST-002-001` (not yet implemented).

## Layout

```text
testdata/
└── <suite>/
    ├── <feature-id>-<NNN>-<slug>.jrxml         # input template
    ├── <feature-id>-<NNN>-<slug>.data.json     # optional input data
    ├── <feature-id>-<NNN>-<slug>.data.csv      # optional input data
    └── golden/
        └── <feature-id>-<NNN>-<slug>.<ext>     # expected output
```

### Suites

A `<suite>` is a short, lowercase area name. Use one of:

`jrxml`, `expression`, `datasource`, `fill`, `layout`, `printdoc`, `pdf`,
`xlsx`, `html`, `csv`.

Add a new suite only when no existing one fits; keep names lowercase and
singular for case-insensitive filesystems.

### Naming

- The base name is `<feature-id>-<NNN>-<slug>`, matching the feature/task that
  introduces the fixture — for example `JRX-001-001-minimal-report.jrxml`. The
  same base name is used for its inputs and its golden so they are easy to
  associate.
- `slug` is lowercase, hyphen-separated, and short (e.g. `multi-band`,
  `group-sum`).
- Inputs use their natural extension (`.jrxml`, `.json`, `.csv`).
- Goldens live under the suite's `golden/` directory and keep the base name with
  the output extension (`pdf`, `xlsx`, `html`, `csv`, `txt`).

## Golden update workflow

Golden files are compared byte-for-byte (or semantically, per the test). They
are **never** rewritten silently.

1. Each test package that uses goldens registers an explicit `-update` flag in
   its test files, for example:

   ```go
   var update = flag.Bool("update", false, "rewrite golden files")
   ```

   The shared helper that implements this is `TST-002-001`; until it exists,
   follow this convention when adding a test that needs it.

2. To (re)generate goldens after an intentional behavior change:

   ```sh
   go test ./<package> -run <TestName> -update
   ```

3. **Review before committing.** Inspect the regenerated files with
   `git diff -- testdata/` and `git status`. Golden changes are part of the
   review surface: a PR that changes a golden must explain why the output
   legitimately changed.

4. Run `make check` (and CI) with the new goldens; the comparison path (no
   `-update`) must pass.

If a golden is missing and `-update` was not passed, the test must fail with a
message telling the author to run with `-update` and review the result.

## Rules

- Keep fixtures **small**, **license-clear**, and free of **sensitive data**.
- Do not add fictional fixtures solely to populate directories; add fixtures with
  the tests that use them.
- Prefer deterministic inputs: pin time, timezone, locale and fonts
  ([`../docs/TEST_STRATEGY.md`](../docs/TEST_STRATEGY.md)).
- Lowercase all directory names.

## References

- [`../docs/TEST_STRATEGY.md`](../docs/TEST_STRATEGY.md) — test strategy and layers.
- [`../docs/COMPATIBILITY.md`](../docs/COMPATIBILITY.md) — the compatibility
  fixture registry that links `Verified` claims to fixtures here.
- Task `TST-001-002` (this convention) and `TST-002-001` (the harness).

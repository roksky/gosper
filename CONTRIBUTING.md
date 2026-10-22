# Contributing

Thanks for your interest in gosper. This repository is a bootstrap scaffold;
contributions should advance it in **small, independently reviewable changes**.

## Scope of changes

- Keep pull requests focused on a single concern.
- Engine features are added one stage at a time. Do not combine unrelated
  stages (e.g. parsing and exporting) in one PR.
- Document anything intentionally deferred rather than smuggling it into an
  unrelated change.

## Before you start

- Open or reference an issue describing the intended change.
- Check `docs/ARCHITECTURE.md` and `docs/architecture/package-boundaries.md`
  for the intended package boundaries.
- Confirm the change fits the current milestone; speculative abstractions and
  large dependency additions are discouraged.

## Formatting and tests

- Run `gofmt` (or `make fmt`) before committing.
- `make check` must pass: formatting check, `go vet ./...`, `go test ./...`.
- Add tests for new behavior. Run `make test-race` for concurrency-sensitive
  code.
- CI runs the same checks across Linux, macOS and Windows.

## Public API

- Prefer internal packages until a stable public API is justified by a real
  consumer.
- Treat exported names as a compatibility commitment. Changes to them should
  be deliberate and documented.
- Keep the package dependency graph acyclic and directional.

## Error handling

- Return errors rather than panicking in library code.
- Wrap errors with context (`fmt.Errorf("...: %w", err)`) so callers can
  inspect them.
- When long-running operations are eventually added, they should accept a
  `context.Context` as their first argument and honor cancellation. No such
  APIs exist yet.

## Dependencies

- The standard library is preferred.
- Justify every new dependency, including its license and maintenance status.
- Do not add a dependency solely for a future feature.

## Review standards

- Changes are reviewed for correctness, scope, tests, documentation and
  adherence to the boundaries above.
- Be responsive to review feedback; keep discussions focused on the change.

## Commit and PR hygiene

- Use clear, descriptive commit messages.
- Fill in the pull request template, including the checklist.

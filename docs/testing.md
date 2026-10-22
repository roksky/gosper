# Testing

> This document records testing **strategy**. It clearly separates what runs
> today from what is planned.

## What runs today

- **`go test ./...`** — unit tests for the CLI metadata/help/version behavior.
- **`go build ./cmd/gosper`** — a compile/smoke check exercised in CI.
- **`gofmt` and `go vet`** — enforced in CI and via `make check`.
- **Race detector** — `go test -race ./...` on Linux in CI.

Run them locally with:

```sh
make check
make test-race
make coverage
```

## Planned strategy (not yet implemented)

When the engine is implemented, testing is expected to grow along these lines:

- **Unit tests per module** for parsing, compilation, expressions, data
  sources, fill, layout and exporters.
- **Golden fixtures** in `testdata/` for stable, reviewable expected output,
  updated deliberately.
- **Semantic assertions** that check report structure and values rather than
  brittle byte-for-byte comparisons where appropriate.
- **Pagination checks** that verify page breaks and band layout across page
  sizes and margins.
- **Cross-platform smoke tests** ensuring consistent behavior on Linux, macOS
  and Windows.
- **Java-reference differential tests** comparing output against
  appropriately licensed JasperReports tooling, as described in
  `docs/compatibility.md`. This is future work and is not present today.

## Fixtures and data

- `testdata/` holds input fixtures and expected results for future tests.
- `examples/` holds runnable example programs (none yet).

Neither directory contains fictional templates or results added solely to
populate the folders.

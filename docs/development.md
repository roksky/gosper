# Development

## Requirements

- **Go**: the version recorded in `go.mod` (currently Go 1.26). The `go.mod`
  `go` directive is the single source of truth; CI reads it via
  `go-version-file`.
- **`make`**: optional but recommended; all common tasks are wrapped there.
- **`git`**.
- No Java/JVM is required for development or tests.

## Getting started

```sh
git clone https://github.com/roksky/gosper.git
cd gosper
make check   # gofmt check, go vet, go test
make build   # builds ./bin/gosper
```

## Common commands

| Command | Purpose |
| --- | --- |
| `make fmt` | Format Go sources with `gofmt`. |
| `make fmt-check` | Fail if any Go file is not `gofmt`-clean. |
| `make vet` | Run `go vet ./...`. |
| `make test` | Run `go test ./...`. |
| `make test-race` | Run `go test -race ./...`. |
| `make coverage` | Write `coverage.out` and print a function summary. |
| `make build` | Build the CLI into `bin/gosper` (gitignored). |
| `make check` | `fmt-check` + `vet` + `test`; no network or optional tools. |
| `make lint` | Run `golangci-lint` if installed. |
| `make clean` | Remove local build and coverage outputs. |

CI executes the same logical checks; see `.github/workflows/ci.yml`.

## Formatting and linting

Formatting and `go vet` are the enforced baseline. `golangci-lint` is optional
and not yet configured because its version and config schema were not pinned
for this environment. To use it locally:

```sh
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
make lint
```

A future change may add a pinned `.golangci.yml` once a specific linter version
is selected and verified.

## Package conventions

- Package documentation lives in `doc.go` with a valid package clause.
- `internal/` packages are implementation details; a stable public API is added
  at the module root only when a real consumer exists.
- Keep the dependency graph acyclic and directional (inputs → output).

## Adding tests

- Prefer table-driven tests using the standard `testing` package.
- Keep tests close to the code they exercise.
- Use `internal` packages freely; there are no external test dependencies.

## How CI runs

`.github/workflows/ci.yml` runs on pull requests and pushes to `main`/`dev`:
a `gofmt` check, vet/test/build on Linux, and a separate race-detector job on
Linux. Workflow permissions are read-only. macOS and Windows matrix entries are
temporarily disabled; see the comment in the workflow to re-enable them.

## Changing the module path

The module path is `github.com/roksky/gosper`, derived from the repository's
`origin` remote. If the project moves to a different hosting path, update:

1. the `module` line in `go.mod`,
2. any import paths referencing the old path,
3. links in `README.md` and `docs/`.

Then run `go mod tidy` and `make check`.

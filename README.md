# gosper

A native Go library for a JasperReports-compatible reporting engine. The goal
is to parse JasperReports-style JRXML templates and run the full report
lifecycle — parse → compile → fill → print document → export — entirely in Go,
with no JVM required at runtime.

> ## ⚠️ Scaffold only — report generation is not implemented
>
> This repository is currently a **bootstrap scaffold**. It contains project
> structure, tooling, documentation and a metadata-only CLI. It **cannot parse,
> compile, fill or export any report**, and it exposes no report API. See the
> status table below for what actually works today.

## Future lifecycle

```
JRXML template → compile → fill → print document → export (PDF/XLSX/HTML/CSV)
```

## Requirements

- Go (the version in `go.mod`, currently **Go 1.26**)
- `make` (optional, for convenience targets)
- No Java/JVM required

## Setup and commands

```sh
git clone https://github.com/roksky/gosper.git
cd gosper

make check    # formatting check + go vet + go test
make test     # go test ./...
make build    # build ./bin/gosper
```

## CLI

The CLI currently exposes metadata only:

```console
$ ./bin/gosper --help
gosper - a Go-native, JasperReports-compatible reporting engine.

Scaffold only: report generation is not implemented yet.
...

$ ./bin/gosper --version
gosper dev
```

`gosper help` and `gosper version` behave the same as the flags above. Unknown
commands print an error to stderr and exit with a nonzero status.

## Status

| Area | State |
| --- | --- |
| Repository scaffold, tooling, docs | ✅ Done |
| CLI help/version metadata | ✅ Done |
| JRXML parsing | ⛔ Not implemented |
| Template compilation | ⛔ Not implemented |
| Expression evaluation | ⛔ Not implemented |
| Data sources | ⛔ Not implemented |
| Fill / iteration / grouping | ⛔ Not implemented |
| Layout / pagination | ⛔ Not implemented |
| Exporters (PDF/XLSX/HTML/CSV) | ⛔ Not implemented |

## Documentation

- [Architecture](docs/architecture.md)
- [Development](docs/development.md)
- [Compatibility](docs/compatibility.md)
- [Testing](docs/testing.md)
- [Licensing (pending)](docs/LICENSING.md)
- [Contributing](CONTRIBUTING.md)

## License

No license has been chosen yet. See [docs/LICENSING.md](docs/LICENSING.md).

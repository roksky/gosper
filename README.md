# gosper

A native Go library for a JasperReports-compatible reporting engine. The goal
is to parse JasperReports-style JRXML templates and run the full report
lifecycle — parse → compile → fill → print document → export — entirely in Go,
with no JVM required at runtime.

Gosper compiles JRXML, fills it from caller-provided rows, and exports PDF,
HTML, CSV, XLSX, or Jasper-style print XML. Embedded SQL can be bound with
`CompiledReport.Query`; the application owns database access.

```go
compiled, err := report.Compile(ctx, strings.NewReader(jrxml))
// Handle err after each operation.
query, args, err := compiled.Query(ctx, report.Params{"school_id": schoolID},
    func(n int) string { return fmt.Sprintf("$%d", n) })
// Execute query with args through your application's repository.
filled, err := report.Fill(ctx, compiled, params, report.FromSlice(rows))
err = report.Export(ctx, filled, output, report.FormatPDF)
```

## Native compatibility

Supported features include standard bands, nested frames, text fonts/colors,
shapes, image parameters (`image.Image`), page numbers, report/group Sum and
Count variables, group headers/footers, typed database values, decimal sums,
ternaries, common string methods, and the date/number formatting expressions
used by Medulla templates. Missing JRXML version attributes are accepted for
standard Jasper 6 templates. No Java/JVM is used.

PDF and HTML preserve positioned elements and embed logos. CSV and XLSX export
a text grid; XLSX currently does not preserve drawings, images, or cell styling.
XML exports the print model and is not a complete JasperPrint interchange
implementation. PDF currently uses core Latin fonts with printable ASCII text;
Unicode font embedding remains unsupported.

Compatibility is scoped to the implemented constructs, not arbitrary Jasper
Java/Groovy code. Subreports, charts, crosstabs, dataset runs, referenced
conditional styles, multi-column layout, band splitting/stretching, and complex
variable calculations remain unsupported. Inert, unreferenced dataset/style
declarations do not execute. Unknown constructs produce diagnostics. SQL
`$P{name}` values use bound arguments; `$P!{}` and `$X{}` clauses are rejected.

Regression fixtures in `testdata/medulla` cover receipts, class-list totals and
grouped balances. Medulla's backend tests additionally compile and export all
shipped JRXML templates with synthetic rows. This does not assert visual or
live-database parity with Jasper.

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

CLI report commands are not implemented yet.
...

$ ./bin/gosper --version
gosper dev
```

`gosper help` and `gosper version` behave the same as the flags above. Unknown
commands print an error to stderr and exit with a nonzero status.

## Documentation

Start at the [documentation hub](docs/README.md).

- [Vision](docs/VISION.md)
- [Status](docs/STATUS.md)
- [Roadmap](docs/ROADMAP.md)
- [Architecture](docs/ARCHITECTURE.md)
- [Feature matrix](docs/FEATURE_MATRIX.md)
- [Compatibility](docs/COMPATIBILITY.md)
- [Test strategy](docs/TEST_STRATEGY.md)
- [Development](docs/development.md)
- [Contribution / task tracking](docs/CONTRIBUTING_TASKS.md)
- [Task index](docs/tasks/INDEX.md)
- [Licensing (pending)](docs/LICENSING.md)
- [Contributing](CONTRIBUTING.md)

## License

No license has been chosen yet. See [docs/LICENSING.md](docs/LICENSING.md).

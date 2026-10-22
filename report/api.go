package report

import (
	"context"
	"errors"
	"io"

	"github.com/roksky/gosper/internal/datasource"
)

// ErrNotImplemented is returned by every entry point in this package until the
// lifecycle tasks implement them. Callers can test for it with errors.Is.
var ErrNotImplemented = errors.New("report: not implemented")

// Format identifies an export format.
type Format string

// Supported format identifiers. An empty Format is invalid.
const (
	FormatPDF  Format = "pdf"
	FormatXLSX Format = "xlsx"
	FormatHTML Format = "html"
	FormatCSV  Format = "csv"
)

// String returns the format's identifier.
func (f Format) String() string { return string(f) }

// Option configures Compile, Fill or Export. No options are defined yet; the
// functional-option shape is fixed so later tasks can add With* constructors
// without breaking the API.
type Option func(*config)

type config struct{}

// Params holds report parameter values keyed by parameter name. A nil Params is
// valid and means "no overrides"; parameter defaults apply.
type Params map[string]any

// DataSource, Cursor and Record are the public aliases of the engine's data
// source contract in internal/datasource. They are aliases so external callers
// can implement data sources without importing an internal package.
type (
	// DataSource yields records for one dataset.
	DataSource = datasource.DataSource
	// Cursor iterates the records of an opened DataSource.
	Cursor = datasource.Cursor
	// Record exposes the fields of one iteration by name.
	Record = datasource.Record
)

// CompiledReport is the immutable, reusable result of Compile. It is safe for
// concurrent Fill calls once the lifecycle is implemented. Construct it only
// through Compile; its representation is unexported and opaque.
type CompiledReport struct {
	_ struct{}
}

// FilledReport is the per-run result of Fill. It is not safe for concurrent
// use and must not be shared across goroutines. Construct it only through Fill.
type FilledReport struct {
	_ struct{}
}

// Compile parses and compiles a JRXML template read from r.
//
// Contract only: it currently returns (nil, ErrNotImplemented). The
// implementation lands in FND-001-002 and the parser/compiler tasks.
func Compile(ctx context.Context, r io.Reader, opts ...Option) (*CompiledReport, error) {
	return nil, ErrNotImplemented
}

// Fill executes a compiled report against ds and the given parameters,
// producing a per-run filled report.
//
// Contract only: it currently returns (nil, ErrNotImplemented). The
// implementation lands in FND-001-002 and the fill tasks.
func Fill(ctx context.Context, c *CompiledReport, params Params, ds DataSource, opts ...Option) (*FilledReport, error) {
	return nil, ErrNotImplemented
}

// Export renders a filled report to w in the given format.
//
// Contract only: it currently returns ErrNotImplemented. The implementation
// lands with the exporter tasks.
func Export(ctx context.Context, f *FilledReport, w io.Writer, format Format, opts ...Option) error {
	return ErrNotImplemented
}

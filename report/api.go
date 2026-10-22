package report

import (
	"context"
	"errors"
	"fmt"
	"github.com/roksky/gosper/internal/datasource"
	"github.com/roksky/gosper/internal/exporter"
	"github.com/roksky/gosper/internal/expression"
	"github.com/roksky/gosper/internal/fill"
	"github.com/roksky/gosper/internal/jrxml"
	"github.com/roksky/gosper/internal/layout"
	"github.com/roksky/gosper/internal/limits"
	"github.com/roksky/gosper/internal/printdoc"
	"io"
	"strings"
)

// ErrNotImplemented identifies export formats outside the M1 PDF subset.
var ErrNotImplemented = errors.New("report: not implemented")

// Format identifies an export format.
type Format string

const (
	FormatPDF  Format = "pdf"
	FormatXLSX Format = "xlsx"
	FormatHTML Format = "html"
	FormatCSV  Format = "csv"
)

func (f Format) String() string { return string(f) }

// Option configures Compile, Fill or Export.
type Option func(*config)
type config struct{ budget limits.Budget }

// Budget sets per-run record, page, step, and retained-field-byte limits.
type Budget = limits.Budget

// WithBudget applies resource limits to a Fill run. Zero or negative values
// mean unlimited for that dimension.
func WithBudget(b Budget) Option { return func(c *config) { c.budget = b } }

// Params holds report parameter overrides. A nil map applies defaults.
type Params map[string]any

// Public aliases let callers implement custom data sources.
type (
	DataSource = datasource.DataSource
	Cursor     = datasource.Cursor
	Record     = datasource.Record
)

// FromSlice exposes the built-in slice source to external callers.
func FromSlice[T any](rows []T) DataSource { return datasource.FromSlice(rows) }

// CompiledReport is immutable after Compile and safe for concurrent Fill calls.
type CompiledReport struct {
	fields fill.Template
	layout *layout.DetailSpec
}

// FilledReport owns the per-run print document for export.
type FilledReport struct{ document *printdoc.Document }

// PageCount returns the number of laid-out pages, including zero for empty data.
func (f *FilledReport) PageCount() int {
	if f == nil || f.document == nil {
		return 0
	}
	return f.document.Metadata.PageCount
}

// Compile parses a JRXML template and prepares its supported M1 bands and
// parameter defaults. Diagnostics for unsupported constructs are returned.
func Compile(ctx context.Context, r io.Reader, opts ...Option) (*CompiledReport, error) {
	if ctx == nil {
		return nil, fmt.Errorf("report: nil context")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if r == nil {
		return nil, fmt.Errorf("report: nil reader")
	}
	tmpl, diags, err := jrxml.Parse(r)
	if err != nil {
		return nil, err
	}
	if len(diags) > 0 {
		return nil, diags.Err()
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	spec, err := layout.FromJRXML(tmpl)
	if err != nil {
		return nil, err
	}
	fields := make([]string, len(tmpl.Fields))
	for i, f := range tmpl.Fields {
		fields[i] = f.Name
	}
	defaults := make([]fill.ParameterDefault, len(tmpl.Parameters))
	for i, p := range tmpl.Parameters {
		defaults[i].Name = p.Name
		if src := strings.TrimSpace(p.DefaultExpr); src != "" {
			defaults[i].Default, err = expression.Compile(src)
			if err != nil {
				return nil, fmt.Errorf("parameter %q default: %w", p.Name, err)
			}
		}
	}
	return &CompiledReport{fields: fill.Template{Fields: fields, Parameters: defaults}, layout: spec}, nil
}

// Fill executes the compiled report and lays out a per-run document.
func Fill(ctx context.Context, c *CompiledReport, params Params, ds DataSource, opts ...Option) (*FilledReport, error) {
	if ctx == nil {
		return nil, fmt.Errorf("report: nil context")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if c == nil || c.layout == nil {
		return nil, fmt.Errorf("report: nil compiled report")
	}
	var cfg config
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	overrides := make(map[string]expression.Value, len(params))
	for name, raw := range params {
		value, err := fill.ValueOf(raw)
		if err != nil {
			return nil, fmt.Errorf("parameter %q: %w", name, err)
		}
		overrides[name] = value
	}
	filled, err := fill.RunWithParams(ctx, &c.fields, ds, overrides, cfg.budget)
	if err != nil {
		return nil, err
	}
	doc, err := layout.LayoutWithBudget(ctx, filled, c.layout, cfg.budget)
	if err != nil {
		return nil, err
	}
	return &FilledReport{document: doc}, nil
}

// Export writes a filled report in the requested format. PDF is the M1 format.
func Export(ctx context.Context, f *FilledReport, w io.Writer, format Format, opts ...Option) error {
	if ctx == nil {
		return fmt.Errorf("report: nil context")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if f == nil || f.document == nil {
		return fmt.Errorf("report: nil filled report")
	}
	if format != FormatPDF {
		return fmt.Errorf("report: format %q: %w", format, ErrNotImplemented)
	}
	return exporter.PDF(ctx, f.document, w)
}

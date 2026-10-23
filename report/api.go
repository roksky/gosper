package report

import (
	"context"
	"errors"
	"fmt"
	"github.com/roksky/gosper/internal/datasource"
	"github.com/roksky/gosper/internal/diag"
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

// ErrNotImplemented identifies an unsupported format or query construct.
var ErrNotImplemented = errors.New("report: not implemented")

// Format identifies an export format.
type Format string

const (
	FormatPDF  Format = "pdf"
	FormatXLSX Format = "xlsx"
	FormatHTML Format = "html"
	FormatCSV  Format = "csv"
	FormatXML  Format = "xml"
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
	query  string
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

// Compile parses a JRXML template and prepares supported bands and
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
	for _, d := range diags {
		// Standard Jasper 6 JRXML does not declare a version attribute.
		if d.Code == diag.CodeUnsupported && strings.Contains(d.Message, "unsupported or missing JRXML version") && tmpl.Dialect.Raw == "" {
			continue
		}
		return nil, diags.Err()
	}
	// Normalize supported 6.x syntax before building the layout and fill model.
	// Uncertain normalization cases remain available to direct jrxml callers as
	// warnings; the existing public Compile signature has no diagnostic channel.
	jrxml.Normalize(tmpl)
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	spec, err := layout.FromJRXML(tmpl)
	if err != nil {
		return nil, err
	}
	fieldTypes := map[string]string{}
	for _, field := range tmpl.Fields {
		fieldTypes[field.Name] = field.Type
	}
	var variables []fill.Variable
	seenVariables := map[string]bool{}
	for _, v := range tmpl.Variables {
		if seenVariables[v.Name] {
			return nil, fmt.Errorf("duplicate variable %q", v.Name)
		}
		seenVariables[v.Name] = true
		if v.Calculation != "Nothing" && v.Calculation != "Count" && v.Calculation != "Sum" {
			return nil, fmt.Errorf("unsupported variable calculation %s", v.Calculation)
		}
		if v.ResetType != "Report" && v.ResetType != "Group" && v.ResetType != "None" {
			return nil, fmt.Errorf("unsupported variable reset %s", v.ResetType)
		}
		variable := fill.Variable{Name: v.Name, Calculation: v.Calculation, ResetType: v.ResetType, ResetGroup: v.ResetGroup}
		if strings.TrimSpace(v.Expr) != "" {
			variable.Expr, err = expression.Compile(v.Expr)
			if err != nil {
				return nil, err
			}
		}
		if strings.TrimSpace(v.InitialExpr) != "" {
			variable.Initial, err = expression.Compile(v.InitialExpr)
			if err != nil {
				return nil, err
			}
		}
		variables = append(variables, variable)
	}
	var groups []fill.Group
	seenGroups := map[string]bool{}
	for _, g := range tmpl.Groups {
		if seenGroups[g.Name] {
			return nil, fmt.Errorf("duplicate group %q", g.Name)
		}
		seenGroups[g.Name] = true
		program, err := expression.Compile(g.Expr)
		if err != nil {
			return nil, err
		}
		groups = append(groups, fill.Group{Name: g.Name, Expr: program})
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
	return &CompiledReport{fields: fill.Template{Fields: fields, FieldTypes: fieldTypes, Parameters: defaults, Variables: variables, Groups: groups}, layout: spec, query: tmpl.Query}, nil
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

// Export writes a filled report as PDF, HTML, CSV, XLSX or Jasper print XML.
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
	switch format {
	case FormatPDF:
		return exporter.PDF(ctx, f.document, w)
	case FormatHTML, FormatCSV, FormatXML, FormatXLSX:
		return exportDocument(ctx, f.document, w, format)
	default:
		return fmt.Errorf("report: format %q: %w", format, ErrNotImplemented)
	}
}

// ParameterNames returns the JRXML parameter names in declaration order.
// The returned slice belongs to the caller and can be safely modified.
func (c *CompiledReport) ParameterNames() []string {
	if c == nil {
		return nil
	}
	names := make([]string, len(c.fields.Parameters))
	for i, p := range c.fields.Parameters {
		names[i] = p.Name
	}
	return names
}

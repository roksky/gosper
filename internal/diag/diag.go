package diag

import (
	"errors"
	"fmt"
	"strings"
)

// Code is a stable, machine-readable classification for a diagnostic.
//
// Codes are intentionally open: unknown codes are permitted and must never
// panic, so producers can add codes without breaking consumers.
type Code string

// Initial diagnostic codes.
const (
	CodeParse       Code = "parse"
	CodeUnsupported Code = "unsupported"
	CodeValidation  Code = "validation"
	CodeCanceled    Code = "canceled"
	CodeLimit       Code = "limit"
	CodeInternal    Code = "internal"
)

// Severity classifies a diagnostic.
type Severity int

const (
	// SeverityError indicates the operation did not succeed.
	SeverityError Severity = iota
	// SeverityWarning indicates a recoverable or notable condition.
	SeverityWarning
)

// String returns the lowercase severity name.
func (s Severity) String() string {
	switch s {
	case SeverityError:
		return "error"
	case SeverityWarning:
		return "warning"
	default:
		return fmt.Sprintf("severity(%d)", int(s))
	}
}

// Span locates a diagnostic in a source input. Line and Col are 1-based;
// Offset is a 0-based byte offset. A nil *Span means "no location".
type Span struct {
	Line   int
	Col    int
	Offset int
}

// Diagnostic is a coded, optionally located error message.
//
// Error's format is stable and documented: the code, the severity, the location
// when a Span is present, and the message, for example:
//
//	parse error at line 1, column 5: unexpected token
//	limit error: record budget exceeded
type Diagnostic struct {
	Code     Code
	Severity Severity
	Message  string
	Span     *Span
	// Cause is an optional underlying error, returned by Unwrap.
	Cause error
}

// New builds a Diagnostic. span may be nil.
func New(code Code, sev Severity, msg string, span *Span) *Diagnostic {
	return &Diagnostic{Code: code, Severity: sev, Message: msg, Span: span}
}

// Error implements error. See Diagnostic for the format.
func (d *Diagnostic) Error() string {
	loc := ""
	if d.Span != nil {
		loc = fmt.Sprintf(" at line %d, column %d", d.Span.Line, d.Span.Col)
	}
	return fmt.Sprintf("%s %s%s: %s", d.Code, d.Severity, loc, d.Message)
}

// Unwrap returns the optional underlying cause.
func (d *Diagnostic) Unwrap() error { return d.Cause }

// As extracts the first *Diagnostic from err's chain, if any.
func As(err error) (*Diagnostic, bool) {
	var d *Diagnostic
	if errors.As(err, &d) {
		return d, true
	}
	return nil, false
}

// List aggregates diagnostics in a deterministic order (slice order).
type List []*Diagnostic

// Add appends diagnostics to the list.
func (l *List) Add(ds ...*Diagnostic) {
	*l = append(*l, ds...)
}

// Err returns nil for an empty list, otherwise an error that wraps every
// diagnostic so errors.As/Is can traverse it.
func (l List) Err() error {
	if len(l) == 0 {
		return nil
	}
	return listError(l)
}

// String formats each diagnostic on its own line, in slice order.
func (l List) String() string {
	parts := make([]string, len(l))
	for i, d := range l {
		parts[i] = d.Error()
	}
	return strings.Join(parts, "\n")
}

type listError List

func (e listError) Error() string { return List(e).String() }

// Unwrap exposes each diagnostic so errors.As/Is can traverse the list.
func (e listError) Unwrap() []error {
	out := make([]error, len(e))
	for i, d := range e {
		out[i] = d
	}
	return out
}

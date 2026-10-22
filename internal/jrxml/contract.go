package jrxml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/roksky/gosper/internal/diag"
)

// Dialect identifies the JRXML dialect declared by a template's root element.
//
// Raw holds the version string exactly as written (for example "6.20.0"). Major
// and Minor are parsed from it. A zero Major means the version was absent or
// unparsable; use Known to test.
type Dialect struct {
	Major int
	Minor int
	Raw   string
}

// Known reports whether a usable major version was detected.
func (d Dialect) Known() bool { return d.Major > 0 }

// String returns the raw version string, or "unknown" when none was declared.
func (d Dialect) String() string {
	if d.Raw == "" {
		return "unknown"
	}
	return d.Raw
}

// Template is the top-level parsed representation of a JRXML template.
//
// Only the dialect is populated today; report metadata, parameters, fields,
// bands and elements are added by later tasks.
type Template struct {
	Dialect Dialect
}

// Parse reads a JRXML template from r.
//
// It returns the parsed Template and a list of non-fatal diagnostics. On a fatal
// error (unreadable input, malformed XML, or no root element) the returned
// Template is nil, the returned list contains a diagnostic with code
// diag.CodeParse, and the returned error wraps that diagnostic so callers can
// recover it with errors.As or diag.As. Parse never panics.
func Parse(r io.Reader) (*Template, diag.List, error) {
	if r == nil {
		return fail(diag.CodeInternal, "nil reader", nil)
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return fail(diag.CodeInternal, "reading template: "+err.Error(), nil)
	}

	dec := xml.NewDecoder(bytes.NewReader(data))
	root, ok, err := firstStartElement(dec)
	if err != nil {
		return fail(diag.CodeParse, "malformed XML: "+err.Error(), syntaxSpan(err))
	}
	if !ok {
		return fail(diag.CodeParse, "no root element", &diag.Span{Line: 1, Col: 1})
	}

	// Continue to EOF to confirm the document is well-formed. The root start
	// element has already been consumed.
	if err := drain(dec); err != nil {
		return fail(diag.CodeParse, "malformed XML: "+err.Error(), syntaxSpan(err))
	}

	version, _ := attribute(root, "version")
	tmpl := &Template{Dialect: parseDialect(version)}

	var list diag.List
	if !tmpl.Dialect.Known() {
		list.Add(diag.New(
			diag.CodeUnsupported,
			diag.SeverityWarning,
			fmt.Sprintf("unsupported or missing JRXML version %q", version),
			spanForOffset(data, 0),
		))
	}
	return tmpl, list, nil
}

// firstStartElement returns the first start element token, skipping the XML
// declaration, comments, directives and whitespace. ok is false if the document
// contains no start element.
func firstStartElement(dec *xml.Decoder) (xml.StartElement, bool, error) {
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			return xml.StartElement{}, false, nil
		}
		if err != nil {
			return xml.StartElement{}, false, err
		}
		if se, ok := tok.(xml.StartElement); ok {
			return se, true, nil
		}
	}
}

// drain reads tokens until EOF to validate well-formedness.
func drain(dec *xml.Decoder) error {
	for {
		_, err := dec.Token()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func attribute(se xml.StartElement, name string) (string, bool) {
	for _, a := range se.Attr {
		if a.Name.Local == name {
			return a.Value, true
		}
	}
	return "", false
}

// parseDialect parses a version string such as "6.20.0", "6.0" or "6".
// An empty or unparsable string yields a Dialect with Major 0 (unknown).
func parseDialect(raw string) Dialect {
	d := Dialect{Raw: raw}
	if raw == "" {
		return d
	}
	parts := strings.SplitN(raw, ".", 3)
	major, err := strconv.Atoi(parts[0])
	if err != nil || major <= 0 {
		return Dialect{Raw: raw}
	}
	d.Major = major
	if len(parts) > 1 {
		if minor, err := strconv.Atoi(parts[1]); err == nil {
			d.Minor = minor
		}
	}
	return d
}

// syntaxSpan derives a best-effort span from an encoding/xml error. Line comes
// from *xml.SyntaxError; column and offset are unknown and left at zero.
func syntaxSpan(err error) *diag.Span {
	var se *xml.SyntaxError
	if errors.As(err, &se) {
		return &diag.Span{Line: se.Line}
	}
	return nil
}

// spanForOffset computes a 1-based line/column span for a byte offset into data.
func spanForOffset(data []byte, off int) *diag.Span {
	if off < 0 {
		off = 0
	}
	if off > len(data) {
		off = len(data)
	}
	line, col := 1, 1
	for i := 0; i < off; i++ {
		if data[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return &diag.Span{Line: line, Col: col, Offset: off}
}

// fail builds a fatal result: a nil template, a single-diagnostic list, and an
// error wrapping that diagnostic.
func fail(code diag.Code, msg string, span *diag.Span) (*Template, diag.List, error) {
	var l diag.List
	l.Add(diag.New(code, diag.SeverityError, msg, span))
	return nil, l, l.Err()
}

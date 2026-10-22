package jrxml

import (
	"bytes"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/roksky/gosper/internal/diag"
	"github.com/roksky/gosper/internal/printdoc"
)

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "testdata", "jrxml", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return b
}

func TestParseMinimalReport(t *testing.T) {
	tmpl, list, err := Parse(bytes.NewReader(readFixture(t, "JRX-001-002-minimal-report.jrxml")))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("unexpected diagnostics: %v", list)
	}
	if tmpl.Dialect.Major != 6 || tmpl.Dialect.Minor != 20 || tmpl.Dialect.Raw != "6.20.0" {
		t.Fatalf("dialect = %+v, want 6.20.0", tmpl.Dialect)
	}
	if tmpl.Name != "minimal" {
		t.Fatalf("name = %q, want minimal", tmpl.Name)
	}
	if tmpl.Language != "java" {
		t.Fatalf("language = %q, want java", tmpl.Language)
	}
	if got := tmpl.Page.Width; got != printdoc.FromPoints(595) {
		t.Fatalf("page width = %d, want %d", got, printdoc.FromPoints(595))
	}
	if len(tmpl.Fields) != 1 || tmpl.Fields[0].Name != "name" {
		t.Fatalf("fields = %+v, want one field named name", tmpl.Fields)
	}
}

func TestHandlerRegistrationSeam(t *testing.T) {
	const elemName = "seamTestElement"
	registerRootChild(elemName, func(p *parser, se xml.StartElement) error {
		p.tmpl.Name = "seam-handled"
		return p.skipElement()
	})
	defer delete(rootChildHandlers, elemName)

	tmpl, _, err := Parse(strings.NewReader(
		`<jasperReport version="6.0"><seamTestElement/></jasperReport>`))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if tmpl.Name != "seam-handled" {
		t.Fatalf("handler was not invoked; name = %q", tmpl.Name)
	}
}

func TestParseUnknownRootElementWarns(t *testing.T) {
	tmpl, list, err := Parse(strings.NewReader(`<otherReport version="6.0"></otherReport>`))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if tmpl == nil {
		t.Fatal("template should not be nil")
	}
	if len(list) != 1 {
		t.Fatalf("want one warning, got %v", list)
	}
	if d, ok := diag.As(list.Err()); !ok || d.Code != diag.CodeUnsupported {
		t.Fatalf("warning = %v, want unsupported", list)
	}
}

func TestParseBOM(t *testing.T) {
	src := append([]byte{0xEF, 0xBB, 0xBF},
		[]byte(`<jasperReport version="6.20.0"></jasperReport>`)...)
	tmpl, list, err := Parse(bytes.NewReader(src))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("unexpected diagnostics: %v", list)
	}
	if !tmpl.Dialect.Known() {
		t.Fatalf("dialect should be known: %+v", tmpl.Dialect)
	}
}

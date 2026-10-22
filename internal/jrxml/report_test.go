package jrxml

import (
	"bytes"
	"strings"
	"testing"

	"github.com/roksky/gosper/internal/diag"
	"github.com/roksky/gosper/internal/printdoc"
)

func TestPageProperties(t *testing.T) {
	tmpl, list, err := Parse(bytes.NewReader(readFixture(t, "JRX-002-001-page-properties.jrxml")))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("unexpected diagnostics: %v", list)
	}
	if tmpl.Name != "page-props" || tmpl.Language != "groovy" {
		t.Fatalf("report metadata = (%q, %q)", tmpl.Name, tmpl.Language)
	}

	p := tmpl.Page
	checks := []struct {
		name string
		got  printdoc.Unit
		want float64
	}{
		{"width", p.Width, 595},
		{"height", p.Height, 842},
		{"columnWidth", p.ColumnWidth, 250},
		{"columnSpacing", p.ColumnSpacing, 15},
		{"topMargin", p.MarginTop, 12},
		{"bottomMargin", p.MarginBottom, 13},
		{"leftMargin", p.MarginLeft, 10},
		{"rightMargin", p.MarginRight, 11},
	}
	for _, c := range checks {
		if want := printdoc.FromPoints(c.want); c.got != want {
			t.Errorf("%s = %d, want %d", c.name, c.got, want)
		}
	}
	if p.Orientation != "Landscape" {
		t.Errorf("orientation = %q, want Landscape", p.Orientation)
	}
	if p.ColumnCount != 2 {
		t.Errorf("columnCount = %d, want 2", p.ColumnCount)
	}
}

func TestPageDefaults(t *testing.T) {
	tmpl, _, err := Parse(strings.NewReader(`<jasperReport version="6.20.0"></jasperReport>`))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	p := tmpl.Page
	if p.Width != printdoc.FromPoints(595) || p.Height != printdoc.FromPoints(842) {
		t.Errorf("default size = %dx%d, want 595x842pt", p.Width, p.Height)
	}
	if p.Orientation != "Portrait" {
		t.Errorf("default orientation = %q, want Portrait", p.Orientation)
	}
	if p.ColumnCount != 1 || p.ColumnWidth != printdoc.FromPoints(555) {
		t.Errorf("default column = %d/%d", p.ColumnCount, p.ColumnWidth)
	}
	if p.MarginTop != printdoc.FromPoints(30) || p.MarginBottom != printdoc.FromPoints(30) {
		t.Errorf("default vertical margins = %d/%d", p.MarginTop, p.MarginBottom)
	}
	if p.MarginLeft != printdoc.FromPoints(20) || p.MarginRight != printdoc.FromPoints(20) {
		t.Errorf("default horizontal margins = %d/%d", p.MarginLeft, p.MarginRight)
	}
}

func TestInvalidPageNumberWarns(t *testing.T) {
	_, list, err := Parse(strings.NewReader(
		`<jasperReport version="6.20.0" pageWidth="wide"></jasperReport>`))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("want one warning, got %v", list)
	}
	if d, ok := diag.As(list.Err()); !ok || d.Code != diag.CodeValidation {
		t.Fatalf("warning = %v, want validation", list)
	}
}

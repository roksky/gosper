package jrxml

import (
	"bytes"
	"strings"
	"testing"

	"github.com/roksky/gosper/internal/diag"
)

func TestUnsupportedConstructs(t *testing.T) {
	_, list, err := Parse(bytes.NewReader(readFixture(t, "JRX-019-001-unsupported.jrxml")))
	if err != nil {
		t.Fatalf("Parse should succeed with warnings, got: %v", err)
	}
	if len(list) == 0 {
		t.Fatal("want unsupported diagnostics, got none")
	}

	wantNames := []string{"chart", "mysteryElement", "crosstab"}
	for _, name := range wantNames {
		d := findUnsupported(list, name)
		if d == nil {
			t.Errorf("no unsupported diagnostic for <%s>; got %v", name, list)
			continue
		}
		if d.Span == nil {
			t.Errorf("diagnostic for <%s> has nil span", name)
		}
		if d.Message == "" || !strings.Contains(d.Message, name) {
			t.Errorf("diagnostic for <%s> has message %q", name, d.Message)
		}
	}

	if !strings.Contains(findUnsupported(list, "chart").Message, "EXT-001") {
		t.Errorf("chart diagnostic should point at its feature: %q", findUnsupported(list, "chart").Message)
	}
}

func TestDeferredConstructIsUnsupportedNotParse(t *testing.T) {
	_, list, err := Parse(strings.NewReader(
		`<jasperReport version="6.20.0"><chart/></jasperReport>`))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	d := findUnsupported(list, "chart")
	if d == nil {
		t.Fatalf("want unsupported diagnostic for <chart>, got %v", list)
	}
	if d.Code != diag.CodeUnsupported {
		t.Errorf("code = %q, want unsupported", d.Code)
	}
}

func findUnsupported(list diag.List, name string) *diag.Diagnostic {
	for _, d := range list {
		if d.Code == diag.CodeUnsupported && strings.Contains(d.Message, "<"+name+">") {
			return d
		}
	}
	return nil
}

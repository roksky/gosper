package jrxml

import (
	"bytes"
	"strings"
	"testing"

	"github.com/roksky/gosper/internal/diag"
)

func TestParameters(t *testing.T) {
	tmpl, list, err := Parse(bytes.NewReader(readFixture(t, "JRX-003-001-parameters.jrxml")))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("unexpected diagnostics: %v", list)
	}
	if len(tmpl.Parameters) != 3 {
		t.Fatalf("parameters = %+v, want 3", tmpl.Parameters)
	}

	title := tmpl.Parameters[0]
	if title.Name != "title" || title.Type != "java.lang.String" || title.DefaultExpr != `"Report"` {
		t.Errorf("title = %+v", title)
	}
	count := tmpl.Parameters[1]
	if count.Name != "count" || count.Type != "java.lang.Integer" || count.DefaultExpr != "1" {
		t.Errorf("count = %+v", count)
	}
	untyped := tmpl.Parameters[2]
	if untyped.Name != "untyped" || untyped.Type != defaultParameterType || untyped.DefaultExpr != "" {
		t.Errorf("untyped = %+v", untyped)
	}
}

func TestDuplicateParameterWarns(t *testing.T) {
	src := `<jasperReport version="6.20.0">
		<parameter name="p"/>
		<parameter name="p"/>
	</jasperReport>`
	_, list, err := Parse(strings.NewReader(src))
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

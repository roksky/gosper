package jrxml

import (
	"bytes"
	"testing"
)

func TestFields(t *testing.T) {
	tmpl, list, err := Parse(bytes.NewReader(readFixture(t, "JRX-004-001-fields.jrxml")))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("unexpected diagnostics: %v", list)
	}
	if len(tmpl.Fields) != 2 {
		t.Fatalf("fields = %+v, want 2", tmpl.Fields)
	}
	if got := tmpl.Fields[0]; got.Name != "id" || got.Type != "java.lang.Integer" || got.Description != "" {
		t.Errorf("field[0] = %+v", got)
	}
	if got := tmpl.Fields[1]; got.Name != "name" || got.Type != "java.lang.String" || got.Description != "Full name" {
		t.Errorf("field[1] = %+v", got)
	}
}

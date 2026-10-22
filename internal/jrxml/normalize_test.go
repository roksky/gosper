package jrxml

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/roksky/gosper/internal/diag"
)

func TestNormalizeEquivalentSixXVariants(t *testing.T) {
	var models []*Template
	for _, name := range []string{"JRX-020-001-legacy.jrxml", "JRX-020-001-modern.jrxml"} {
		f, err := os.Open("../../testdata/jrxml/versions/" + name)
		if err != nil {
			t.Fatal(err)
		}
		model, diags, err := Parse(f)
		f.Close()
		if err != nil || len(diags) != 0 {
			t.Fatalf("Parse %s: %v, %v", name, err, diags)
		}
		if warnings := Normalize(model); len(warnings) != 0 {
			t.Fatalf("Normalize %s: %v", name, warnings)
		}
		if again := Normalize(model); len(again) != 0 {
			t.Fatalf("second Normalize %s: %v", name, again)
		}
		models = append(models, model)
	}
	if !reflect.DeepEqual(models[0], models[1]) {
		t.Fatalf("normalized models differ:\n%+v\n%+v", models[0], models[1])
	}
	if !models[0].Bands[0].Elements[0].Stretch || models[0].Language != "java" || models[0].Dialect.Raw != "6" {
		t.Fatalf("canonical model: %+v", models[0])
	}
}

func TestNormalizeWarnings(t *testing.T) {
	for _, tc := range []struct {
		source string
		code   diag.Code
	}{
		{`<jasperReport><detail><band height="10"><textField textAdjust="ScaleFont"/></band></detail></jasperReport>`, diag.CodeUnsupported},
		{`<jasperReport version="7.0"/>`, diag.CodeUnsupported},
		{`<jasperReport version="6.20"><detail><band height="10"><textField textAdjust="CutText" isStretchWithOverflow="true"/></band></detail></jasperReport>`, diag.CodeValidation},
	} {
		model, _, err := Parse(strings.NewReader(tc.source))
		if err != nil {
			t.Fatal(err)
		}
		warnings := Normalize(model)
		found := false
		for _, d := range warnings {
			if d.Code == tc.code && d.Severity == diag.SeverityWarning {
				found = true
			}
		}
		if !found {
			t.Errorf("Normalize(%q) = %v", tc.source, warnings)
		}
	}
}

package jrxml

import (
	"bytes"
	"strings"
	"testing"

	"github.com/roksky/gosper/internal/diag"
	"github.com/roksky/gosper/internal/printdoc"
)

func TestElementSettings(t *testing.T) {
	tmpl, list, err := Parse(bytes.NewReader(readFixture(t, "JRX-014-001-element-settings.jrxml")))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("unexpected diagnostics: %v", list)
	}
	els := tmpl.Bands[0].Elements
	if len(els) != 2 {
		t.Fatalf("elements = %+v, want 2", els)
	}

	st := els[0]
	if st.X != printdoc.FromPoints(1) || st.Y != printdoc.FromPoints(2) {
		t.Errorf("position = (%d, %d)", st.X, st.Y)
	}
	if st.Width != printdoc.FromPoints(30) || st.Height != printdoc.FromPoints(4) {
		t.Errorf("size = (%d, %d)", st.Width, st.Height)
	}
	if !st.Stretch {
		t.Error("stretch = false, want true")
	}
	if st.EvalTime != "Report" {
		t.Errorf("evaluationTime = %q, want Report", st.EvalTime)
	}
	if st.PrintWhenExpr != "$F{show}" {
		t.Errorf("printWhenExpression = %q", st.PrintWhenExpr)
	}

	tf := els[1]
	if tf.EvalTime != "Now" {
		t.Errorf("default evaluationTime = %q, want Now", tf.EvalTime)
	}
	if tf.Stretch {
		t.Error("default stretch = true, want false")
	}
}

func TestUnknownEvaluationTimeWarns(t *testing.T) {
	src := `<jasperReport version="6.20.0">
		<detail><band height="10">
			<textField evaluationTime="Bogus"><reportElement x="0" y="0" width="1" height="1"/></textField>
		</band></detail>
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

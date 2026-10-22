package jrxml

import (
	"bytes"
	"strings"
	"testing"

	"github.com/roksky/gosper/internal/printdoc"
)

func TestBasicElements(t *testing.T) {
	tmpl, list, err := Parse(bytes.NewReader(readFixture(t, "JRX-010-001-basic-elements.jrxml")))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("unexpected diagnostics: %v", list)
	}
	if len(tmpl.Bands) != 1 {
		t.Fatalf("bands = %+v, want 1", tmpl.Bands)
	}
	band := tmpl.Bands[0]
	if band.Kind != BandDetail {
		t.Fatalf("band kind = %q, want detail", band.Kind)
	}
	if band.Height != printdoc.FromPoints(30) {
		t.Errorf("band height = %d, want %d", band.Height, printdoc.FromPoints(30))
	}
	if len(band.Elements) != 2 {
		t.Fatalf("elements = %+v, want 2", band.Elements)
	}

	st := band.Elements[0]
	if st.Kind != KindStaticText {
		t.Errorf("element[0].Kind = %q", st.Kind)
	}
	if st.Text != "Hello, world" {
		t.Errorf("static text = %q, want %q", st.Text, "Hello, world")
	}
	if st.X != 0 || st.Y != 0 || st.Width != printdoc.FromPoints(100) || st.Height != printdoc.FromPoints(20) {
		t.Errorf("static text geometry = %+v", st.elementAttrs)
	}

	tf := band.Elements[1]
	if tf.Kind != KindTextField {
		t.Errorf("element[1].Kind = %q", tf.Kind)
	}
	if tf.Expr != "$F{amount}" {
		t.Errorf("expression = %q, want $F{amount}", tf.Expr)
	}
	if tf.Pattern != "#,##0.00" {
		t.Errorf("pattern = %q", tf.Pattern)
	}
	if !tf.IsBlankWhenNull {
		t.Error("isBlankWhenNull = false, want true")
	}
	if tf.X != printdoc.FromPoints(100) || tf.Width != printdoc.FromPoints(80) {
		t.Errorf("textField geometry = %+v", tf.elementAttrs)
	}
}

func TestStaticTextVerbatim(t *testing.T) {
	src := `<jasperReport version="6.20.0">
		<detail><band height="10">
			<staticText><reportElement x="0" y="0" width="1" height="1"/><text> a b </text></staticText>
		</band></detail>
	</jasperReport>`
	tmpl, _, err := Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got := tmpl.Bands[0].Elements[0].Text
	if got != " a b " {
		t.Errorf("text = %q, want %q (verbatim)", got, " a b ")
	}
}

func TestElementUnknownChildWarns(t *testing.T) {
	src := `<jasperReport version="6.20.0">
		<detail><band height="10">
			<staticText><reportElement x="0" y="0" width="1" height="1"/><text>x</text><bogus/></staticText>
		</band></detail>
	</jasperReport>`
	_, list, err := Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("want one warning, got %v", list)
	}
}

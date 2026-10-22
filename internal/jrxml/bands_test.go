package jrxml

import (
	"bytes"
	"strings"
	"testing"

	"github.com/roksky/gosper/internal/diag"
	"github.com/roksky/gosper/internal/printdoc"
)

func TestCoreBands(t *testing.T) {
	tmpl, list, err := Parse(bytes.NewReader(readFixture(t, "JRX-008-001-core-bands.jrxml")))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("unexpected diagnostics: %v", list)
	}

	wantKinds := []BandKind{
		BandTitle, BandPageHeader, BandColumnHeader, BandDetail,
		BandColumnFooter, BandPageFooter, BandSummary,
	}
	wantHeights := []float64{40, 20, 15, 25, 15, 20, 30}

	if len(tmpl.Bands) != len(wantKinds) {
		t.Fatalf("bands = %+v, want %d", tmpl.Bands, len(wantKinds))
	}
	for i, kind := range wantKinds {
		b := tmpl.Bands[i]
		if b.Kind != kind {
			t.Errorf("band[%d].Kind = %q, want %q", i, b.Kind, kind)
		}
		if want := printdoc.FromPoints(wantHeights[i]); b.Height != want {
			t.Errorf("band %q height = %d, want %d", kind, b.Height, want)
		}
	}

	detail := tmpl.Bands[3]
	if len(detail.Elements) != 1 || detail.Elements[0].Kind != KindTextField {
		t.Errorf("detail elements = %+v, want one textField", detail.Elements)
	}
}

func TestDuplicateBandWarns(t *testing.T) {
	src := `<jasperReport version="6.20.0">
		<detail><band height="10"/></detail>
		<detail><band height="20"/></detail>
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

func TestMissingDetailAllowed(t *testing.T) {
	tmpl, list, err := Parse(strings.NewReader(
		`<jasperReport version="6.20.0"><title><band height="10"/></title></jasperReport>`))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("unexpected diagnostics: %v", list)
	}
	if len(tmpl.Bands) != 1 || tmpl.Bands[0].Kind != BandTitle {
		t.Fatalf("bands = %+v", tmpl.Bands)
	}
}

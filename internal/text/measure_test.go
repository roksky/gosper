package text

import (
	"github.com/roksky/gosper/internal/printdoc"
	"strings"
	"testing"
)

func TestMeasureFixedStrings(t *testing.T) {
	f, _ := Resolve("SansSerif", 0, 10)
	m, err := Measure("Mi ", f)
	if err != nil {
		t.Fatal(err)
	}
	if m.Width != 278 || m.Ascent != 160 || m.Descent != 40 || m.LineHeight != 240 {
		t.Fatalf("metrics=%+v", m)
	}
	empty, err := Measure("", f)
	if err != nil || empty.Width != 0 || empty.LineHeight != 240 {
		t.Fatalf("empty=%+v,%v", empty, err)
	}
	mono, _ := Resolve("Monospaced", 0, 10)
	m, _ = Measure("abc", mono)
	if m.Width != 360 {
		t.Fatalf("mono width=%d", m.Width)
	}
}
func TestMeasureDeterminismAndScale(t *testing.T) {
	small, _ := Resolve("Serif", Italic, 8)
	large, _ := Resolve("Serif", Italic, 16)
	first, _ := Measure("Invoice 123", small)
	for i := 0; i < 100; i++ {
		again, _ := Measure("Invoice 123", small)
		if again != first {
			t.Fatal("measurement drift")
		}
	}
	bigger, _ := Measure("Invoice 123", large)
	if bigger.Width <= first.Width || bigger.LineHeight <= first.LineHeight {
		t.Fatalf("sizes: %+v %+v", first, bigger)
	}
}
func TestMeasureBoundsAndInvalid(t *testing.T) {
	f, _ := Resolve("SansSerif", 0, 12)
	for _, s := range []string{"a\nb", string([]byte{0xff}), strings.Repeat("a", maxMeasureBytes+1)} {
		if _, err := Measure(s, f); err == nil {
			t.Errorf("accepted invalid text of length %d", len(s))
		}
	}
	if _, err := Measure("x", Font{}); err == nil {
		t.Fatal("accepted zero font")
	}
	m, err := Measure(strings.Repeat("a", 100000), f)
	if err != nil || m.Width <= printdoc.Unit(0) {
		t.Fatalf("long text=%+v,%v", m, err)
	}
}

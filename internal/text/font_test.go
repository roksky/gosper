package text

import "testing"

func TestResolveCoreFontsAndFallback(t *testing.T) {
	for _, tc := range []struct {
		name, want string
		advance    int
	}{{"SansSerif", "SansSerif", 833}, {"serif", "Serif", 889}, {"Monospaced", "Monospaced", 600}, {"unknown", "SansSerif", 833}, {"", "SansSerif", 833}} {
		f, err := Resolve(tc.name, Bold, 12)
		if err != nil {
			t.Fatal(err)
		}
		if f.Name != tc.want || f.GlyphAdvance('M') != tc.advance {
			t.Errorf("%q -> %s, %d", tc.name, f.Name, f.GlyphAdvance('M'))
		}
	}
}
func TestResolveInvalidSizeAndStyle(t *testing.T) {
	for _, size := range []float64{0, -1} {
		if _, err := Resolve("SansSerif", 0, size); err == nil {
			t.Errorf("accepted size %v", size)
		}
	}
	if _, err := Resolve("SansSerif", Style(8), 12); err == nil {
		t.Fatal("accepted unknown style")
	}
}
func TestMetricTableStable(t *testing.T) {
	f, _ := Resolve("SansSerif", 0, 10)
	if got := f.GlyphAdvance('M') + f.GlyphAdvance('i') + f.GlyphAdvance(' '); got != 1389 {
		t.Fatalf("Mi space width=%d", got)
	}
	if got := f.GlyphAdvance('界'); got != 600 {
		t.Fatalf("fallback=%d", got)
	}
}

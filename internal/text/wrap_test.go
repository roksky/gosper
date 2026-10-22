package text

import (
	"reflect"
	"testing"

	"github.com/roksky/gosper/internal/printdoc"
)

func TestWrapGreedyAndLongWord(t *testing.T) {
	font, _ := Resolve("Monospaced", 0, 10)
	width := printdoc.FromPoints(18) // three glyphs
	for _, tc := range []struct {
		input string
		want  []string
	}{
		{"one two", []string{"one", "two"}},
		{"abcde", []string{"abc", "de"}},
		{"a  b\n\nc", []string{"a b", "", "c"}},
		{"", []string{""}},
	} {
		got, err := Wrap(tc.input, font, width)
		if err != nil || !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Wrap(%q) = %q, %v, want %q", tc.input, got, err, tc.want)
		}
		for _, line := range got {
			m, _ := Measure(line, font)
			if m.Width > width {
				t.Errorf("overwide line %q", line)
			}
		}
	}
	if _, err := Wrap("A", font, 1); err == nil {
		t.Fatal("accepted glyph wider than box")
	}
}

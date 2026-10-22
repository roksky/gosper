package expression

import (
	"github.com/roksky/gosper/internal/diag"
	"strings"
	"testing"
)

func TestParseErrorsHavePositions(t *testing.T) {
	for _, src := range []string{"1 +", "(1 + 2", "1 + @", "1\n+ * 2", "foo()", "a.b", "$F{a{b}}", strings.Repeat("(", 130) + "1" + strings.Repeat(")", 130)} {
		_, err := Compile(src)
		if err == nil {
			t.Errorf("%q unexpectedly compiled", src)
			continue
		}
		d, ok := diag.As(err)
		if !ok || d.Span == nil || d.Span.Col < 1 {
			t.Errorf("%q: missing position: %v", src, err)
		}
	}
}
func TestSandboxRejectsCallsAndAccess(t *testing.T) {
	for _, src := range []string{"exec(1)", "system.exit(0)", "a[0]", "a; b", "a = 1"} {
		if _, err := Compile(src); err == nil {
			t.Errorf("accepted %q", src)
		}
	}
}
func FuzzCompile(f *testing.F) {
	for _, s := range []string{"1+2*3", "$F{name}", "not (true and false)", "1/0", "\xff", strings.Repeat("(", 200)} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) { _, _ = Compile(s) })
}

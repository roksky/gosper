package jrxml

import (
	"strings"
	"testing"

	"github.com/roksky/gosper/internal/diag"
)

func TestParseVersion(t *testing.T) {
	cases := []struct {
		name  string
		src   string
		major int
		minor int
		raw   string
	}{
		{
			name:  "full version",
			src:   `<jasperReport version="6.20.0"><field name="x"/></jasperReport>`,
			major: 6,
			minor: 20,
			raw:   "6.20.0",
		},
		{
			name:  "major.minor",
			src:   `<jasperReport version="6.0"/>`,
			major: 6,
			minor: 0,
			raw:   "6.0",
		},
		{
			name:  "major only",
			src:   `<jasperReport version="6"/>`,
			major: 6,
			minor: 0,
			raw:   "6",
		},
		{
			name:  "leading xml declaration and namespace",
			src:   `<?xml version="1.0" encoding="UTF-8"?><jasperReport xmlns="http://jasperreports.sourceforge.net/jasperreports" xmlns:j="http://x" version="6.4.1"></jasperReport>`,
			major: 6,
			minor: 4,
			raw:   "6.4.1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmpl, list, err := Parse(strings.NewReader(tc.src))
			if err != nil {
				t.Fatalf("Parse error = %v, want nil", err)
			}
			if tmpl == nil {
				t.Fatal("Parse returned nil template")
			}
			if len(list) != 0 {
				t.Fatalf("unexpected diagnostics: %v", list)
			}
			if !tmpl.Dialect.Known() {
				t.Fatalf("dialect should be known: %+v", tmpl.Dialect)
			}
			if tmpl.Dialect.Major != tc.major || tmpl.Dialect.Minor != tc.minor || tmpl.Dialect.Raw != tc.raw {
				t.Fatalf("dialect = %+v, want {%d %d %q}", tmpl.Dialect, tc.major, tc.minor, tc.raw)
			}
			if got := tmpl.Dialect.String(); got != tc.raw {
				t.Fatalf("Dialect.String() = %q, want %q", got, tc.raw)
			}
		})
	}
}

func TestParseMissingVersionWarns(t *testing.T) {
	tmpl, list, err := Parse(strings.NewReader(`<jasperReport></jasperReport>`))
	if err != nil {
		t.Fatalf("Parse error = %v, want nil (warning only)", err)
	}
	if tmpl == nil {
		t.Fatal("Parse returned nil template despite a warning, not an error")
	}
	if tmpl.Dialect.Known() {
		t.Fatalf("dialect should be unknown, got %+v", tmpl.Dialect)
	}
	if got := tmpl.Dialect.String(); got != "unknown" {
		t.Fatalf("Dialect.String() = %q, want unknown", got)
	}
	if len(list) != 1 {
		t.Fatalf("want exactly one warning, got %v", list)
	}
	d, ok := diag.As(list.Err())
	if !ok {
		t.Fatal("warning is not a diagnostic")
	}
	if d.Code != diag.CodeUnsupported || d.Severity != diag.SeverityWarning {
		t.Fatalf("warning = %+v, want unsupported/warning", d)
	}
}

func TestParseMalformed(t *testing.T) {
	cases := map[string]string{
		"mismatched tag": `<jasperReport><field></jasperReport>`,
		"unclosed":       `<jasperReport><field>`,
		"garbage":        `<jasperReport <<>`,
	}
	for name, src := range cases {
		t.Run(name, func(t *testing.T) {
			tmpl, list, err := Parse(strings.NewReader(src))
			if err == nil {
				t.Fatal("Parse should fail on malformed XML")
			}
			if tmpl != nil {
				t.Fatalf("template should be nil, got %+v", tmpl)
			}
			d, ok := diag.As(err)
			if !ok {
				t.Fatalf("error is not a diagnostic: %v", err)
			}
			if d.Code != diag.CodeParse {
				t.Fatalf("code = %q, want %q", d.Code, diag.CodeParse)
			}
			if len(list) != 1 || list[0].Code != diag.CodeParse {
				t.Fatalf("list = %v, want one parse diagnostic", list)
			}
		})
	}
}

func TestParseEmpty(t *testing.T) {
	for _, src := range []string{"", "   \n\t", "<!-- comment only -->"} {
		tmpl, _, err := Parse(strings.NewReader(src))
		if err == nil || tmpl != nil {
			t.Fatalf("Parse(%q) = (%v, %v), want nil template and error", src, tmpl, err)
		}
		if d, ok := diag.As(err); !ok || d.Code != diag.CodeParse {
			t.Fatalf("Parse(%q) error = %v, want parse diagnostic", src, err)
		}
	}
}

func TestParseDialectUnknownVersion(t *testing.T) {
	tmpl, list, err := Parse(strings.NewReader(`<jasperReport version="banana"/>`))
	if err != nil {
		t.Fatalf("Parse error = %v, want nil", err)
	}
	if tmpl.Dialect.Known() {
		t.Fatalf("dialect should be unknown: %+v", tmpl.Dialect)
	}
	if len(list) != 1 {
		t.Fatalf("want one warning, got %v", list)
	}
}

func TestParseNilReader(t *testing.T) {
	// A nil reader must not panic; io.ReadAll returns an error.
	_, _, err := Parse(nil)
	if err == nil {
		t.Fatal("Parse(nil) should return an error")
	}
}

func TestParseNoPanic(t *testing.T) {
	inputs := []string{
		"<", "<>", "</>", "<a", "<a b='c'", "<?xml?>", "<a><b/></a>extra",
		"\xef\xbb\xbf<jasperReport version=\"6.0\"></jasperReport>", // BOM
	}
	for _, in := range inputs {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("Parse(%q) panicked: %v", in, r)
				}
			}()
			_, _, _ = Parse(strings.NewReader(in))
		}()
	}
}

func FuzzParse(f *testing.F) {
	for _, seed := range []string{
		`<jasperReport version="6.20.0"></jasperReport>`,
		`<jasperReport>`,
		``,
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, src string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Parse panicked on %q: %v", src, r)
			}
		}()
		_, _, _ = Parse(strings.NewReader(src))
	})
}

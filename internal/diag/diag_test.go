package diag

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestCodeOpenSet(t *testing.T) {
	// Unknown codes are permitted and must not panic.
	custom := Code("custom-code")
	d := New(custom, SeverityError, "boom", nil)
	if got := d.Error(); !strings.Contains(got, "custom-code") {
		t.Fatalf("Error() = %q, want it to contain the code", got)
	}
}

func TestDiagnosticErrorFormatting(t *testing.T) {
	cases := []struct {
		name string
		diag *Diagnostic
		want string
	}{
		{
			name: "with span",
			diag: New(CodeParse, SeverityError, "unexpected token", &Span{Line: 1, Col: 5, Offset: 4}),
			want: "parse error at line 1, column 5: unexpected token",
		},
		{
			name: "without span",
			diag: New(CodeLimit, SeverityError, "record budget exceeded", nil),
			want: "limit error: record budget exceeded",
		},
		{
			name: "warning",
			diag: New(CodeUnsupported, SeverityWarning, "ignored", &Span{Line: 2, Col: 3}),
			want: "unsupported warning at line 2, column 3: ignored",
		},
		{
			name: "empty message",
			diag: New(CodeInternal, SeverityError, "", nil),
			want: "internal error: ",
		},
		{
			name: "nil span does not panic",
			diag: New(CodeValidation, SeverityError, "bad", nil),
			want: "validation error: bad",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.diag.Error(); got != tc.want {
				t.Fatalf("Error() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestAsThroughWrap(t *testing.T) {
	orig := New(CodeParse, SeverityError, "bad", &Span{Line: 3, Col: 7, Offset: 20})
	wrapped := fmt.Errorf("compiling: %w", orig)

	got, ok := As(wrapped)
	if !ok {
		t.Fatal("As() did not find the diagnostic")
	}
	if got.Code != CodeParse {
		t.Fatalf("Code = %q, want %q", got.Code, CodeParse)
	}
	if got.Span == nil || got.Span.Line != 3 || got.Span.Col != 7 {
		t.Fatalf("Span = %+v, want line 3 col 7", got.Span)
	}
}

func TestAsMiss(t *testing.T) {
	if _, ok := As(errors.New("plain")); ok {
		t.Fatal("As() found a diagnostic in a plain error")
	}
}

func TestDiagnosticCause(t *testing.T) {
	cause := errors.New("io")
	d := New(CodeInternal, SeverityError, "read failed", nil)
	d.Cause = cause
	if !errors.Is(d, cause) {
		t.Fatal("Diagnostic does not unwrap to its cause")
	}
}

func TestListErr(t *testing.T) {
	var empty List
	if empty.Err() != nil {
		t.Fatal("empty List.Err() must be nil")
	}

	l := List{
		New(CodeParse, SeverityError, "first", nil),
		New(CodeValidation, SeverityError, "second", nil),
	}
	err := l.Err()
	if err == nil {
		t.Fatal("non-empty List.Err() must be non-nil")
	}

	// errors.As must traverse the aggregated list.
	if _, ok := As(err); !ok {
		t.Fatal("As() did not traverse List error")
	}

	want := "parse error: first\nvalidation error: second"
	if got := err.Error(); got != want {
		t.Fatalf("List error = %q, want %q", got, want)
	}
}

func TestListDeterministicOrder(t *testing.T) {
	a := New(CodeParse, SeverityError, "a", nil)
	b := New(CodeParse, SeverityError, "b", nil)

	var l1 List
	l1.Add(a, b)
	var l2 List
	l2.Add(a, b)

	if l1.String() != l2.String() {
		t.Fatalf("List.String not deterministic: %q vs %q", l1.String(), l2.String())
	}
	if l1.String() != "parse error: a\nparse error: b" {
		t.Fatalf("List.String = %q", l1.String())
	}

	var reversed List
	reversed.Add(b, a)
	if reversed.String() == l1.String() {
		t.Fatal("slice order should be preserved")
	}
}

func TestSeverityString(t *testing.T) {
	if got := SeverityError.String(); got != "error" {
		t.Fatalf("SeverityError = %q", got)
	}
	if got := SeverityWarning.String(); got != "warning" {
		t.Fatalf("SeverityWarning = %q", got)
	}
	if got := Severity(9).String(); got != "severity(9)" {
		t.Fatalf("unknown severity = %q", got)
	}
}

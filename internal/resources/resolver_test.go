package resources

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/roksky/gosper/internal/diag"
)

var _ Resolver = (*fsResolver)(nil)

func testFS() fstest.MapFS {
	return fstest.MapFS{
		"templates/basic.jrxml": &fstest.MapFile{Data: []byte("<jasperReport/>")},
		"secret.txt":            &fstest.MapFile{Data: []byte("top secret")},
	}
}

func readAll(t *testing.T, r Resolver, ref string) (string, error) {
	t.Helper()
	rc, err := r.Open(context.Background(), ref)
	if err != nil {
		return "", err
	}
	defer rc.Close()
	b, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	return string(b), nil
}

func TestOpenAllowedLocalRead(t *testing.T) {
	r := New(testFS())
	got, err := readAll(t, r, "templates/basic.jrxml")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if want := "<jasperReport/>"; got != want {
		t.Fatalf("content = %q, want %q", got, want)
	}
}

func TestDeniedSchemesByDefault(t *testing.T) {
	r := New(testFS())
	refs := []string{
		"file:///etc/passwd",
		"file://templates/basic.jrxml",
		"http://example.com/template.jrxml",
		"https://example.com/template.jrxml",
	}
	for _, ref := range refs {
		if _, err := r.Open(context.Background(), ref); err == nil {
			t.Fatalf("Open(%q) = nil error, want unsupported", ref)
		} else if d, ok := diag.As(err); !ok || d.Code != diag.CodeUnsupported {
			t.Fatalf("Open(%q) code = %v, want %q", ref, d, diag.CodeUnsupported)
		}
	}
}

func TestDeniedTraversal(t *testing.T) {
	r := New(testFS())
	refs := []string{
		"../secret.txt",
		"templates/../../secret.txt",
		"a/../b",
		"/etc/passwd",
	}
	for _, ref := range refs {
		if _, err := r.Open(context.Background(), ref); err == nil {
			t.Fatalf("Open(%q) = nil error, want rejection", ref)
		} else if d, ok := diag.As(err); !ok || d.Code != diag.CodeValidation {
			t.Fatalf("Open(%q) code = %v, want %q", ref, d, diag.CodeValidation)
		}
	}
}

func TestAllowedFileScheme(t *testing.T) {
	r := New(testFS(), WithAllowedSchemes("file"))
	got, err := readAll(t, r, "file://templates/basic.jrxml")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if want := "<jasperReport/>"; got != want {
		t.Fatalf("content = %q, want %q", got, want)
	}
}

func TestAllowedNonFileSchemeHasNoTransport(t *testing.T) {
	r := New(testFS(), WithAllowedSchemes("http"))
	_, err := r.Open(context.Background(), "http://example.com/t.jrxml")
	if d, ok := diag.As(err); !ok || d.Code != diag.CodeUnsupported {
		t.Fatalf("code = %v, want %q", d, diag.CodeUnsupported)
	}
}

func TestAllowedFileSchemeStillRooted(t *testing.T) {
	r := New(testFS(), WithAllowedSchemes("file"))
	if _, err := r.Open(context.Background(), "file://../../secret.txt"); err == nil {
		t.Fatal("Open escaped root, want rejection")
	} else if d, ok := diag.As(err); !ok || d.Code != diag.CodeValidation {
		t.Fatalf("code = %v, want %q", d, diag.CodeValidation)
	}
}

func TestEmptyReference(t *testing.T) {
	r := New(testFS())
	if _, err := r.Open(context.Background(), "   "); err == nil {
		t.Fatal("Open(empty) = nil error, want rejection")
	} else if d, ok := diag.As(err); !ok || d.Code != diag.CodeValidation {
		t.Fatalf("code = %v, want %q", d, diag.CodeValidation)
	}
}

func TestMissingFile(t *testing.T) {
	r := New(testFS())
	_, err := r.Open(context.Background(), "templates/missing.jrxml")
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Open error = %v, want fs.ErrNotExist", err)
	}
}

func TestOpenCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	r := New(testFS())
	if _, err := r.Open(ctx, "templates/basic.jrxml"); !errors.Is(err, context.Canceled) {
		t.Fatalf("Open with canceled ctx = %v, want context.Canceled", err)
	}
}

func TestFromDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "t.jrxml"), []byte("x"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	r := FromDir(dir)
	if _, err := r.Open(context.Background(), "t.jrxml"); err != nil {
		t.Fatalf("Open: %v", err)
	}
	if _, err := r.Open(context.Background(), "../outside.jrxml"); err == nil {
		t.Fatal("Open escaped root, want rejection")
	}
}

func TestSchemeOf(t *testing.T) {
	tests := []struct {
		ref    string
		scheme string
		ok     bool
	}{
		{"templates/a.jrxml", "", false},
		{"file:///etc/passwd", "file", true},
		{"HTTP://example.com", "http", true},
		{"C:/reports/a.jrxml", "", false},
	}
	for _, tc := range tests {
		scheme, ok := schemeOf(tc.ref)
		if scheme != tc.scheme || ok != tc.ok {
			t.Errorf("schemeOf(%q) = (%q, %v), want (%q, %v)", tc.ref, scheme, ok, tc.scheme, tc.ok)
		}
	}
}

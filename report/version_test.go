package report

import "testing"

func TestVersionDefaultAndOverride(t *testing.T) {
	original := version
	t.Cleanup(func() { version = original })
	if Version() == "" {
		t.Fatal("empty default")
	}
	version = "1.2.3"
	if Version() != "1.2.3" {
		t.Fatalf("override=%q", Version())
	}
	version = ""
	if Version() != "dev" {
		t.Fatalf("empty fallback=%q", Version())
	}
	version = "  "
	if Version() != "dev" {
		t.Fatalf("blank fallback=%q", Version())
	}
}

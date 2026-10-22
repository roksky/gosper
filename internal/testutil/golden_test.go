package testutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoldenUpdateCompareAndDiff(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "sample.golden")
	if err := goldenCheck([]byte("a\nb\n"), path, false); err == nil || !strings.Contains(err.Error(), "-update") {
		t.Fatalf("missing file error=%v", err)
	}
	if err := goldenCheck([]byte("a\nb\n"), path, true); err != nil {
		t.Fatal(err)
	}
	if err := goldenCheck([]byte("a\nb\n"), path, false); err != nil {
		t.Fatal(err)
	}
	err := goldenCheck([]byte("a\nc\n"), path, false)
	if err == nil || !strings.Contains(err.Error(), "line 2, column 1") || !strings.Contains(err.Error(), "- expected") || !strings.Contains(err.Error(), "+ actual") {
		t.Fatalf("diff=%v", err)
	}
	data, readErr := os.ReadFile(path)
	if readErr != nil || string(data) != "a\nb\n" {
		t.Fatalf("mismatch changed fixture: %q,%v", data, readErr)
	}
}
func TestGoldenBinaryDiff(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.bin")
	if err := goldenCheck([]byte{0, 1, 2}, path, true); err != nil {
		t.Fatal(err)
	}
	err := goldenCheck([]byte{0, 1, 3}, path, false)
	if err == nil || !strings.Contains(err.Error(), "byte 2") {
		t.Fatalf("binary diff=%v", err)
	}
}

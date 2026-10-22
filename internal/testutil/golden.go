package testutil

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

var updateGolden = flag.Bool("update", false, "rewrite golden fixtures with current test output")

// Golden compares got with a fixture at path. The fixture is changed only
// when the test binary receives -update. Paths are relative to the calling
// package, matching Go's normal test working directory.
func Golden(t testing.TB, got []byte, path string) {
	t.Helper()
	if err := goldenCheck(got, path, *updateGolden); err != nil {
		t.Fatal(err)
	}
}
func goldenCheck(got []byte, path string, update bool) error {
	if update {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("create golden directory: %w", err)
		}
		if err := os.WriteFile(path, got, 0644); err != nil {
			return fmt.Errorf("write golden %s: %w", path, err)
		}
		return nil
	}
	want, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("golden file missing: %s (run go test -args -update)", path)
	}
	if err != nil {
		return fmt.Errorf("read golden %s: %w", path, err)
	}
	if string(want) == string(got) {
		return nil
	}
	offset := 0
	for offset < len(want) && offset < len(got) && want[offset] == got[offset] {
		offset++
	}
	line, col := 1, 1
	lineStart := 0
	for i := 0; i < offset && i < len(want); i++ {
		if want[i] == '\n' {
			line++
			col = 1
			lineStart = i + 1
		} else {
			col++
		}
	}
	expected := lineAt(want, lineStart)
	actual := lineAt(got, lineStart)
	return fmt.Errorf("golden mismatch %s at byte %d (line %d, column %d):\n- expected %s\n+ actual   %s", path, offset, line, col, strconv.QuoteToASCII(expected), strconv.QuoteToASCII(actual))
}
func lineAt(data []byte, start int) string {
	if start >= len(data) {
		return "<EOF>"
	}
	end := start
	for end < len(data) && data[end] != '\n' && end-start < 120 {
		end++
	}
	return string(data[start:end])
}

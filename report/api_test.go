package report_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/roksky/gosper/report"
	"io"
	"os"
	"strings"
	"testing"
)

type fakeCursor struct{}

func (fakeCursor) Next(context.Context) (report.Record, bool, error) { return nil, false, nil }
func (fakeCursor) Close() error                                      { return nil }

type fakeSource struct{}

func (fakeSource) Open(context.Context) (report.Cursor, error) { return fakeCursor{}, nil }

var (
	_ report.DataSource = fakeSource{}
	_ report.Cursor     = fakeCursor{}
	_ report.Record     = record{}
)

type record struct{}

func (record) Value(string) (any, bool) { return nil, false }
func fixture(t *testing.T) *report.CompiledReport {
	t.Helper()
	file, err := os.Open("../testdata/m1/hello.jrxml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	compiled, err := report.Compile(context.Background(), file)
	if err != nil {
		t.Fatal(err)
	}
	return compiled
}
func TestFormatConstants(t *testing.T) {
	for f, want := range map[report.Format]string{report.FormatPDF: "pdf", report.FormatXLSX: "xlsx", report.FormatHTML: "html", report.FormatCSV: "csv"} {
		if f.String() != want {
			t.Errorf("%q != %q", f, want)
		}
	}
}
func TestPublicWalkingSkeleton(t *testing.T) {
	compiled := fixture(t)
	source := report.FromSlice([]map[string]string{{"name": "Alice"}, {"name": "Bob"}, {"name": "Cara"}})
	filled, err := report.Fill(context.Background(), compiled, nil, source)
	if err != nil {
		t.Fatal(err)
	}
	if filled.PageCount() != 2 {
		t.Fatalf("pages=%d", filled.PageCount())
	}
	var out bytes.Buffer
	if err := report.Export(context.Background(), filled, &out, report.FormatPDF); err != nil {
		t.Fatal(err)
	}
	for _, part := range []string{"%PDF-1.4", "/Count 2", "(Alice!) Tj", "(Cara!) Tj"} {
		if !strings.Contains(out.String(), part) {
			t.Errorf("PDF missing %q", part)
		}
	}
	overridden, err := report.Fill(context.Background(), compiled, report.Params{"suffix": "?"}, source)
	if err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if err := report.Export(context.Background(), overridden, &out, report.FormatPDF); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "(Alice?) Tj") {
		t.Fatal("parameter override missing")
	}
}
func TestFillWithoutExportAndEmptyExport(t *testing.T) {
	compiled := fixture(t)
	filled, err := report.Fill(context.Background(), compiled, nil, fakeSource{})
	if err != nil || filled.PageCount() != 0 {
		t.Fatalf("empty fill=%v,%v", filled, err)
	}
	var out bytes.Buffer
	if err := report.Export(context.Background(), filled, &out, report.FormatPDF); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "/Count 0") {
		t.Fatal("empty PDF count missing")
	}
	if err := report.Export(context.Background(), filled, io.Discard, report.Format("unknown")); !errors.Is(err, report.ErrNotImplemented) {
		t.Fatalf("HTML error=%v", err)
	}
}

type cancelCursor struct {
	cancel context.CancelFunc
	closed bool
}

func (c *cancelCursor) Next(ctx context.Context) (report.Record, bool, error) {
	c.cancel()
	return nil, false, ctx.Err()
}
func (c *cancelCursor) Close() error { c.closed = true; return nil }

type cancelSource struct{ cursor *cancelCursor }

func (s cancelSource) Open(context.Context) (report.Cursor, error) { return s.cursor, nil }
func TestCancellationClosesCursor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cur := &cancelCursor{cancel: cancel}
	_, err := report.Fill(ctx, fixture(t), nil, cancelSource{cur})
	if !errors.Is(err, context.Canceled) || !cur.closed {
		t.Fatalf("cancel=%v closed=%v", err, cur.closed)
	}
}
func TestCompileValidationAndNilInputs(t *testing.T) {
	if _, err := report.Compile(context.Background(), strings.NewReader("<broken>")); err == nil {
		t.Fatal("invalid XML compiled")
	}
	if _, err := report.Fill(context.Background(), nil, nil, fakeSource{}); err == nil {
		t.Fatal("nil compiled report filled")
	}
	if err := report.Export(context.Background(), nil, io.Discard, report.FormatPDF); err == nil {
		t.Fatal("nil filled report exported")
	}
}

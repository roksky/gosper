package report_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/roksky/gosper/report"
)

// fakeCursor and fakeSource prove the exported aliases can be implemented by an
// external caller without importing internal packages.
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

func TestFormatConstants(t *testing.T) {
	cases := map[report.Format]string{
		report.FormatPDF:  "pdf",
		report.FormatXLSX: "xlsx",
		report.FormatHTML: "html",
		report.FormatCSV:  "csv",
	}
	for f, want := range cases {
		if got := f.String(); got != want {
			t.Fatalf("Format %q.String() = %q, want %q", f, got, want)
		}
	}
}

func TestCompileContract(t *testing.T) {
	got, err := report.Compile(context.Background(), strings.NewReader("<jasperReport/>"))
	if got != nil {
		t.Fatalf("Compile returned a non-nil report: %v", got)
	}
	if !errors.Is(err, report.ErrNotImplemented) {
		t.Fatalf("Compile error = %v, want ErrNotImplemented", err)
	}
}

func TestFillContract(t *testing.T) {
	got, err := report.Fill(context.Background(), &report.CompiledReport{}, report.Params{"p": 1}, fakeSource{})
	if got != nil {
		t.Fatalf("Fill returned a non-nil report: %v", got)
	}
	if !errors.Is(err, report.ErrNotImplemented) {
		t.Fatalf("Fill error = %v, want ErrNotImplemented", err)
	}
}

func TestExportContract(t *testing.T) {
	err := report.Export(context.Background(), &report.FilledReport{}, io.Discard, report.FormatPDF)
	if !errors.Is(err, report.ErrNotImplemented) {
		t.Fatalf("Export error = %v, want ErrNotImplemented", err)
	}
}

func TestNilParamsValid(t *testing.T) {
	var p report.Params
	if _, err := report.Fill(context.Background(), &report.CompiledReport{}, p, fakeSource{}); !errors.Is(err, report.ErrNotImplemented) {
		t.Fatalf("Fill with nil Params error = %v, want ErrNotImplemented", err)
	}
}

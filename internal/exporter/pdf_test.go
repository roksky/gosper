package exporter

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/roksky/gosper/internal/datasource"
	"github.com/roksky/gosper/internal/fill"
	"github.com/roksky/gosper/internal/jrxml"
	"github.com/roksky/gosper/internal/layout"
	"github.com/roksky/gosper/internal/limits"
	"github.com/roksky/gosper/internal/printdoc"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func testPDFDocument(t *testing.T) *printdoc.Document {
	t.Helper()
	file, err := os.Open("../../testdata/jrxml/m1-basic.jrxml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	tmpl, diags, err := jrxml.Parse(file)
	if err != nil || len(diags) > 0 {
		t.Fatalf("parse: %v %v", err, diags)
	}
	spec, err := layout.FromJRXML(tmpl)
	if err != nil {
		t.Fatal(err)
	}
	source := datasource.FromSlice([]map[string]string{{"name": "Alice"}, {"name": "Bob"}, {"name": "Cara"}})
	filled, err := fill.Run(context.Background(), &fill.Template{Fields: []string{"name"}}, source, limits.Budget{})
	if err != nil {
		t.Fatal(err)
	}
	doc, err := layout.Layout(context.Background(), filled, spec)
	if err != nil {
		t.Fatal(err)
	}
	return doc
}
func TestPDFStructureAndText(t *testing.T) {
	doc := testPDFDocument(t)
	var out bytes.Buffer
	if err := PDF(context.Background(), doc, &out); err != nil {
		t.Fatal(err)
	}
	data := out.Bytes()
	if !bytes.HasPrefix(data, []byte("%PDF-1.4\n")) {
		t.Fatal("missing PDF header")
	}
	for _, part := range []string{"/Count 2", "/MediaBox [0 0 200.00 100.00]", "(Alice) Tj", "(Bob) Tj", "(Cara) Tj", "%%EOF"} {
		if !bytes.Contains(data, []byte(part)) {
			t.Errorf("missing %q", part)
		}
	}
	checkXref(t, data)
	if _, err := exec.LookPath("pdfinfo"); err == nil {
		path := filepath.Join(t.TempDir(), "out.pdf")
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
		info, err := exec.Command("pdfinfo", path).Output()
		if err != nil || !bytes.Contains(info, []byte("Pages:           2")) {
			t.Fatalf("pdfinfo: %v %s", err, info)
		}
		if _, err := exec.LookPath("pdftotext"); err == nil {
			text, err := exec.Command("pdftotext", path, "-").Output()
			if err != nil || !bytes.Contains(text, []byte("Alice")) || !bytes.Contains(text, []byte("Cara")) {
				t.Fatalf("pdftotext: %v %s", err, text)
			}
		}
	}
}
func checkXref(t *testing.T, data []byte) {
	t.Helper()
	match := regexp.MustCompile(`(?s)startxref\n([0-9]+)\n%%EOF`).FindSubmatch(data)
	if match == nil {
		t.Fatal("missing startxref")
	}
	offset, _ := strconv.Atoi(string(match[1]))
	if offset >= len(data) || !bytes.HasPrefix(data[offset:], []byte("xref\n")) {
		t.Fatalf("bad xref offset %d", offset)
	}
	lines := strings.Split(string(data[offset:]), "\n")
	if len(lines) < 4 {
		t.Fatal("short xref")
	}
	var start, count int
	if _, err := fmt.Sscanf(lines[1], "%d %d", &start, &count); err != nil || start != 0 {
		t.Fatalf("xref header=%q", lines[1])
	}
	for id := 1; id < count; id++ {
		entry := lines[id+2]
		if len(entry) < 10 {
			t.Fatalf("xref entry %d=%q", id, entry)
		}
		pos, err := strconv.Atoi(entry[:10])
		if err != nil || pos >= len(data) || !bytes.HasPrefix(data[pos:], []byte(strconv.Itoa(id)+" 0 obj")) {
			t.Fatalf("xref entry %d points to %d", id, pos)
		}
	}
}

type brokenWriter struct {
	limit int
	wrote int
	err   error
}

func (w *brokenWriter) Write(p []byte) (int, error) {
	if w.wrote >= w.limit {
		return 0, w.err
	}
	n := len(p)
	if n > w.limit-w.wrote {
		n = w.limit - w.wrote
	}
	w.wrote += n
	if n < len(p) {
		return n, w.err
	}
	return n, nil
}
func TestPDFEmptyWriterErrorAndUnsupported(t *testing.T) {
	var out bytes.Buffer
	empty := &printdoc.Document{}
	if err := PDF(context.Background(), empty, &out); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(out.Bytes(), []byte("/Count 0")) {
		t.Fatal("empty PDF missing page count")
	}
	boom := errors.New("write failed")
	w := &brokenWriter{limit: 30, err: boom}
	if err := PDF(context.Background(), testPDFDocument(t), w); !errors.Is(err, boom) {
		t.Fatalf("writer error=%v", err)
	}
	bad := &printdoc.Document{}
	_ = bad.AppendPage(printdoc.Page{Width: 100, Height: 100, Elements: []printdoc.Element{printdoc.NewShapeElement(printdoc.Bounds{}, printdoc.ShapeLine)}})
	if err := PDF(context.Background(), bad, &out); err == nil {
		t.Fatal("unsupported shape accepted")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := PDF(ctx, empty, &out); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancel=%v", err)
	}
}
func TestPDFEscapesTextAndUsesResolvedStyle(t *testing.T) {
	doc := &printdoc.Document{}
	e := printdoc.NewTextElement(printdoc.Bounds{X: 20, Y: 20, Width: 300, Height: 50}, `(A) \\ B`)
	e.Style = printdoc.Style{Font: "Serif", Size: 240, Bold: true, Color: 0xFF0000}
	if err := doc.AppendPage(printdoc.Page{Width: 1000, Height: 1000, Elements: []printdoc.Element{e}}); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := PDF(context.Background(), doc, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "/BaseFont /Times-Bold") || !strings.Contains(out.String(), "1.0000 0.0000 0.0000 rg") || !strings.Contains(out.String(), `\(A\) \\`) {
		t.Fatalf("style or escape missing: %s", out.String())
	}
}

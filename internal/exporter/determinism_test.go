package exporter

import (
	"bytes"
	"context"
	"crypto/sha256"
	"github.com/roksky/gosper/internal/printdoc"
	"sync"
	"testing"
)

func deterministicDocument(t *testing.T) *printdoc.Document {
	t.Helper()
	doc := testPDFDocument(t)
	extra := printdoc.NewTextElement(printdoc.Bounds{X: 200, Y: 200, Width: 1000, Height: 300}, "Stable")
	extra.Style = printdoc.Style{Font: "Serif", Size: 240, Bold: true, Color: 0x336699}
	doc.Pages[0].Elements = append(doc.Pages[0].Elements, extra)
	other := printdoc.NewTextElement(printdoc.Bounds{X: 200, Y: 200, Width: 1000, Height: 300}, "Order")
	other.Style = printdoc.Style{Font: "Monospaced", Size: 200, Italic: true}
	doc.Pages[1].Elements = append(doc.Pages[1].Elements, other)
	return doc
}
func TestPDFByteIdenticalRepeated(t *testing.T) {
	doc := deterministicDocument(t)
	var first bytes.Buffer
	if err := PDF(context.Background(), doc, &first); err != nil {
		t.Fatal(err)
	}
	reference := sha256.Sum256(first.Bytes())
	for i := 0; i < 25; i++ {
		var next bytes.Buffer
		if err := PDF(context.Background(), doc, &next); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(first.Bytes(), next.Bytes()) || sha256.Sum256(next.Bytes()) != reference {
			t.Fatalf("export %d differs", i)
		}
	}
	for _, key := range [][]byte{[]byte("/CreationDate"), []byte("/ModDate"), []byte("/ID [")} {
		if bytes.Contains(first.Bytes(), key) {
			t.Fatalf("variable metadata %q found", key)
		}
	}
}
func TestPDFByteIdenticalParallel(t *testing.T) {
	doc := deterministicDocument(t)
	var reference bytes.Buffer
	if err := PDF(context.Background(), doc, &reference); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var result bytes.Buffer
			if err := PDF(context.Background(), doc, &result); err != nil {
				t.Errorf("PDF: %v", err)
				return
			}
			if !bytes.Equal(reference.Bytes(), result.Bytes()) {
				t.Error("parallel export differs")
			}
		}()
	}
	wg.Wait()
}

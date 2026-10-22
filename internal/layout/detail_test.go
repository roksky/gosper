package layout

import (
	"context"
	"github.com/roksky/gosper/internal/datasource"
	"github.com/roksky/gosper/internal/fill"
	"github.com/roksky/gosper/internal/jrxml"
	"github.com/roksky/gosper/internal/limits"
	"github.com/roksky/gosper/internal/printdoc"
	"os"
	"testing"
)

func basicSpec(t *testing.T) *DetailSpec {
	t.Helper()
	file, err := os.Open("../../testdata/jrxml/m1-basic.jrxml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	tmpl, diags, err := jrxml.Parse(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(diags) > 0 {
		t.Fatalf("parse diagnostics: %v", diags)
	}
	spec, err := FromJRXML(tmpl)
	if err != nil {
		t.Fatal(err)
	}
	return spec
}
func TestDetailLayoutFromJRXML(t *testing.T) {
	spec := basicSpec(t)
	source := datasource.FromSlice([]map[string]string{{"name": "Alice"}, {"name": "Bob"}, {"name": "Cara"}})
	filled, err := fill.Run(context.Background(), &fill.Template{Fields: []string{"name"}}, source, limits.Budget{})
	if err != nil {
		t.Fatal(err)
	}
	doc, err := Layout(context.Background(), filled, spec)
	if err != nil {
		t.Fatal(err)
	}
	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}
	if doc.Metadata.PageCount != 2 {
		t.Fatalf("pages=%d", doc.Metadata.PageCount)
	}
	if doc.Pages[0].Width != printdoc.FromPoints(200) || doc.Pages[1].Height != printdoc.FromPoints(100) {
		t.Fatalf("page sizes=%+v", doc.Pages)
	}
	for _, tc := range []struct {
		page int
		text string
		x, y printdoc.Unit
	}{{0, "Header", 200, 200}, {0, "Footer", 200, 1600}, {0, "Items", 200, 400}, {0, "Alice", 300, 600}, {0, "Bob", 300, 1000}, {1, "Header", 200, 200}, {1, "Footer", 200, 1600}, {1, "Cara", 300, 400}} {
		found := false
		for _, el := range doc.Pages[tc.page].Elements {
			if el.Text == tc.text {
				found = true
				if el.Bounds.X != tc.x || el.Bounds.Y != tc.y || el.Bounds.Width <= 0 {
					t.Errorf("%s bounds=%+v", tc.text, el.Bounds)
				}
			}
		}
		if !found {
			t.Errorf("missing %q on page %d", tc.text, tc.page+1)
		}
	}
}
func TestDetailLayoutEmptyAndOversize(t *testing.T) {
	spec := basicSpec(t)
	doc, err := Layout(context.Background(), &fill.Filled{}, spec)
	if err != nil || doc.Metadata.PageCount != 1 {
		t.Fatalf("title-only=%+v,%v", doc, err)
	}
	spec.Title = nil
	doc, err = Layout(context.Background(), &fill.Filled{}, spec)
	if err != nil || doc.Metadata.PageCount != 0 {
		t.Fatalf("empty=%+v,%v", doc, err)
	}
	spec.Detail.Height = printdoc.FromPoints(70)
	_, err = Layout(context.Background(), &fill.Filled{Records: []fill.FilledRecord{{}}}, spec)
	if err == nil {
		t.Fatal("oversize detail accepted")
	}
}
func TestDetailExactlyFullPageAndCancel(t *testing.T) {
	spec := &DetailSpec{Page: PageGeometry{Width: 200, Height: 100, Top: 10, Bottom: 10, Left: 10, Right: 10}, Detail: Band{Height: 40}}
	filled := &fill.Filled{Records: []fill.FilledRecord{{}, {}, {}}}
	doc, err := Layout(context.Background(), filled, spec)
	if err != nil || doc.Metadata.PageCount != 2 {
		t.Fatalf("full page result=%+v,%v", doc, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := Layout(ctx, filled, spec); err != context.Canceled {
		t.Fatalf("cancel error=%v", err)
	}
}

func TestLayoutCarriesResolvedTextStyle(t *testing.T) {
	spec := &DetailSpec{
		Page:   PageGeometry{Width: 2000, Height: 1000, Top: 100, Bottom: 100, Left: 100, Right: 100},
		Detail: Band{Height: 400, Texts: []TextItem{{Bounds: printdoc.Bounds{Width: 1000, Height: 300}, Literal: "Hi", Font: "Times", Size: 12}}},
	}
	doc, err := Layout(context.Background(), &fill.Filled{Records: []fill.FilledRecord{{}}}, spec)
	if err != nil {
		t.Fatal(err)
	}
	style := doc.Pages[0].Elements[0].Style
	if style.Font != "Serif" || style.Size != 240 || style.Color != 0 {
		t.Fatalf("style=%+v", style)
	}
}

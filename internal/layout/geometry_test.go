package layout

import (
	"github.com/roksky/gosper/internal/printdoc"
	"testing"
)

func testGeometry() PageGeometry {
	return PageGeometry{Width: 200, Height: 300, Top: 20, Bottom: 30, Left: 10, Right: 15}
}
func TestPrintableAreaAndBandPlacement(t *testing.T) {
	c, err := NewCanvas(testGeometry())
	if err != nil {
		t.Fatal(err)
	}
	if c.Doc.Metadata.PageCount != 0 {
		t.Fatal("canvas started with a page")
	}
	if got := c.PrintableArea(); got != (printdoc.Bounds{X: 10, Y: 20, Width: 175, Height: 250}) {
		t.Fatalf("area=%+v", got)
	}
	for _, tc := range []struct {
		h    printdoc.Unit
		page int
		y    printdoc.Unit
	}{{100, 1, 20}, {150, 1, 120}, {1, 2, 20}} {
		p, err := c.Place(tc.h)
		if err != nil || p.Page != tc.page || p.X != 10 || p.Y != tc.y {
			t.Fatalf("Place(%d)=%+v,%v", tc.h, p, err)
		}
	}
	if c.Doc.Metadata.PageCount != 2 || c.Doc.Pages[1].Width != 200 {
		t.Fatalf("document=%+v", c.Doc)
	}
}
func TestGeometryInvalidAndOversizeBand(t *testing.T) {
	for _, g := range []PageGeometry{{Width: 0, Height: 10}, {Width: 10, Height: 10, Left: 5, Right: 5}, {Width: 10, Height: 10, Top: -1}} {
		if _, err := g.PrintableArea(); err == nil {
			t.Errorf("accepted %+v", g)
		}
	}
	c, _ := NewCanvas(testGeometry())
	if _, err := c.Place(251); err == nil {
		t.Fatal("oversize band accepted")
	}
	if _, err := c.Place(-1); err == nil {
		t.Fatal("negative band accepted")
	}
	if c.Doc.Metadata.PageCount != 0 {
		t.Fatal("failed placement changed document")
	}
}
func TestExactlyFullPage(t *testing.T) {
	c, _ := NewCanvas(testGeometry())
	first, _ := c.Place(250)
	second, _ := c.Place(0)
	if first.Page != 1 || second.Page != 1 {
		t.Fatalf("placements=%+v,%+v", first, second)
	}
	third, _ := c.Place(1)
	if third.Page != 2 {
		t.Fatalf("new page=%+v", third)
	}
}

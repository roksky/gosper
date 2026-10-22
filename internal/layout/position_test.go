package layout

import (
	"github.com/roksky/gosper/internal/printdoc"
	"testing"
)

func TestAbsoluteTextPlacement(t *testing.T) {
	c, _ := NewCanvas(testGeometry())
	p, _ := c.Place(80)
	relative := printdoc.Bounds{X: 5, Y: 7, Width: 40, Height: 20}
	if err := c.PlaceText(p, relative, "bound"); err != nil {
		t.Fatal(err)
	}
	got := c.Doc.Pages[0].Elements[0]
	if got.Bounds != (printdoc.Bounds{X: 15, Y: 27, Width: 40, Height: 20}) || got.Text != "bound" {
		t.Fatalf("element=%+v", got)
	}
	next, _ := c.Place(200)
	if next.Page != 2 {
		t.Fatalf("page=%d", next.Page)
	}
	if err := c.PlaceText(next, relative, "second"); err != nil {
		t.Fatal(err)
	}
	if c.Doc.Pages[1].Elements[0].Bounds.Y != 27 {
		t.Fatal("page two origin wrong")
	}
}
func TestAbsolutePositionRejectsOverflow(t *testing.T) {
	c, _ := NewCanvas(testGeometry())
	p, _ := c.Place(30)
	for _, b := range []printdoc.Bounds{{X: -1, Width: 1, Height: 1}, {X: 170, Width: 6, Height: 1}, {Y: 29, Width: 1, Height: 2}, {Width: -1}, {Y: -1}} {
		if err := c.PlaceText(p, b, "bad"); err == nil {
			t.Errorf("accepted %+v", b)
		}
	}
	if len(c.Doc.Pages[0].Elements) != 0 {
		t.Fatal("rejected element appended")
	}
}

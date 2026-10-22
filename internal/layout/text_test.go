package layout

import (
	"context"
	"testing"

	"github.com/roksky/gosper/internal/fill"
	"github.com/roksky/gosper/internal/printdoc"
)

func TestWrappedTextAlignment(t *testing.T) {
	box := printdoc.Bounds{X: 10, Y: 20, Width: printdoc.FromPoints(30), Height: printdoc.FromPoints(50)}
	spec := &DetailSpec{
		Page:   PageGeometry{Width: 2000, Height: 2000, Top: 100, Bottom: 100, Left: 100, Right: 100},
		Detail: Band{Height: 1200, Texts: []TextItem{{Bounds: box, Literal: "one two", Font: "Monospaced", Size: 10, Horizontal: printdoc.AlignRight, Vertical: AlignBottom}}},
	}
	doc, err := Layout(context.Background(), &fill.Filled{Records: []fill.FilledRecord{{}}}, spec)
	if err != nil {
		t.Fatal(err)
	}
	els := doc.Pages[0].Elements
	if len(els) != 2 || els[0].Text != "one" || els[1].Text != "two" {
		t.Fatalf("lines = %+v", els)
	}
	// Three 6pt monospaced glyphs occupy 18pt; two 12pt lines fit in 50pt.
	wantX := printdoc.Unit(100) + box.X + printdoc.FromPoints(12)
	wantY := printdoc.Unit(100) + box.Y + printdoc.FromPoints(26)
	if els[0].Bounds.X != wantX || els[0].Bounds.Y != wantY || els[1].Bounds.Y != wantY+printdoc.FromPoints(12) {
		t.Fatalf("bounds = %+v, %+v", els[0].Bounds, els[1].Bounds)
	}
	spec.Detail.Texts[0].Horizontal = printdoc.AlignCenter
	spec.Detail.Texts[0].Vertical = AlignMiddle
	doc, err = Layout(context.Background(), &fill.Filled{Records: []fill.FilledRecord{{}}}, spec)
	if err != nil {
		t.Fatal(err)
	}
	if got := doc.Pages[0].Elements[0].Bounds; got.X != 100+box.X+printdoc.FromPoints(6) || got.Y != 100+box.Y+printdoc.FromPoints(13) {
		t.Fatalf("center/middle = %+v", got)
	}
}

func TestWrappedTextTooTall(t *testing.T) {
	_, err := laidOutText(context.Background(), TextItem{Bounds: printdoc.Bounds{Width: printdoc.FromPoints(18), Height: printdoc.FromPoints(12)}, Literal: "one two", Font: "Monospaced", Size: 10}, nil)
	if err == nil {
		t.Fatal("overflow accepted")
	}
}

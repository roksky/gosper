package layout

import (
	"fmt"
	"github.com/roksky/gosper/internal/printdoc"
)

// PageGeometry defines page dimensions and margins in print-document twips.
type PageGeometry struct {
	Width, Height            printdoc.Unit
	Top, Bottom, Left, Right printdoc.Unit
}

// PrintableArea returns the rectangle available for band placement.
func (g PageGeometry) PrintableArea() (printdoc.Bounds, error) {
	if g.Width <= 0 || g.Height <= 0 || g.Top < 0 || g.Bottom < 0 || g.Left < 0 || g.Right < 0 || g.Left >= g.Width || g.Right >= g.Width-g.Left || g.Top >= g.Height || g.Bottom >= g.Height-g.Top {
		return printdoc.Bounds{}, fmt.Errorf("layout: invalid page geometry")
	}
	return printdoc.Bounds{X: g.Left, Y: g.Top, Width: g.Width - g.Left - g.Right, Height: g.Height - g.Top - g.Bottom}, nil
}

// Placement is the absolute origin of one band on a one-based page number.
type Placement struct {
	Page int
	X, Y printdoc.Unit
}

// Canvas owns pagination state and the resulting print document for one run.
type Canvas struct {
	Doc      *printdoc.Document
	Geometry PageGeometry
	area     printdoc.Bounds
	cursor   printdoc.Unit
}

// NewCanvas validates geometry and starts with an empty document. The first
// successful Place creates page one.
func NewCanvas(g PageGeometry) (*Canvas, error) {
	area, err := g.PrintableArea()
	if err != nil {
		return nil, err
	}
	return &Canvas{Doc: &printdoc.Document{}, Geometry: g, area: area, cursor: area.Y}, nil
}

// Place reserves a fixed-height band and returns its absolute page origin.
// A band taller than the printable area fails without changing the canvas.
func (c *Canvas) Place(height printdoc.Unit) (Placement, error) {
	if c == nil || c.Doc == nil {
		return Placement{}, fmt.Errorf("layout: nil canvas")
	}
	if height < 0 || height > c.area.Height {
		return Placement{}, fmt.Errorf("layout: band height %d exceeds printable height %d", height, c.area.Height)
	}
	if len(c.Doc.Pages) == 0 || height > c.area.Y+c.area.Height-c.cursor {
		if err := c.Doc.AppendPage(printdoc.Page{Width: c.Geometry.Width, Height: c.Geometry.Height}); err != nil {
			return Placement{}, err
		}
		c.cursor = c.area.Y
	}
	result := Placement{Page: len(c.Doc.Pages), X: c.area.X, Y: c.cursor}
	c.cursor += height
	return result, nil
}

// PrintableArea returns the validated printable rectangle used by this canvas.
func (c *Canvas) PrintableArea() printdoc.Bounds { return c.area }

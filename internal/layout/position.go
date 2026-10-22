package layout

import (
	"fmt"
	"github.com/roksky/gosper/internal/printdoc"
)

// AbsoluteBounds translates band-relative bounds to a page rectangle. An
// element outside the printable width or its band height is rejected.
func (c *Canvas) AbsoluteBounds(p Placement, relative printdoc.Bounds) (printdoc.Bounds, error) {
	if c == nil || c.Doc == nil || p.Page < 1 || p.Page > len(c.Doc.Pages) {
		return printdoc.Bounds{}, fmt.Errorf("layout: invalid placement page")
	}
	if p.X != c.area.X || p.Y < c.area.Y || p.Y > c.area.Y+c.area.Height || p.Height < 0 || p.Height > c.area.Y+c.area.Height-p.Y {
		return printdoc.Bounds{}, fmt.Errorf("layout: invalid band origin")
	}
	if relative.X < 0 || relative.Y < 0 || relative.Width < 0 || relative.Height < 0 || relative.X > c.area.Width || relative.Width > c.area.Width-relative.X || relative.Y > p.Height || relative.Height > p.Height-relative.Y {
		return printdoc.Bounds{}, fmt.Errorf("layout: element outside printable band")
	}
	return printdoc.Bounds{X: p.X + relative.X, Y: p.Y + relative.Y, Width: relative.Width, Height: relative.Height}, nil
}

// PlaceText appends already-bound text at an absolute page position.
func (c *Canvas) PlaceText(p Placement, relative printdoc.Bounds, value string) error {
	return c.PlaceStyledText(p, relative, value, printdoc.DefaultStyle())
}

// PlaceStyledText appends already-bound text with a resolved style.
func (c *Canvas) PlaceStyledText(p Placement, relative printdoc.Bounds, value string, style printdoc.Style) error {
	if err := style.Validate(); err != nil {
		return err
	}
	absolute, err := c.AbsoluteBounds(p, relative)
	if err != nil {
		return err
	}
	element := printdoc.NewTextElement(absolute, value)
	element.Style = style
	c.Doc.Pages[p.Page-1].Elements = append(c.Doc.Pages[p.Page-1].Elements, element)
	return nil
}

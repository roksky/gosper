package layout

import (
	"context"
	"fmt"
	"github.com/roksky/gosper/internal/expression"
	"github.com/roksky/gosper/internal/fill"
	"github.com/roksky/gosper/internal/jrxml"
	"github.com/roksky/gosper/internal/limits"
	"github.com/roksky/gosper/internal/printdoc"
	"github.com/roksky/gosper/internal/text"
	"strconv"
	"time"
)

// TextItem is one fixed text box within a band. Expr, when present, is
// evaluated against the current filled record; otherwise Literal is used.
type TextItem struct {
	Bounds     printdoc.Bounds
	Literal    string
	Expr       *expression.Program
	Font       string
	Size       float64
	Horizontal printdoc.Align
	Vertical   VerticalAlign
}

// VerticalAlign controls placement of a text block within its box.
type VerticalAlign uint8

const (
	AlignTop VerticalAlign = iota
	AlignMiddle
	AlignBottom
)

// Band is a fixed-height collection of text items.
type Band struct {
	Height printdoc.Unit
	Texts  []TextItem
}

// DetailSpec describes the M1 bands and page geometry.
type DetailSpec struct {
	Page                          PageGeometry
	Title, PageHeader, PageFooter *Band
	Detail                        Band
}

// FromJRXML converts the supported JRXML page and text bands into a layout
// specification. Columns and other band kinds remain outside the M1 subset.
func FromJRXML(t *jrxml.Template) (*DetailSpec, error) {
	if t == nil {
		return nil, fmt.Errorf("layout: nil JRXML template")
	}
	if t.Page.ColumnCount != 1 {
		return nil, fmt.Errorf("layout: M1 supports one column")
	}
	spec := &DetailSpec{Page: PageGeometry{Width: t.Page.Width, Height: t.Page.Height, Top: t.Page.MarginTop, Bottom: t.Page.MarginBottom, Left: t.Page.MarginLeft, Right: t.Page.MarginRight}}
	seen := map[jrxml.BandKind]bool{}
	for _, source := range t.Bands {
		if seen[source.Kind] {
			return nil, fmt.Errorf("layout: duplicate %s band", source.Kind)
		}
		seen[source.Kind] = true
		band := Band{Height: source.Height}
		for _, el := range source.Elements {
			item := TextItem{Bounds: printdoc.Bounds{X: el.X, Y: el.Y, Width: el.Width, Height: el.Height}, Literal: el.Text, Font: "SansSerif", Size: 10}
			switch el.Kind {
			case jrxml.KindStaticText:
			case jrxml.KindTextField:
				program, err := expression.Compile(el.Expr)
				if err != nil {
					return nil, fmt.Errorf("layout: %s expression: %w", source.Kind, err)
				}
				item.Expr = program
			default:
				return nil, fmt.Errorf("layout: unsupported element %s", el.Kind)
			}
			band.Texts = append(band.Texts, item)
		}
		switch source.Kind {
		case jrxml.BandTitle:
			spec.Title = &band
		case jrxml.BandPageHeader:
			spec.PageHeader = &band
		case jrxml.BandPageFooter:
			spec.PageFooter = &band
		case jrxml.BandDetail:
			spec.Detail = band
		default:
			return nil, fmt.Errorf("layout: unsupported band %s", source.Kind)
		}
	}
	return spec, nil
}

// Layout places a fixed-height detail band for each filled record. A title is
// placed once, while page headers and footers occur on every produced page.
// Empty data without a title yields an empty document.
func Layout(ctx context.Context, filled *fill.Filled, spec *DetailSpec) (*printdoc.Document, error) {
	return LayoutWithBudget(ctx, filled, spec, limits.Budget{})
}

// LayoutWithBudget enforces the page budget as pages are created.
func LayoutWithBudget(ctx context.Context, filled *fill.Filled, spec *DetailSpec, budget limits.Budget) (*printdoc.Document, error) {
	if ctx == nil || filled == nil || spec == nil {
		return nil, fmt.Errorf("layout: nil input")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	headerH, footerH := bandHeight(spec.PageHeader), bandHeight(spec.PageFooter)
	if headerH < 0 || footerH < 0 || spec.Page.Top < 0 || spec.Page.Bottom < 0 || spec.Page.Height <= 0 || headerH > spec.Page.Height-spec.Page.Top || footerH > spec.Page.Height-spec.Page.Bottom-headerH-spec.Page.Top {
		return nil, fmt.Errorf("layout: headers and footers exceed page")
	}
	body := spec.Page
	body.Top += headerH
	body.Bottom += footerH
	canvas, err := NewCanvasWithGuard(body, budget.Guard(ctx))
	if err != nil {
		return nil, err
	}
	decorate := func(page int, record *fill.FilledRecord) error {
		if spec.PageHeader != nil {
			if err := appendBand(ctx, canvas, page, spec.Page.Left, spec.Page.Top, *spec.PageHeader, record); err != nil {
				return err
			}
		}
		if spec.PageFooter != nil {
			y := spec.Page.Height - spec.Page.Bottom - footerH
			if err := appendBand(ctx, canvas, page, spec.Page.Left, y, *spec.PageFooter, record); err != nil {
				return err
			}
		}
		return nil
	}
	if spec.Title != nil {
		place, err := canvas.Place(spec.Title.Height)
		if err != nil {
			return nil, err
		}
		var record *fill.FilledRecord
		if len(filled.Records) > 0 {
			record = &filled.Records[0]
		}
		if err := decorate(place.Page, record); err != nil {
			return nil, err
		}
		if err := appendBodyBand(ctx, canvas, place, *spec.Title, record); err != nil {
			return nil, err
		}
	}
	if len(filled.Records) > 0 && spec.Detail.Height <= 0 {
		return nil, fmt.Errorf("layout: detail height must be positive")
	}
	for i := range filled.Records {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		place, err := canvas.Place(spec.Detail.Height)
		if err != nil {
			return nil, err
		}
		if place.Page > canvas.Doc.Metadata.PageCount {
			return nil, fmt.Errorf("layout: inconsistent page count")
		}
		if place.Page > 1 || spec.Title == nil {
			if isNewPagePlacement(canvas, place) {
				if err := decorate(place.Page, &filled.Records[i]); err != nil {
					return nil, err
				}
			}
		}
		if err := appendBodyBand(ctx, canvas, place, spec.Detail, &filled.Records[i]); err != nil {
			return nil, err
		}
	}
	return canvas.Doc, nil
}
func bandHeight(b *Band) printdoc.Unit {
	if b == nil {
		return 0
	}
	return b.Height
}
func isNewPagePlacement(c *Canvas, p Placement) bool { return p.Y == c.area.Y }
func appendBodyBand(ctx context.Context, c *Canvas, p Placement, b Band, record *fill.FilledRecord) error {
	for _, item := range b.Texts {
		lines, err := laidOutText(ctx, item, record)
		if err != nil {
			return err
		}
		for _, line := range lines {
			if err := c.PlaceStyledText(p, line.bounds, line.value, line.style); err != nil {
				return err
			}
		}
	}
	return nil
}
func appendBand(ctx context.Context, c *Canvas, page int, x, y printdoc.Unit, b Band, record *fill.FilledRecord) error {
	area, err := c.Geometry.PrintableArea()
	if err != nil {
		return err
	}
	for _, item := range b.Texts {
		lines, err := laidOutText(ctx, item, record)
		if err != nil {
			return err
		}
		for _, line := range lines {
			bounds := line.bounds
			if bounds.X < 0 || bounds.Y < 0 || bounds.Width < 0 || bounds.Height < 0 || bounds.X > area.Width || bounds.Width > area.Width-bounds.X || bounds.Y > b.Height || bounds.Height > b.Height-bounds.Y {
				return fmt.Errorf("layout: element outside header or footer band")
			}
			absolute := printdoc.Bounds{X: x + bounds.X, Y: y + bounds.Y, Width: bounds.Width, Height: bounds.Height}
			element := printdoc.NewTextElement(absolute, line.value)
			element.Style = line.style
			c.Doc.Pages[page-1].Elements = append(c.Doc.Pages[page-1].Elements, element)
		}
	}
	return nil
}

type positionedLine struct {
	value  string
	bounds printdoc.Bounds
	style  printdoc.Style
}

func laidOutText(ctx context.Context, item TextItem, record *fill.FilledRecord) ([]positionedLine, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	value := item.Literal
	if item.Expr != nil {
		var scope expression.Scope
		if record != nil {
			scope = expression.Scope{*record}
		}
		result, err := item.Expr.Eval(ctx, scope)
		if err != nil {
			return nil, err
		}
		value = valueString(result)
	}
	size := item.Size
	if size == 0 {
		size = 10
	}
	font, err := text.Resolve(item.Font, 0, size)
	if err != nil {
		return nil, err
	}
	if item.Horizontal > printdoc.AlignRight || item.Vertical > AlignBottom {
		return nil, fmt.Errorf("layout: invalid text alignment")
	}
	lines, err := text.Wrap(value, font, item.Bounds.Width)
	if err != nil {
		return nil, err
	}
	metrics, err := text.Measure("", font)
	if err != nil {
		return nil, err
	}
	blockHeight := printdoc.Unit(len(lines)) * metrics.LineHeight
	if len(lines) == 1 {
		blockHeight = metrics.Ascent + metrics.Descent
	}
	if blockHeight > item.Bounds.Height {
		return nil, fmt.Errorf("layout: wrapped text height %d exceeds box height %d", blockHeight, item.Bounds.Height)
	}
	y := item.Bounds.Y
	switch item.Vertical {
	case AlignMiddle:
		y += (item.Bounds.Height - blockHeight) / 2
	case AlignBottom:
		y += item.Bounds.Height - blockHeight
	}
	style := printdoc.DefaultStyle()
	style.Font = font.Name
	style.Size = printdoc.FromPoints(size)
	style.Align = item.Horizontal
	result := make([]positionedLine, 0, len(lines))
	for _, line := range lines {
		m, err := text.Measure(line, font)
		if err != nil {
			return nil, err
		}
		x := item.Bounds.X
		switch item.Horizontal {
		case printdoc.AlignCenter:
			x += (item.Bounds.Width - m.Width) / 2
		case printdoc.AlignRight:
			x += item.Bounds.Width - m.Width
		}
		height := metrics.LineHeight
		if len(lines) == 1 {
			height = blockHeight
		}
		result = append(result, positionedLine{line, printdoc.Bounds{X: x, Y: y, Width: m.Width, Height: height}, style})
		y += metrics.LineHeight
	}
	return result, nil
}
func valueString(v expression.Value) string {
	switch x := v.(type) {
	case nil, expression.Null:
		return ""
	case expression.String:
		return string(x)
	case expression.Int:
		return strconv.FormatInt(int64(x), 10)
	case expression.Decimal:
		return string(x)
	case expression.Bool:
		return strconv.FormatBool(bool(x))
	case expression.Time:
		return time.Time(x).Format(time.RFC3339)
	default:
		return fmt.Sprint(x)
	}
}

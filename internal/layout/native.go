package layout

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"reflect"
	"strings"

	"github.com/roksky/gosper/internal/expression"
	"github.com/roksky/gosper/internal/fill"
	"github.com/roksky/gosper/internal/jrxml"
	"github.com/roksky/gosper/internal/limits"
	"github.com/roksky/gosper/internal/printdoc"
	"github.com/roksky/gosper/internal/text"
)

// NativeSpec extends the fixed text layout with standard Jasper sections,
// nested frames, text styling, shapes, images and deferred page evaluation.
type NativeSpec struct {
	Template *jrxml.Template
	Programs map[string]*expression.Program
}

func useNative(t *jrxml.Template) bool {
	if len(t.Groups) > 0 || len(t.Variables) > 0 {
		return true
	}
	for _, b := range t.Bands {
		if b.Kind != jrxml.BandTitle && b.Kind != jrxml.BandPageHeader && b.Kind != jrxml.BandPageFooter && b.Kind != jrxml.BandDetail {
			return true
		}
		for _, el := range b.Elements {
			if el.Kind != jrxml.KindStaticText && el.Kind != jrxml.KindTextField || el.FontSize != 0 || el.TextAlignment != "" || el.PrintWhenExpr != "" || el.EvalTime != "Now" {
				return true
			}
		}
	}
	return false
}

func nativeSpec(t *jrxml.Template) (*DetailSpec, error) {
	if t.Page.ColumnCount != 1 {
		return nil, fmt.Errorf("layout: multiple columns are not implemented")
	}
	n := &NativeSpec{Template: t, Programs: map[string]*expression.Program{}}
	var visit func(jrxml.Element) error
	visit = func(el jrxml.Element) error {
		if el.EvalTime != "" && el.EvalTime != "Now" && el.EvalTime != "Report" && el.EvalTime != "Page" {
			return fmt.Errorf("layout: evaluationTime %q is not implemented", el.EvalTime)
		}
		if style, ok := t.Styles[el.StyleName]; ok && style.StyleName == "conditional" {
			return fmt.Errorf("layout: conditional style %q is not implemented", el.StyleName)
		}
		if (el.Kind == jrxml.KindImage || el.Kind == jrxml.KindTextField) && strings.TrimSpace(el.Expr) == "" {
			return fmt.Errorf("layout: %s expression is required", el.Kind)
		}
		for _, src := range []string{el.Expr, el.PrintWhenExpr} {
			if strings.TrimSpace(src) != "" {
				program, err := expression.Compile(src)
				if err != nil {
					return err
				}
				n.Programs[src] = program
			}
		}
		for _, child := range el.Children {
			if err := visit(child); err != nil {
				return err
			}
		}
		return nil
	}
	seen := map[jrxml.BandKind]bool{}
	for _, b := range t.Bands {
		if seen[b.Kind] {
			return nil, fmt.Errorf("layout: duplicate %s band", b.Kind)
		}
		seen[b.Kind] = true
		for _, el := range b.Elements {
			if err := visit(el); err != nil {
				return nil, err
			}
		}
	}
	for _, g := range t.Groups {
		for _, b := range append(append([]jrxml.Band{}, g.Headers...), g.Footers...) {
			for _, el := range b.Elements {
				if err := visit(el); err != nil {
					return nil, err
				}
			}
		}
	}
	return &DetailSpec{Native: n}, nil
}

type nativePlacement struct {
	Page   int
	X, Y   printdoc.Unit
	Band   jrxml.Band
	Record fill.FilledRecord
}

func layoutNative(ctx context.Context, filled *fill.Filled, spec *NativeSpec, budget limits.Budget) (*printdoc.Document, error) {
	t := spec.Template
	p := t.Page
	bands := map[jrxml.BandKind]jrxml.Band{}
	for _, b := range t.Bands {
		bands[b.Kind] = b
	}
	h := func(kind jrxml.BandKind) printdoc.Unit { return bands[kind].Height }
	top := p.MarginTop
	bottom := p.Height - p.MarginBottom - h(jrxml.BandPageFooter) - h(jrxml.BandColumnFooter)
	if p.Width <= 0 || p.Height <= 0 || bottom <= top+h(jrxml.BandPageHeader)+h(jrxml.BandColumnHeader) {
		return nil, fmt.Errorf("layout: invalid page geometry")
	}
	guard := budget.Guard(ctx)
	doc := &printdoc.Document{}
	var placements []nativePlacement
	empty := fill.FilledRecord{Fields: map[string]expression.Value{}, Parameters: filled.Parameters, Variables: map[string]expression.Value{}}
	current := empty
	if len(filled.Records) > 0 {
		current = filled.Records[0]
	}
	y := top
	page := 0
	add := func(kind jrxml.BandKind, at printdoc.Unit, record fill.FilledRecord) {
		if b, ok := bands[kind]; ok {
			placements = append(placements, nativePlacement{page, p.MarginLeft, at, b, record})
		}
	}
	newPage := func(record fill.FilledRecord) error {
		if err := guard.Page(); err != nil {
			return err
		}
		if err := doc.AppendPage(printdoc.Page{Width: p.Width, Height: p.Height}); err != nil {
			return err
		}
		page = len(doc.Pages)
		y = top
		add(jrxml.BandBackground, top, record)
		if page == 1 {
			if title, ok := bands[jrxml.BandTitle]; ok {
				if y+title.Height+h(jrxml.BandPageHeader)+h(jrxml.BandColumnHeader) > bottom {
					return fmt.Errorf("layout: title exceeds page")
				}
				add(jrxml.BandTitle, y, record)
				y += title.Height
			}
		}
		add(jrxml.BandPageHeader, y, record)
		y += h(jrxml.BandPageHeader)
		add(jrxml.BandColumnHeader, y, record)
		y += h(jrxml.BandColumnHeader)
		return nil
	}
	finish := func(record fill.FilledRecord) {
		add(jrxml.BandColumnFooter, bottom, record)
		add(jrxml.BandPageFooter, p.Height-p.MarginBottom-h(jrxml.BandPageFooter), record)
	}
	if len(filled.Records) == 0 {
		if h(jrxml.BandNoData) > 0 {
			if err := newPage(empty); err != nil {
				return nil, err
			}
			add(jrxml.BandNoData, y, empty)
			finish(empty)
		}
	} else {
		if err := newPage(current); err != nil {
			return nil, err
		}

		detail := bands[jrxml.BandDetail]
		if detail.Height < 0 || detail.Height > bottom-top-h(jrxml.BandPageHeader)-h(jrxml.BandColumnHeader) {
			return nil, fmt.Errorf("layout: detail exceeds page")
		}
		placeGroupBand := func(band jrxml.Band, record fill.FilledRecord) error {
			if band.Height > bottom-top-h(jrxml.BandPageHeader)-h(jrxml.BandColumnHeader) {
				return fmt.Errorf("layout: group band exceeds page")
			}
			if y+band.Height > bottom {
				finish(current)
				if err := newPage(record); err != nil {
					return err
				}
			}
			placements = append(placements, nativePlacement{page, p.MarginLeft, y, band, record})
			y += band.Height
			return nil
		}
		for i, record := range filled.Records {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			firstChanged := len(t.Groups)
			for j, g := range t.Groups {
				if i == 0 || !reflect.DeepEqual(record.Groups[g.Name], current.Groups[g.Name]) {
					firstChanged = j
					break
				}
			}
			if i > 0 {
				for j := len(t.Groups) - 1; j >= firstChanged; j-- {
					for _, band := range t.Groups[j].Footers {
						if err := placeGroupBand(band, current); err != nil {
							return nil, err
						}
					}
				}
			}
			for j := firstChanged; j < len(t.Groups); j++ {
				for _, band := range t.Groups[j].Headers {
					if err := placeGroupBand(band, record); err != nil {
						return nil, err
					}
				}
			}
			if y+detail.Height > bottom {
				finish(current)
				if err := newPage(record); err != nil {
					return nil, err
				}
			}
			add(jrxml.BandDetail, y, record)
			y += detail.Height
			current = record
		}
		for j := len(t.Groups) - 1; j >= 0; j-- {
			for _, band := range t.Groups[j].Footers {
				if err := placeGroupBand(band, current); err != nil {
					return nil, err
				}
			}
		}
		if summary, ok := bands[jrxml.BandSummary]; ok && summary.Height > 0 {
			if summary.Height > bottom-top-h(jrxml.BandPageHeader)-h(jrxml.BandColumnHeader) {
				return nil, fmt.Errorf("layout: summary exceeds page")
			}
			if y+summary.Height > bottom {
				finish(current)
				if err := newPage(current); err != nil {
					return nil, err
				}
			}
			add(jrxml.BandSummary, y, current)
		}
		finish(current)
	}
	total := len(doc.Pages)
	for _, place := range placements {
		// Last-page footer replaces the normal footer.
		if place.Page == total && place.Band.Kind == jrxml.BandPageFooter {
			if last, ok := bands[jrxml.BandLastPageFooter]; ok {
				place.Band = last
				place.Y = p.Height - p.MarginBottom - last.Height
			}
		}
		for _, el := range place.Band.Elements {
			if err := renderNativeElement(ctx, doc, spec, place, el, place.X, place.Y, total, filled, guard); err != nil {
				return nil, err
			}
		}
	}
	if err := doc.Validate(); err != nil {
		return nil, err
	}
	return doc, nil
}

func renderNativeElement(ctx context.Context, doc *printdoc.Document, spec *NativeSpec, place nativePlacement, el jrxml.Element, x, y printdoc.Unit, total int, filled *fill.Filled, guard *limits.Guard) error {
	if err := guard.Step(1); err != nil {
		return err
	}
	record := place.Record
	if el.EvalTime == "Report" && len(filled.Records) > 0 {
		record = filled.Records[len(filled.Records)-1]
	}
	variables := map[string]expression.Value{}
	for name, value := range record.Variables {
		variables[name] = value
	}
	number := place.Page
	if el.EvalTime == "Report" {
		number = total
	}
	variables["PAGE_NUMBER"] = expression.Int(number)
	record.Variables = variables
	scope := expression.Scope{record}
	if el.PrintWhenExpr != "" {
		v, err := spec.Programs[el.PrintWhenExpr].Eval(ctx, scope)
		if err != nil {
			return err
		}
		if v != expression.Bool(true) {
			return nil
		}
	}
	bounds := printdoc.Bounds{X: x + el.X, Y: y + el.Y, Width: el.Width, Height: el.Height}
	if bounds.Width < 0 || bounds.Height < 0 {
		return fmt.Errorf("layout: negative element size")
	}
	appendElement := func(element printdoc.Element) {
		doc.Pages[place.Page-1].Elements = append(doc.Pages[place.Page-1].Elements, element)
	}
	if el.Opaque {
		shape := printdoc.NewShapeElement(bounds, printdoc.ShapeRectangle)
		shape.Shape.Fill = true
		shape.Shape.FillColor = el.Backcolor
		appendElement(shape)
	}
	switch el.Kind {
	case jrxml.KindFrame:
		for _, child := range el.Children {
			if err := renderNativeElement(ctx, doc, spec, place, child, bounds.X, bounds.Y, total, filled, guard); err != nil {
				return err
			}
		}
		return nil
	case jrxml.KindImage:
		v, err := spec.Programs[el.Expr].Eval(ctx, scope)
		if err != nil {
			return err
		}
		if v == nil || v == expression.NullValue {
			return nil
		}
		imageValue, ok := v.(expression.Image)
		if !ok {
			return fmt.Errorf("layout: image expression requires an image parameter")
		}
		var pngData bytes.Buffer
		if err := png.Encode(&pngData, imageValue.Value); err != nil {
			return err
		}
		appendElement(printdoc.NewImageElement(bounds, "data:image/png;base64,"+base64.StdEncoding.EncodeToString(pngData.Bytes())))
		return nil
	case jrxml.KindLine, jrxml.KindRectangle, jrxml.KindEllipse:
		kind := printdoc.ShapeLine
		if el.Kind == jrxml.KindRectangle {
			kind = printdoc.ShapeRectangle
		} else if el.Kind == jrxml.KindEllipse {
			kind = printdoc.ShapeEllipse
		}
		shape := printdoc.NewShapeElement(bounds, kind)
		shape.Shape.Stroke = true
		shape.Shape.StrokeColor = el.Forecolor
		shape.Shape.StrokeWidth = el.LineWidth
		if shape.Shape.StrokeWidth == 0 {
			shape.Shape.StrokeWidth = printdoc.FromPoints(1)
		}
		appendElement(shape)
		return nil
	}
	value := el.Text
	if el.Kind == jrxml.KindTextField {
		v, err := spec.Programs[el.Expr].Eval(ctx, scope)
		if err != nil {
			return err
		}
		value = valueString(v)
		if el.Pattern != "" {
			formatted, err := expression.FormatValue(v, el.Pattern)
			if err != nil {
				return err
			}
			value = formatted
		}
	}
	if el.Border {
		shape := printdoc.NewShapeElement(bounds, printdoc.ShapeRectangle)
		shape.Shape.Stroke = true
		shape.Shape.StrokeColor = el.Forecolor
		shape.Shape.StrokeWidth = el.LineWidth
		appendElement(shape)
	}
	bounds.X += el.Padding.Left
	bounds.Y += el.Padding.Top
	bounds.Width -= el.Padding.Left + el.Padding.Right
	bounds.Height -= el.Padding.Top + el.Padding.Bottom
	size := el.FontSize
	if size == 0 {
		size = 10
	}
	flags := text.Style(0)
	if el.Bold {
		flags |= text.Bold
	}
	if el.Italic {
		flags |= text.Italic
	}
	font, err := text.Resolve(el.Font, flags, size)
	if err != nil {
		return err
	}
	lines, err := text.Wrap(value, font, bounds.Width)
	if err != nil {
		return err
	}
	metrics, err := text.Measure("", font)
	if err != nil {
		return err
	}
	style := printdoc.DefaultStyle()
	style.Font = font.Name
	style.Size = printdoc.FromPoints(size)
	style.Bold = el.Bold
	style.Italic = el.Italic
	style.Color = el.Forecolor
	textY := bounds.Y
	lineHeight := metrics.LineHeight
	if el.VerticalAlignment == "Middle" {
		textY += (bounds.Height - printdoc.Unit(len(lines))*lineHeight) / 2
	} else if el.VerticalAlignment == "Bottom" {
		textY += bounds.Height - printdoc.Unit(len(lines))*lineHeight
	}
	for _, line := range lines {
		if textY+metrics.Ascent+metrics.Descent > bounds.Y+bounds.Height {
			break
		}
		m, err := text.Measure(line, font)
		if err != nil {
			return err
		}
		textX := bounds.X
		if el.TextAlignment == "Center" {
			textX += (bounds.Width - m.Width) / 2
		} else if el.TextAlignment == "Right" {
			textX += bounds.Width - m.Width
		}
		element := printdoc.NewTextElement(printdoc.Bounds{X: textX, Y: textY, Width: m.Width, Height: metrics.Ascent + metrics.Descent}, line)
		element.Style = style
		element.GridX = bounds.X
		element.HasGridX = true
		appendElement(element)
		textY += lineHeight
	}
	return nil
}

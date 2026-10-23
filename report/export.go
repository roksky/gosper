package report

import (
	"archive/zip"
	"context"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"sort"
	"strings"

	"github.com/roksky/gosper/internal/printdoc"
	"github.com/roksky/gosper/internal/text"
)

// exportDocument exports the supported positioned-text print model. Spreadsheet
// and CSV rows are grouped by vertical position; columns use global X positions.
func exportDocument(ctx context.Context, doc *printdoc.Document, w io.Writer, format Format) error {
	if w == nil {
		return fmt.Errorf("report: nil writer")
	}
	if err := doc.Validate(); err != nil {
		return err
	}
	var b strings.Builder
	if format == FormatHTML {
		b.WriteString(`<!doctype html><html><head><meta charset="utf-8"><title>Report</title><style>html,body{margin:0;padding:0;background:#fff}section{box-sizing:border-box;overflow:hidden;isolation:isolate}@media print{section{break-after:page}section:last-child{break-after:auto}}</style></head><body>`)
		for _, p := range doc.Pages {
			if err := ctx.Err(); err != nil {
				return err
			}
			fmt.Fprintf(&b, `<section style="position:relative;width:%gpt;height:%gpt;page-break-after:always">`, p.Width.Points(), p.Height.Points())
			for _, e := range p.Elements {
				if e.Kind != printdoc.ElementText {
					if err := htmlGraphic(&b, e); err != nil {
						return err
					}
					continue
				}
				if err := htmlText(&b, e); err != nil {
					return err
				}
			}
			b.WriteString("</section>")
		}
		b.WriteString("</body></html>")
	} else if format == FormatXML {
		// Jasper print XML, rather than the JRXML source or raw query data.
		b.WriteString(xml.Header + `<jasperPrint name="report">`)
		for _, p := range doc.Pages {
			if err := ctx.Err(); err != nil {
				return err
			}
			b.WriteString("<page>")
			for _, e := range p.Elements {
				if e.Kind != printdoc.ElementText {
					fmt.Fprintf(&b, `<%s x="%g" y="%g" width="%g" height="%g"/>`, e.Kind, e.Bounds.X.Points(), e.Bounds.Y.Points(), e.Bounds.Width.Points(), e.Bounds.Height.Points())
					continue
				}
				fmt.Fprintf(&b, `<text><reportElement x="%g" y="%g" width="%g" height="%g"/><textContent>`, e.Bounds.X.Points(), e.Bounds.Y.Points(), e.Bounds.Width.Points(), e.Bounds.Height.Points())
				if err := xml.EscapeText(&b, []byte(e.Text)); err != nil {
					return err
				}
				b.WriteString("</textContent></text>")
			}
			b.WriteString("</page>")
		}
		b.WriteString("</jasperPrint>")
	} else {
		rows, err := grid(ctx, doc)
		if err != nil {
			return err
		}
		if format == FormatCSV {
			cw := csv.NewWriter(w)
			return cw.WriteAll(rows)
		}
		return writeXLSX(ctx, w, rows)
	}
	_, err := io.WriteString(w, b.String())
	return err
}

func grid(ctx context.Context, doc *printdoc.Document) ([][]string, error) {
	xs := map[printdoc.Unit]bool{}
	for _, p := range doc.Pages {
		for _, e := range p.Elements {
			if e.Kind != printdoc.ElementText {
				continue
			}
			xs[gridX(e)] = true
		}
	}
	var columns []printdoc.Unit
	for x := range xs {
		columns = append(columns, x)
	}
	sort.Slice(columns, func(i, j int) bool { return columns[i] < columns[j] })
	index := map[printdoc.Unit]int{}
	for i, x := range columns {
		index[x] = i
	}
	var result [][]string
	for _, p := range doc.Pages {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		byY := map[printdoc.Unit][]string{}
		var ys []printdoc.Unit
		for _, e := range p.Elements {
			if e.Kind != printdoc.ElementText {
				continue
			}
			y := e.Bounds.Y
			if _, ok := byY[y]; !ok {
				byY[y] = make([]string, len(columns))
				ys = append(ys, y)
			}
			col := index[gridX(e)]
			if byY[y][col] != "" {
				byY[y][col] += "\n"
			}
			byY[y][col] += e.Text
		}
		sort.Slice(ys, func(i, j int) bool { return ys[i] < ys[j] })
		for _, y := range ys {
			result = append(result, byY[y])
		}
	}
	return result, nil
}

func writeXLSX(ctx context.Context, w io.Writer, rows [][]string) error {
	z := zip.NewWriter(w)
	files := []struct{ name, body string }{
		{"[Content_Types].xml", `<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/><Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/></Types>`},
		{"_rels/.rels", `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/></Relationships>`},
		{"xl/workbook.xml", `<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets><sheet name="Report" sheetId="1" r:id="rId1"/></sheets></workbook>`},
		{"xl/_rels/workbook.xml.rels", `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>`},
	}
	for _, f := range files {
		out, err := z.Create(f.name)
		if err != nil {
			return err
		}
		if _, err = io.WriteString(out, xml.Header+f.body); err != nil {
			return err
		}
	}
	out, err := z.Create("xl/worksheets/sheet1.xml")
	if err != nil {
		return err
	}
	var b strings.Builder
	b.WriteString(xml.Header + `<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>`)
	for i, row := range rows {
		if err := ctx.Err(); err != nil {
			return err
		}
		fmt.Fprintf(&b, `<row r="%d">`, i+1)
		for j, value := range row {
			col := ""
			for n := j + 1; n > 0; n = (n - 1) / 26 {
				col = string(rune('A'+(n-1)%26)) + col
			}
			fmt.Fprintf(&b, `<c r="%s%d" t="inlineStr"><is><t xml:space="preserve">`, col, i+1)
			if err := xml.EscapeText(&b, []byte(value)); err != nil {
				return err
			}
			b.WriteString("</t></is></c>")
		}
		b.WriteString("</row>")
	}
	b.WriteString("</sheetData></worksheet>")
	if _, err = io.WriteString(out, b.String()); err != nil {
		return err
	}
	return z.Close()
}

func htmlGraphic(b *strings.Builder, e printdoc.Element) error {
	bounds := e.Bounds
	style := fmt.Sprintf("position:absolute;left:%gpt;top:%gpt;width:%gpt;height:%gpt", bounds.X.Points(), bounds.Y.Points(), bounds.Width.Points(), bounds.Height.Points())
	switch e.Kind {
	case printdoc.ElementImage:
		if !strings.HasPrefix(e.Image.Resource, "data:image/png;base64,") {
			return fmt.Errorf("report: unsupported image resource")
		}
		fmt.Fprintf(b, `<img alt="" style="%s;display:block;max-width:none;object-fit:contain" src="%s">`, style, html.EscapeString(e.Image.Resource))
	case printdoc.ElementShape:
		shape := e.Shape
		fill := "none"
		if shape.Fill {
			fill = fmt.Sprintf("#%06x", shape.FillColor)
		}
		stroke := "none"
		if shape.Stroke || !shape.Fill {
			stroke = fmt.Sprintf("#%06x", shape.StrokeColor)
		}
		width := shape.StrokeWidth.Points()
		if width == 0 {
			width = 1
		}
		w, h := bounds.Width.Points(), bounds.Height.Points()
		fmt.Fprintf(b, `<svg xmlns="http://www.w3.org/2000/svg" style="%s;overflow:visible" viewBox="0 0 %g %g" aria-hidden="true">`, style, max(w, 1), max(h, 1))
		switch shape.Kind {
		case printdoc.ShapeLine:
			fmt.Fprintf(b, `<line x1="0" y1="0" x2="%g" y2="%g" stroke="%s" stroke-width="%g"/>`, w, h, stroke, width)
		case printdoc.ShapeEllipse:
			fmt.Fprintf(b, `<ellipse cx="%g" cy="%g" rx="%g" ry="%g" fill="%s" stroke="%s" stroke-width="%g"/>`, w/2, h/2, w/2, h/2, fill, stroke, width)
		case printdoc.ShapeRectangle:
			fmt.Fprintf(b, `<rect width="%g" height="%g" fill="%s" stroke="%s" stroke-width="%g"/>`, w, h, fill, stroke, width)
		}
		b.WriteString("</svg>")
	default:
		return fmt.Errorf("report: unsupported HTML element %s", e.Kind)
	}
	return nil
}

func gridX(e printdoc.Element) printdoc.Unit {
	if e.HasGridX {
		return e.GridX
	}
	return e.Bounds.X
}

// Layout already split and positioned each line. SVG supplies the same text
// baseline as PDF without allowing the browser to wrap the line a second time.
func htmlText(b *strings.Builder, e printdoc.Element) error {
	s := e.EffectiveStyle()
	flags := text.Style(0)
	weight, italic := "normal", "normal"
	if s.Bold {
		flags |= text.Bold
		weight = "bold"
	}
	if s.Italic {
		flags |= text.Italic
		italic = "italic"
	}
	font, err := text.Resolve(s.Font, flags, s.Size.Points())
	if err != nil {
		return err
	}
	metrics, err := text.Measure(e.Text, font)
	if err != nil {
		return err
	}
	x := s.Padding.Left
	available := e.Bounds.Width - s.Padding.Left - s.Padding.Right
	if available < 0 {
		return fmt.Errorf("report: padding exceeds HTML element width")
	}
	if available > metrics.Width {
		switch s.Align {
		case printdoc.AlignCenter:
			x += (available - metrics.Width) / 2
		case printdoc.AlignRight:
			x += available - metrics.Width
		}
	}
	family := map[string]string{"SansSerif": "Arial,Helvetica,sans-serif", "Serif": "Times New Roman,Times,serif", "Monospaced": "Courier New,Courier,monospace"}[s.Font]
	fmt.Fprintf(b, `<svg xmlns="http://www.w3.org/2000/svg" style="position:absolute;overflow:visible;left:%gpt;top:%gpt;width:%gpt;height:%gpt" viewBox="0 0 %g %g"><text x="%g" y="%g" xml:space="preserve" style="white-space:pre;font-family:%s;font-size:%gpx;font-weight:%s;font-style:%s;fill:#%06x">%s</text></svg>`, e.Bounds.X.Points(), e.Bounds.Y.Points(), max(e.Bounds.Width.Points(), 1), max(e.Bounds.Height.Points(), 1), max(e.Bounds.Width.Points(), 1), max(e.Bounds.Height.Points(), 1), x.Points(), (s.Padding.Top + metrics.Ascent).Points(), family, s.Size.Points(), weight, italic, s.Color, html.EscapeString(e.Text))
	return nil
}

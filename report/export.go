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
		b.WriteString("<!doctype html><html><head><meta charset=\"utf-8\"><title>Report</title></head><body>")
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
				s := e.EffectiveStyle()
				weight := "normal"
				if s.Bold {
					weight = "bold"
				}
				italic := "normal"
				if s.Italic {
					italic = "italic"
				}
				align := []string{"left", "center", "right"}[s.Align]
				font := map[string]string{"SansSerif": "sans-serif", "Serif": "serif", "Monospaced": "monospace"}[s.Font]
				fmt.Fprintf(&b, `<div style="position:absolute;white-space:pre-wrap;left:%gpt;top:%gpt;width:%gpt;height:%gpt;font-family:%s;font-size:%gpt;font-weight:%s;font-style:%s;text-align:%s;color:#%06x">%s</div>`, e.Bounds.X.Points(), e.Bounds.Y.Points(), e.Bounds.Width.Points(), e.Bounds.Height.Points(), font, s.Size.Points(), weight, italic, align, s.Color, html.EscapeString(e.Text))
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
		fmt.Fprintf(b, `<img alt="" style="%s;object-fit:contain" src="%s">`, style, html.EscapeString(e.Image.Resource))
	case printdoc.ElementShape:
		shape := e.Shape
		if shape.Kind == printdoc.ShapeEllipse {
			style += ";border-radius:50%"
		}
		if shape.Fill {
			style += fmt.Sprintf(";background:#%06x", shape.FillColor)
		}
		if shape.Stroke {
			style += fmt.Sprintf(";border:%gpt solid #%06x", shape.StrokeWidth.Points(), shape.StrokeColor)
		}
		fmt.Fprintf(b, `<div style="%s"></div>`, style)
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

package exporter

import (
	"context"
	"fmt"
	"github.com/roksky/gosper/internal/printdoc"
	"github.com/roksky/gosper/internal/text"
	"io"
	"strconv"
	"strings"
)

type countedWriter struct {
	w io.Writer
	n int64
}

func (c *countedWriter) write(s string) error {
	n, err := io.WriteString(c.w, s)
	c.n += int64(n)
	if err != nil {
		return err
	}
	if n != len(s) {
		return io.ErrShortWrite
	}
	return nil
}

type pdfFont struct {
	key  string
	base string
}

// PDF streams a valid PDF 1.4 document to w. It supports text elements using
// the PDF core Latin fonts and printable ASCII text; images, shapes, and font
// embedding belong to later exporter tasks. No creation timestamp is emitted.
func PDF(ctx context.Context, doc *printdoc.Document, w io.Writer) error {
	if ctx == nil || w == nil {
		return fmt.Errorf("pdf: nil context or writer")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := doc.Validate(); err != nil {
		return err
	}
	fonts, index, err := scanFonts(doc)
	if err != nil {
		return err
	}
	pageCount := len(doc.Pages)
	fontBase := 3 + 3*pageCount
	objectCount := fontBase + len(fonts) - 1
	offsets := make([]int64, objectCount+1)
	out := &countedWriter{w: w}
	write := func(s string) error { return out.write(s) }
	object := func(id int) error { offsets[id] = out.n; return write(fmt.Sprintf("%d 0 obj\n", id)) }
	if err := write("%PDF-1.4\n%\xE2\xE3\xCF\xD3\n"); err != nil {
		return err
	}
	if err := object(1); err != nil {
		return err
	}
	if err := write("<< /Type /Catalog /Pages 2 0 R >>\nendobj\n"); err != nil {
		return err
	}
	if err := object(2); err != nil {
		return err
	}
	var kids strings.Builder
	for i := 0; i < pageCount; i++ {
		fmt.Fprintf(&kids, " %d 0 R", 3+3*i)
	}
	if err := write(fmt.Sprintf("<< /Type /Pages /Kids [%s ] /Count %d >>\nendobj\n", kids.String(), pageCount)); err != nil {
		return err
	}
	for i, page := range doc.Pages {
		if err := ctx.Err(); err != nil {
			return err
		}
		pageID, streamID, lengthID := 3+3*i, 4+3*i, 5+3*i
		if err := object(pageID); err != nil {
			return err
		}
		var resources strings.Builder
		for j := range fonts {
			fmt.Fprintf(&resources, " /F%d %d 0 R", j+1, fontBase+j)
		}
		shell := fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %s %s] /Resources << /Font <<%s >> >> /Contents %d 0 R >>\nendobj\n", pdfUnit(page.Width), pdfUnit(page.Height), resources.String(), streamID)
		if err := write(shell); err != nil {
			return err
		}
		if err := object(streamID); err != nil {
			return err
		}
		if err := write(fmt.Sprintf("<< /Length %d 0 R >>\nstream\n", lengthID)); err != nil {
			return err
		}
		start := out.n
		for _, elem := range page.Elements {
			if err := ctx.Err(); err != nil {
				return err
			}
			if elem.Kind != printdoc.ElementText {
				if err := writeGraphic(out, page, elem); err != nil {
					return err
				}
				continue
			}
			if err := writeText(out, page, elem, index); err != nil {
				return err
			}
		}
		length := out.n - start
		if err := write("endstream\nendobj\n"); err != nil {
			return err
		}
		if err := object(lengthID); err != nil {
			return err
		}
		if err := write(fmt.Sprintf("%d\nendobj\n", length)); err != nil {
			return err
		}
	}
	for i, f := range fonts {
		if err := object(fontBase + i); err != nil {
			return err
		}
		if err := write(fmt.Sprintf("<< /Type /Font /Subtype /Type1 /BaseFont /%s /Encoding /WinAnsiEncoding >>\nendobj\n", f.base)); err != nil {
			return err
		}
	}
	xref := out.n
	if err := write(fmt.Sprintf("xref\n0 %d\n0000000000 65535 f \n", objectCount+1)); err != nil {
		return err
	}
	for id := 1; id <= objectCount; id++ {
		if offsets[id] > 9999999999 {
			return fmt.Errorf("pdf: object offset exceeds PDF xref range")
		}
		if err := write(fmt.Sprintf("%010d 00000 n \n", offsets[id])); err != nil {
			return err
		}
	}
	return write(fmt.Sprintf("trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", objectCount+1, xref))
}
func scanFonts(doc *printdoc.Document) ([]pdfFont, map[string]int, error) {
	fonts := []pdfFont{}
	index := map[string]int{}
	for _, page := range doc.Pages {
		for _, elem := range page.Elements {
			if elem.Kind != printdoc.ElementText {
				if elem.Kind == printdoc.ElementShape || elem.Kind == printdoc.ElementImage {
					continue
				}
				return nil, nil, fmt.Errorf("pdf: unsupported element kind %s", elem.Kind)
			}
			for _, r := range elem.Text {
				if r < 32 || r > 126 {
					return nil, nil, fmt.Errorf("pdf: text outside printable ASCII subset")
				}
			}
			style := elem.EffectiveStyle()
			if err := style.Validate(); err != nil {
				return nil, nil, err
			}
			base := baseFont(style)
			if _, ok := index[base]; !ok {
				fonts = append(fonts, pdfFont{key: base, base: base})
				index[base] = len(fonts)
			}
		}
	}
	return fonts, index, nil
}
func baseFont(s printdoc.Style) string {
	switch s.Font {
	case "Serif":
		if s.Bold && s.Italic {
			return "Times-BoldItalic"
		}
		if s.Bold {
			return "Times-Bold"
		}
		if s.Italic {
			return "Times-Italic"
		}
		return "Times-Roman"
	case "Monospaced":
		if s.Bold && s.Italic {
			return "Courier-BoldOblique"
		}
		if s.Bold {
			return "Courier-Bold"
		}
		if s.Italic {
			return "Courier-Oblique"
		}
		return "Courier"
	default:
		if s.Bold && s.Italic {
			return "Helvetica-BoldOblique"
		}
		if s.Bold {
			return "Helvetica-Bold"
		}
		if s.Italic {
			return "Helvetica-Oblique"
		}
		return "Helvetica"
	}
}
func writeText(out *countedWriter, page printdoc.Page, elem printdoc.Element, index map[string]int) error {
	style := elem.EffectiveStyle()
	flags := text.Style(0)
	if style.Bold {
		flags |= text.Bold
	}
	if style.Italic {
		flags |= text.Italic
	}
	font, err := text.Resolve(style.Font, flags, style.Size.Points())
	if err != nil {
		return err
	}
	metrics, err := text.Measure(elem.Text, font)
	if err != nil {
		return err
	}
	x := elem.Bounds.X + style.Padding.Left
	available := elem.Bounds.Width - style.Padding.Left - style.Padding.Right
	if available < 0 {
		return fmt.Errorf("pdf: padding exceeds element width")
	}
	switch style.Align {
	case printdoc.AlignCenter:
		if available > metrics.Width {
			x += (available - metrics.Width) / 2
		}
	case printdoc.AlignRight:
		if available > metrics.Width {
			x += available - metrics.Width
		}
	}
	y := page.Height - elem.Bounds.Y - style.Padding.Top - metrics.Ascent
	if y < 0 {
		return fmt.Errorf("pdf: text baseline outside page")
	}
	color := style.Color
	red, green, blue := float64((color>>16)&0xff)/255, float64((color>>8)&0xff)/255, float64(color&0xff)/255
	command := fmt.Sprintf("BT /F%d %s Tf %.4f %.4f %.4f rg 1 0 0 1 %s %s Tm (%s) Tj ET\n", index[baseFont(style)], pdfUnit(style.Size), red, green, blue, pdfUnit(x), pdfUnit(y), escapePDF(elem.Text))
	return out.write(command)
}
func pdfUnit(u printdoc.Unit) string { return strconv.FormatFloat(u.Points(), 'f', 2, 64) }
func escapePDF(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == '(' || r == ')' || r == '\\' {
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

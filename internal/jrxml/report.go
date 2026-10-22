package jrxml

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/roksky/gosper/internal/diag"
	"github.com/roksky/gosper/internal/printdoc"
)

// PageProperties holds the report-level geometry declared on <jasperReport>.
//
// Numeric attributes in JRXML are interpreted as points (72 dpi), the unit
// JasperReports uses for page and element geometry, and converted to twips via
// the printdoc helpers. Fields left at their defaults reflect the
// JasperReports JRXML defaults (A4 portrait, single 555pt column, 30pt
// top/bottom and 20pt left/right margins).
type PageProperties struct {
	Width         printdoc.Unit
	Height        printdoc.Unit
	Orientation   string
	ColumnCount   int
	ColumnWidth   printdoc.Unit
	ColumnSpacing printdoc.Unit
	MarginTop     printdoc.Unit
	MarginBottom  printdoc.Unit
	MarginLeft    printdoc.Unit
	MarginRight   printdoc.Unit
}

func defaultPage() PageProperties {
	return PageProperties{
		Width:         printdoc.FromPoints(595),
		Height:        printdoc.FromPoints(842),
		Orientation:   "Portrait",
		ColumnCount:   1,
		ColumnWidth:   printdoc.FromPoints(555),
		ColumnSpacing: printdoc.FromPoints(0),
		MarginTop:     printdoc.FromPoints(30),
		MarginBottom:  printdoc.FromPoints(30),
		MarginLeft:    printdoc.FromPoints(20),
		MarginRight:   printdoc.FromPoints(20),
	}
}

func init() { registerRootStart(parseReportAttrs) }

// parseReportAttrs reads the report-level attributes from <jasperReport>.
func parseReportAttrs(p *parser, se xml.StartElement) {
	for _, a := range se.Attr {
		switch a.Name.Local {
		case "name":
			p.tmpl.Name = a.Value
		case "language":
			p.tmpl.Language = a.Value
		case "pageWidth":
			p.tmpl.Page.Width = p.points(a)
		case "pageHeight":
			p.tmpl.Page.Height = p.points(a)
		case "orientation":
			p.tmpl.Page.Orientation = a.Value
		case "columnCount":
			p.tmpl.Page.ColumnCount = p.intAttr(a)
		case "columnWidth":
			p.tmpl.Page.ColumnWidth = p.points(a)
		case "columnSpacing":
			p.tmpl.Page.ColumnSpacing = p.points(a)
		case "topMargin":
			p.tmpl.Page.MarginTop = p.points(a)
		case "bottomMargin":
			p.tmpl.Page.MarginBottom = p.points(a)
		case "leftMargin":
			p.tmpl.Page.MarginLeft = p.points(a)
		case "rightMargin":
			p.tmpl.Page.MarginRight = p.points(a)
		}
	}
}

// points converts a numeric attribute in points to the internal unit, emitting
// a validation diagnostic when it is not a number.
func (p *parser) points(a xml.Attr) printdoc.Unit {
	return p.pointsValue(a.Name.Local, a.Value)
}

func (p *parser) pointsValue(name, raw string) printdoc.Unit {
	v, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		p.invalidAttr(name, raw, "number")
		return 0
	}
	return printdoc.FromPoints(v)
}

func (p *parser) intAttr(a xml.Attr) int {
	v, err := strconv.Atoi(strings.TrimSpace(a.Value))
	if err != nil {
		p.invalidAttr(a.Name.Local, a.Value, "integer")
		return 0
	}
	return v
}

func (p *parser) invalidAttr(name, raw, want string) {
	p.diags.Add(diag.New(
		diag.CodeValidation,
		diag.SeverityWarning,
		fmt.Sprintf("invalid %s value %q for attribute %q", want, raw, name),
		p.span(),
	))
}

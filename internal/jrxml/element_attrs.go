package jrxml

import (
	"encoding/xml"
	"fmt"

	"github.com/roksky/gosper/internal/diag"
	"github.com/roksky/gosper/internal/printdoc"
)

// elementAttrs holds the geometry and evaluation settings shared by all
// drawable elements. In JRXML the geometry lives on the nested <reportElement>
// child; the evaluation settings live on the element itself. Both locations are
// accepted.
type elementAttrs struct {
	X             printdoc.Unit
	Y             printdoc.Unit
	Width         printdoc.Unit
	Height        printdoc.Unit
	Stretch       bool
	PrintWhenExpr string
	EvalTime      string
	EvalGroup     string
}

// evalTimes are the valid JasperReports evaluationTime values.
var evalTimes = map[string]bool{
	"Now":    true,
	"Report": true,
	"Page":   true,
	"Column": true,
	"Group":  true,
	"Auto":   true,
	"Band":   true,
}

func defaultElementAttrs() elementAttrs {
	return elementAttrs{EvalTime: "Now"}
}

// applyElementStartAttrs reads geometry/evaluation attributes from an element or
// <reportElement> start tag.
func (p *parser) applyElementStartAttrs(se xml.StartElement, el *Element) {
	for _, a := range se.Attr {
		switch a.Name.Local {
		case "x":
			el.X = p.points(a)
		case "y":
			el.Y = p.points(a)
		case "width":
			el.Width = p.points(a)
		case "height":
			el.Height = p.points(a)
		case "isStretchWithOverflow":
			el.Stretch = parseBool(a.Value)
		case "printWhenExpression":
			el.PrintWhenExpr = a.Value
		case "evaluationTime":
			p.setEvalTime(a.Value, el)
		case "evaluationGroup":
			el.EvalGroup = a.Value
		case "uuid":
			// Ignored: identity is not needed before layout.
		}
	}
}

func (p *parser) setEvalTime(v string, el *Element) {
	el.EvalTime = v
	if !evalTimes[v] {
		p.diags.Add(newAttrDiag(p, "unknown evaluationTime %q", v))
	}
}

// parseReportElement consumes a <reportElement> child: geometry attributes plus
// an optional nested <printWhenExpression>.
func (p *parser) parseReportElement(se xml.StartElement, el *Element) error {
	p.applyElementStartAttrs(se, el)
	for {
		tok, err := p.dec.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "printWhenExpression":
				v, err := p.readText()
				if err != nil {
					return err
				}
				el.PrintWhenExpr = v
			default:
				p.unsupported(t, t.Name.Local)
				if err := p.skipElement(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			return nil
		}
	}
}

// parseElementCommon parses the start-tag attributes and children shared by all
// drawable elements: reportElement, text and expression children, and
// unsupported diagnostics for anything else.
func (p *parser) parseElementCommon(se xml.StartElement, el *Element) error {
	p.applyElementStartAttrs(se, el)
	for {
		tok, err := p.dec.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "reportElement":
				if err := p.parseReportElement(t, el); err != nil {
					return err
				}
			case "text":
				v, err := p.readText()
				if err != nil {
					return err
				}
				el.Text = v
			case "textFieldExpression":
				v, err := p.readText()
				if err != nil {
					return err
				}
				el.Expr = v
			case "printWhenExpression":
				v, err := p.readText()
				if err != nil {
					return err
				}
				el.PrintWhenExpr = v
			default:
				p.unsupported(t, t.Name.Local)
				if err := p.skipElement(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			return nil
		}
	}
}

func newAttrDiag(p *parser, format string, args ...any) *diag.Diagnostic {
	return diag.New(diag.CodeValidation, diag.SeverityWarning, fmt.Sprintf(format, args...), p.span())
}

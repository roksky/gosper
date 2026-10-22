package jrxml

import (
	"encoding/xml"
	"fmt"
)

// Parameter is a report parameter declaration (<parameter>).
//
// Type is the declared Java class name (for example "java.lang.String"). It
// defaults to "java.lang.String" when the class attribute is absent. DefaultExpr
// is the raw text of the optional <defaultValueExpression>; it is stored
// unevaluated and compiled later.
type Parameter struct {
	Name        string
	Type        string
	NestedType  string
	DefaultExpr string
}

const defaultParameterType = "java.lang.String"

func init() { registerRootChild("parameter", parseParameter) }

func parseParameter(p *parser, se xml.StartElement) error {
	prm := Parameter{Type: defaultParameterType}
	for _, a := range se.Attr {
		switch a.Name.Local {
		case "name":
			prm.Name = a.Value
		case "class":
			prm.Type = a.Value
		case "nestedType":
			prm.NestedType = a.Value
		}
	}

	for {
		tok, err := p.dec.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "defaultValueExpression":
				v, err := p.readText()
				if err != nil {
					return err
				}
				prm.DefaultExpr = v
			default:
				p.unsupported(t, t.Name.Local)
				if err := p.skipElement(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			p.appendParameter(prm)
			return nil
		}
	}
}

func (p *parser) appendParameter(prm Parameter) {
	for _, existing := range p.tmpl.Parameters {
		if existing.Name == prm.Name {
			p.warnValidation(fmt.Sprintf("duplicate parameter %q", prm.Name))
			break
		}
	}
	p.tmpl.Parameters = append(p.tmpl.Parameters, prm)
}

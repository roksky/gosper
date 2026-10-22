package jrxml

import (
	"encoding/xml"
	"fmt"
)

type Variable struct{ Name, Type, Calculation, ResetType, ResetGroup, Expr, InitialExpr string }
type Group struct {
	Name, Expr       string
	Headers, Footers []Band
}

func init() { registerRootChild("variable", parseVariable); registerRootChild("group", parseGroup) }
func parseVariable(p *parser, se xml.StartElement) error {
	v := Variable{Calculation: "Nothing", ResetType: "Report"}
	v.Name, _ = attribute(se, "name")
	v.Type, _ = attribute(se, "class")
	if value, ok := attribute(se, "calculation"); ok {
		v.Calculation = value
	}
	if value, ok := attribute(se, "resetType"); ok {
		v.ResetType = value
	}
	v.ResetGroup, _ = attribute(se, "resetGroup")
	if v.Name == "" {
		return fmt.Errorf("jrxml: variable name required")
	}
	for {
		tok, err := p.dec.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "variableExpression":
				v.Expr, err = p.readText()
			case "initialValueExpression":
				v.InitialExpr, err = p.readText()
			default:
				p.unsupported(t, t.Name.Local)
				err = p.skipElement()
			}
			if err != nil {
				return err
			}
		case xml.EndElement:
			p.tmpl.Variables = append(p.tmpl.Variables, v)
			return nil
		}
	}
}
func parseGroup(p *parser, se xml.StartElement) error {
	g := Group{}
	g.Name, _ = attribute(se, "name")
	if g.Name == "" {
		return fmt.Errorf("jrxml: group name required")
	}
	for {
		tok, err := p.dec.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "groupExpression":
				g.Expr, err = p.readText()
			case "groupHeader", "groupFooter":
				var bands []Band
				bands, err = p.parseGroupBands(t)
				if t.Name.Local == "groupHeader" {
					g.Headers = bands
				} else {
					g.Footers = bands
				}
			default:
				p.unsupported(t, t.Name.Local)
				err = p.skipElement()
			}
			if err != nil {
				return err
			}
		case xml.EndElement:
			p.tmpl.Groups = append(p.tmpl.Groups, g)
			return nil
		}
	}
}
func (p *parser) parseGroupBands(se xml.StartElement) ([]Band, error) {
	var bands []Band
	for {
		tok, err := p.dec.Token()
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local != "band" {
				return nil, fmt.Errorf("jrxml: expected group band")
			}
			b, err := p.parseBandElement(BandKind(se.Name.Local), t)
			if err != nil {
				return nil, err
			}
			bands = append(bands, b)
		case xml.EndElement:
			return bands, nil
		}
	}
}

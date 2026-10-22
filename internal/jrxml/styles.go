package jrxml

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

func init() {
	registerRootChild("style", parseStyle)
	// Unreferenced datasets have no effect on report execution. Dataset runs and
	// components still produce unsupported diagnostics rather than losing rows.
	registerRootChild("subDataset", func(p *parser, se xml.StartElement) error { return p.skipElement() })
}
func parseStyle(p *parser, se xml.StartElement) error {
	name, _ := attribute(se, "name")
	if name == "" {
		return fmt.Errorf("jrxml: style name required")
	}
	el := Element{Backcolor: 0xffffff}
	p.applyElementStartAttrs(se, &el)
	for _, a := range se.Attr {
		switch a.Name.Local {
		case "fontName":
			el.Font = a.Value
		case "fontSize":
			el.FontSize, _ = strconv.ParseFloat(a.Value, 64)
		case "isBold":
			el.Bold = parseBool(a.Value)
		case "isItalic":
			el.Italic = parseBool(a.Value)
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
			case "box":
				err = p.parseBox(t, &el)
			case "conditionalStyle":
				// Preserve diagnostics for conditional styles that are actually referenced.
				// The parser retains the declaration as an unsupported style marker.
				el.StyleName = "conditional"
				err = p.skipElement()
			default:
				p.unsupported(t, t.Name.Local)
				err = p.skipElement()
			}
			if err != nil {
				return err
			}
		case xml.EndElement:
			if p.tmpl.Styles == nil {
				p.tmpl.Styles = map[string]Element{}
			}
			p.tmpl.Styles[name] = el
			return nil
		}
	}
}
func (p *parser) parseBox(se xml.StartElement, el *Element) error {
	for _, a := range se.Attr {
		value := p.points(a)
		switch a.Name.Local {
		case "padding":
			el.Padding.Top = value
			el.Padding.Bottom = value
			el.Padding.Left = value
			el.Padding.Right = value
		case "leftPadding":
			el.Padding.Left = value
		case "rightPadding":
			el.Padding.Right = value
		case "topPadding":
			el.Padding.Top = value
		case "bottomPadding":
			el.Padding.Bottom = value
		}
	}
	for {
		tok, err := p.dec.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local == "pen" {
				p.readPen(t, el)
				el.Border = true
			} else if t.Name.Local != "topPen" && t.Name.Local != "leftPen" && t.Name.Local != "rightPen" && t.Name.Local != "bottomPen" {
				p.unsupported(t, t.Name.Local)
			}
			if err := p.skipElement(); err != nil {
				return err
			}
		case xml.EndElement:
			return nil
		}
	}
}
func (p *parser) readPen(se xml.StartElement, el *Element) {
	if value, ok := attribute(se, "lineWidth"); ok {
		el.LineWidth = p.pointsValue("lineWidth", value)
	}
	if value, ok := attribute(se, "lineColor"); ok {
		el.Forecolor = parseColor(value)
	}
}
func (p *parser) parseGraphic(se xml.StartElement, el *Element) error {
	for {
		tok, err := p.dec.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local == "pen" {
				p.readPen(t, el)
			} else {
				p.unsupported(t, t.Name.Local)
			}
			if err := p.skipElement(); err != nil {
				return err
			}
		case xml.EndElement:
			return nil
		}
	}
}

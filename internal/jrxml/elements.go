package jrxml

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// ElementKind identifies a drawable element.
type ElementKind string

const (
	// KindStaticText is <staticText>.
	KindStaticText ElementKind = "staticText"
	// KindTextField is <textField>.
	KindTextField ElementKind = "textField"
)

// Element is a drawable element parsed from a band.
//
// Text holds staticText content; Expr holds a textFieldExpression; both are
// stored raw and unevaluated. Coordinates are in the internal unit (see
// PageProperties).
type Element struct {
	Kind ElementKind

	elementAttrs

	Text            string
	Expr            string
	Pattern         string
	IsBlankWhenNull bool
}

func init() {
	registerChildElement("staticText", parseStaticText)
	registerChildElement("textField", parseTextField)
}

func parseStaticText(p *parser, se xml.StartElement) (Element, error) {
	el := Element{Kind: KindStaticText, elementAttrs: defaultElementAttrs()}
	if err := p.parseElementCommon(se, &el); err != nil {
		return Element{}, err
	}
	return el, nil
}

func parseTextField(p *parser, se xml.StartElement) (Element, error) {
	el := Element{Kind: KindTextField, elementAttrs: defaultElementAttrs()}
	for _, a := range se.Attr {
		switch a.Name.Local {
		case "pattern":
			el.Pattern = a.Value
		case "isBlankWhenNull":
			el.IsBlankWhenNull = parseBool(a.Value)
		}
	}
	if err := p.parseElementCommon(se, &el); err != nil {
		return Element{}, err
	}
	return el, nil
}

// parseBool parses an XML boolean, treating "true" (case-insensitive) as true.
// JasperReports uses the strings "true"/"false".
func parseBool(v string) bool {
	b, err := strconv.ParseBool(strings.TrimSpace(v))
	return err == nil && b
}

package jrxml

import "encoding/xml"

func init() {
	for _, kind := range []ElementKind{KindFrame, KindImage, KindLine, KindRectangle, KindEllipse} {
		k := kind
		registerChildElement(string(k), func(p *parser, se xml.StartElement) (Element, error) {
			el := Element{Kind: k, elementAttrs: defaultElementAttrs(), Backcolor: 0xffffff}
			err := p.parseElementCommon(se, &el)
			return el, err
		})
	}
	registerRootChild("property", func(p *parser, se xml.StartElement) error { return p.skipElement() })
}

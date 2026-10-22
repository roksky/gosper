package jrxml

import (
	"encoding/xml"
	"fmt"

	"github.com/roksky/gosper/internal/printdoc"
)

// BandKind identifies a report section that carries drawable elements.
type BandKind string

// Core band kinds (JRX-008). Advanced bands (background, lastPageFooter,
// noData) and group bands are later tasks (JRX-009, JRX-015).
const (
	BandTitle        BandKind = "title"
	BandPageHeader   BandKind = "pageHeader"
	BandColumnHeader BandKind = "columnHeader"
	BandDetail       BandKind = "detail"
	BandColumnFooter BandKind = "columnFooter"
	BandPageFooter   BandKind = "pageFooter"
	BandSummary      BandKind = "summary"
)

// coreBands lists the sections parsed in declaration order.
var coreBands = []BandKind{
	BandTitle, BandPageHeader, BandColumnHeader, BandDetail,
	BandColumnFooter, BandPageFooter, BandSummary,
}

// Band is one report section with its height and elements.
type Band struct {
	Kind     BandKind
	Height   printdoc.Unit
	Elements []Element
}

func init() {
	for _, k := range coreBands {
		kind := k
		registerRootChild(string(kind), func(p *parser, se xml.StartElement) error {
			return p.parseSection(kind, se)
		})
	}
}

// parseSection parses a section element (for example <detail>). JasperReports
// nests a <band> child that carries the height and elements; a section with
// direct element children and a height attribute is also accepted.
func (p *parser) parseSection(kind BandKind, se xml.StartElement) error {
	if p.seenSections[kind] {
		p.warnValidation(fmt.Sprintf("duplicate band <%s>", kind))
	}
	p.seenSections[kind] = true

	var direct *Band
	for {
		tok, err := p.dec.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local == "band" {
				b, err := p.parseBandElement(kind, t)
				if err != nil {
					return err
				}
				p.tmpl.Bands = append(p.tmpl.Bands, b)
				continue
			}
			el, ok, err := p.parseElement(t)
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
			if direct == nil {
				direct = &Band{Kind: kind}
				if v, ok := attribute(se, "height"); ok {
					direct.Height = p.pointsValue("height", v)
				}
			}
			direct.Elements = append(direct.Elements, el)
		case xml.EndElement:
			if direct != nil {
				p.tmpl.Bands = append(p.tmpl.Bands, *direct)
			}
			return nil
		}
	}
}

// parseBandElement parses a <band> child: its height attribute and elements.
func (p *parser) parseBandElement(kind BandKind, se xml.StartElement) (Band, error) {
	b := Band{Kind: kind}
	if v, ok := attribute(se, "height"); ok {
		b.Height = p.pointsValue("height", v)
	}
	for {
		tok, err := p.dec.Token()
		if err != nil {
			return Band{}, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			el, ok, err := p.parseElement(t)
			if err != nil {
				return Band{}, err
			}
			if ok {
				b.Elements = append(b.Elements, el)
			}
		case xml.EndElement:
			return b, nil
		}
	}
}

// parseElement parses one drawable element. Unknown elements produce an
// unsupported diagnostic and are skipped; ok reports whether an Element was
// produced.
func (p *parser) parseElement(se xml.StartElement) (Element, bool, error) {
	if fn, ok := childElementParsers[se.Name.Local]; ok {
		el, err := fn(p, se)
		if err != nil {
			return Element{}, false, err
		}
		return el, true, nil
	}
	p.unsupported(se, se.Name.Local)
	if err := p.skipElement(); err != nil {
		return Element{}, false, err
	}
	return Element{}, false, nil
}

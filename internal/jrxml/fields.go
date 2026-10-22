package jrxml

import "encoding/xml"

// Field is a report field declaration (<field>).
//
// Type is the declared Java class name. Unknown class names are stored verbatim;
// mapping them to engine types is deferred to DAT-007.
type Field struct {
	Name        string
	Type        string
	Description string
}

func init() { registerRootChild("field", parseField) }

func parseField(p *parser, se xml.StartElement) error {
	f := Field{}
	for _, a := range se.Attr {
		switch a.Name.Local {
		case "name":
			f.Name = a.Value
		case "class":
			f.Type = a.Value
		case "description":
			f.Description = a.Value
		}
	}
	if err := p.skipElement(); err != nil {
		return err
	}
	p.tmpl.Fields = append(p.tmpl.Fields, f)
	return nil
}

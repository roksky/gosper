package jrxml

import (
	"encoding/xml"
	"fmt"
	"strings"
)

func init() { registerRootChild("queryString", parseQuery) }

func parseQuery(p *parser, se xml.StartElement) error {
	language, _ := attribute(se, "language")
	if language != "" && !strings.EqualFold(language, "SQL") {
		return fmt.Errorf("jrxml: unsupported query language %q", language)
	}
	query, err := p.readText()
	if err != nil {
		return err
	}
	if p.tmpl.Query != "" {
		return fmt.Errorf("jrxml: duplicate queryString")
	}
	p.tmpl.Query = strings.TrimSpace(query)
	return nil
}

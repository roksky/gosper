package jrxml

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// Dialect identifies the JRXML dialect declared by a template's root element.
//
// Raw holds the version string exactly as written (for example "6.20.0"). Major
// and Minor are parsed from it. A zero Major means the version was absent or
// unparsable; use Known to test.
type Dialect struct {
	Major int
	Minor int
	Raw   string
}

// Known reports whether a usable major version was detected.
func (d Dialect) Known() bool { return d.Major > 0 }

// String returns the raw version string, or "unknown" when none was declared.
func (d Dialect) String() string {
	if d.Raw == "" {
		return "unknown"
	}
	return d.Raw
}

// Template is the top-level parsed representation of a JRXML template.
//
// The zero value is not ready for use; Parse always populates Page with the
// JasperReports default page properties.
type Template struct {
	Dialect    Dialect
	Name       string
	Language   string
	Page       PageProperties
	Parameters []Parameter
	Fields     []Field
	Bands      []Band
}

// parseDialect parses a version string such as "6.20.0", "6.0" or "6".
// An empty or unparsable string yields a Dialect with Major 0 (unknown).
func parseDialect(raw string) Dialect {
	d := Dialect{Raw: raw}
	if raw == "" {
		return d
	}
	parts := strings.SplitN(raw, ".", 3)
	major, err := strconv.Atoi(parts[0])
	if err != nil || major <= 0 {
		return Dialect{Raw: raw}
	}
	d.Major = major
	if len(parts) > 1 {
		if minor, err := strconv.Atoi(parts[1]); err == nil {
			d.Minor = minor
		}
	}
	return d
}

// attribute returns the value of the named attribute (matched by local name)
// and whether it was present.
func attribute(se xml.StartElement, name string) (string, bool) {
	for _, a := range se.Attr {
		if a.Name.Local == name {
			return a.Value, true
		}
	}
	return "", false
}

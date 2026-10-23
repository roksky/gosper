package jrxml

import (
	"encoding/xml"
	"fmt"

	"github.com/roksky/gosper/internal/diag"
)

// deferredElements are constructs recognized by name but deliberately not
// implemented in this milestone. They are reported as unsupported (with a
// pointer to the owning feature) rather than silently ignored or treated as
// parse errors; their subtrees are skipped.
var deferredElements = map[string]string{
	"subreport": "subreports (JRX-017)",
	"chart":     "charts (JRX-024, EXT-001)",
	"crosstab":  "crosstabs (JRX-024, EXT-002)",
	"barcode":   "barcodes (JRX-024, EXT-003)",
	"component": "report components (JRX-024)",
}

// unsupported records an unsupported-construct diagnostic. The message names
// the element and, for known-but-deferred constructs, points at the owning
// feature. The span is always non-nil.
func (p *parser) unsupported(se xml.StartElement, name string) {
	msg := fmt.Sprintf("unsupported element <%s>", name)
	if reason, ok := deferredElements[name]; ok {
		msg = fmt.Sprintf("unsupported element <%s>: %s", name, reason)
	}
	p.diags.Add(diag.New(diag.CodeUnsupported, diag.SeverityWarning, msg, p.span()))
}

// Package jrxml parses JasperReports-style JRXML templates from XML into an
// in-memory template model.
//
// Parse is total: it never panics on malformed input. Fatal problems are
// returned as an error that wraps a *diag.Diagnostic with code diag.CodeParse;
// non-fatal conditions (an unknown or missing version, an unsupported element,
// an invalid attribute, or a duplicate declaration) are returned as warnings in
// the diag.List while the template is still produced.
//
// # Handler seam
//
// Parsing is driven by a token-stream dispatcher. Constructs register handlers
// from their own files via registerRootChild (children of <jasperReport> such as
// <parameter>, <field> and bands) and registerChildElement (drawable elements
// such as <staticText> and <textField>), and report-level attributes via
// registerRootStart. Adding a construct therefore means adding a file, not
// editing the dispatcher; unrecognized elements fall through to an unsupported
// diagnostic (see unsupported.go).
//
// # Units and defaults
//
// Numeric geometry in JRXML is interpreted as points (72 dpi) and converted to
// the printdoc twip unit. Absent page attributes take the JasperReports JRXML
// defaults: A4 portrait (595x842pt), one 555pt column, 30pt top/bottom and 20pt
// left/right margins. Text and expression content is stored verbatim (no
// trimming) and evaluated later.
//
// Later tasks add variables, queries, styles, groups, images and version
// normalization (JRX-005, JRX-006, JRX-007, JRX-011, JRX-012, JRX-015, JRX-020).
package jrxml

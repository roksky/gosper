// Package jrxml parses JasperReports-style JRXML templates from XML into an
// in-memory template model.
//
// This package currently defines the parser contract: the Parse entry point,
// the top-level Template type, and the Dialect (version) representation.
// Element, band, field and parameter parsing are added by later tasks
// (JRX-002..JRX-019); normalization is JRX-020.
//
// Parse is total: it never panics on malformed input. Fatal problems are
// returned as an error that wraps a *diag.Diagnostic with code diag.CodeParse;
// non-fatal conditions (for example an unknown or missing version) are returned
// as warnings in the diag.List while the template is still produced.
package jrxml

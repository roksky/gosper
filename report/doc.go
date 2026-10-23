// Package report is the public, stable API of the gosper reporting engine.
//
// It defines the surface a caller uses to compile a JRXML template, fill it
// with data, and export a document:
//
//	template -> Compile -> *CompiledReport
//	*CompiledReport + data -> Fill -> *FilledReport
//	*FilledReport -> Export -> output
//
// Compile prepares JRXML bands, expressions, groups, and parameter defaults.
// Query binds embedded SQL for a caller-owned database. Fill accepts rows and
// image parameters; Export supports PDF, HTML, CSV, XLSX and print XML.
// Unsupported constructs return diagnostics rather than executing Java code.
//
// The package lives at the module root (import path
// "github.com/roksky/gosper/report"). See docs/adr/0005-public-api-package.md.
//
// # Concurrency
//
// A *CompiledReport is immutable and safe for concurrent Fill calls. A
// *FilledReport owns per-run output and must not be mutated during export.
package report

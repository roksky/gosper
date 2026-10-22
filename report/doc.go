// Package report is the public, stable API of the gosper reporting engine.
//
// It defines the surface a caller uses to compile a JRXML template, fill it
// with data, and export a document:
//
//	template -> Compile -> *CompiledReport
//	*CompiledReport + data -> Fill -> *FilledReport
//	*FilledReport -> Export -> output
//
// The M1 path compiles supported JRXML text bands, fills them from a data
// source, and exports a PDF. Other formats return ErrNotImplemented.
//
// The package lives at the module root (import path
// "github.com/roksky/gosper/report"). See docs/adr/0005-public-api-package.md.
//
// # Concurrency
//
// A *CompiledReport is immutable and safe for concurrent Fill calls. A
// *FilledReport owns per-run output and must not be mutated during export.
package report

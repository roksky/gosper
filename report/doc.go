// Package report is the public, stable API of the gosper reporting engine.
//
// It defines the surface a caller uses to compile a JRXML template, fill it
// with data, and export a document:
//
//	template -> Compile -> *CompiledReport
//	*CompiledReport + data -> Fill -> *FilledReport
//	*FilledReport -> Export -> output
//
// # Status
//
// This package currently contains the API contract only: exported types and
// entry-point signatures. Every entry point returns ErrNotImplemented until the
// corresponding lifecycle task lands (FND-001-002 and later). No report
// behavior is implemented here.
//
// The package lives at the module root (import path
// "github.com/roksky/gosper/report"). See docs/adr/0005-public-api-package.md.
//
// # Concurrency
//
// A *CompiledReport is intended to be immutable and safe for concurrent use. A
// *FilledReport is per-run and must not be shared across goroutines. These
// guarantees become effective when the lifecycle is implemented.
package report

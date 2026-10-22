// Package expression parses and evaluates report expressions used by compiled
// templates, independent of any JVM-based implementation.
//
// This package currently defines the expression value and symbol contracts
// (Value, Symbols, Kind, Scope); the grammar and evaluator are deferred to
// EXP-003. Evaluation is sandboxed: expressions are interpreted by a closed
// implementation and never execute arbitrary Go or Java code. See the sandbox
// policy in docs/architecture/expressions-and-types.md.
package expression

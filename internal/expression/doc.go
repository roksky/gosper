// Package expression parses and evaluates report expressions used by compiled
// templates, independent of any JVM-based implementation.
//
// The package defines the expression value and symbol contracts and a bounded
// parser and evaluator. Evaluation is sandboxed: expressions use a closed
// implementation and never execute arbitrary Go or Java code. See the sandbox
// policy in docs/architecture/expressions-and-types.md.
package expression

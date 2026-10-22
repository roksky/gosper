// Package printdoc defines the engine-neutral print document model that
// exporters consume, decoupling layout from output formats.
//
// A print document is a paginated, positioned representation produced by
// layout. It is independent of JRXML, expressions and data sources, and uses a
// single internal unit: the twip (1/1440 inch), stored as int64 to avoid
// floating-point drift. See docs/adr/0004-print-ir-unit-system.md.
package printdoc

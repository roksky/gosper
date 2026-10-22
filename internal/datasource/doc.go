// Package datasource defines pluggable sources of business data that feed the
// fill stage.
//
// It provides the DataSource/Cursor/Record contract and slice, JSON, and CSV sources for
// structs, string-keyed maps, and positional rows. SQL sources are deferred to later changes.
package datasource

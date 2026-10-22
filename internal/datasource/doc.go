// Package datasource defines pluggable sources of business data that feed the
// fill stage.
//
// It provides the DataSource/Cursor/Record contract and a slice source for
// structs, string-keyed maps, and positional rows. CSV and SQL sources are deferred to later changes.
package datasource

// Package datasource defines pluggable sources of business data that feed the
// fill stage.
//
// This package currently provides only the DataSource/Cursor/Record contract
// and its lifecycle rules (see datasource.go). Concrete sources — slices,
// JSON, CSV and SQL — are deferred to later changes (DAT-001, DAT-003,
// DAT-004, DAT-005).
package datasource

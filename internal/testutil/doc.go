// Package testutil provides deterministic helpers for tests: a fixed clock, a
// pinned location and locale, and stable map iteration order. It exists so
// engine tests do not depend on wall-clock time, the host timezone or Go map
// order. Golden compares byte-exact fixtures and updates them only with the
// explicit -update test flag.
package testutil

package testutil

import (
	"sort"
	"time"
)

// DefaultLocation is the timezone tests must assume.
var DefaultLocation = time.UTC

// DefaultLocale is the locale tests must assume for number/date formatting.
const DefaultLocale = "en-US"

// FixedClock returns a clock that always reports t. Use it to make
// time-dependent output deterministic.
func FixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

// UTC returns the UTC location, the default location for tests.
func UTC() *time.Location { return time.UTC }

// SortedKeys returns the keys of m in ascending order, giving tests a stable
// iteration order regardless of insertion order.
func SortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

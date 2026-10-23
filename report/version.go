package report

import "strings"

// version is the shared build identifier. Override it at link time with
// -X github.com/roksky/gosper/report.version=<value>.
var version = "0.1.0-beta.1"

// Version reports the library and CLI build identifier. An empty override
// falls back to dev.
func Version() string {
	v := strings.TrimSpace(version)
	if v == "" {
		return "dev"
	}
	return v
}

// Command gosper is the optional command-line entry point for the gosper
// reporting engine.
//
// In this bootstrap milestone the command only exposes help and version
// metadata. It performs no report parsing, compilation, filling or exporting.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/roksky/gosper/report"
)

const usage = `gosper - a Go-native, JasperReports-compatible reporting engine.

CLI report commands are not implemented yet.

Usage:
  gosper [command]

Commands:
  help       Show this help text.
  version    Show version information.

Flags:
  -h, --help      Show this help text.
  -v, --version   Show version information.

This build exposes CLI metadata only; no report commands are available.
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run executes the CLI with the given arguments and writers, returning the
// process exit code. It is separated from main so it can be tested directly.
func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stdout, usage)
		return 0
	}

	switch args[0] {
	case "-h", "--help", "help":
		fmt.Fprint(stdout, usage)
		return 0
	case "-v", "--version", "version":
		fmt.Fprintf(stdout, "gosper %s\n", report.Version())
		return 0
	default:
		fmt.Fprintf(stderr, "gosper: unknown command %q\n\n", args[0])
		fmt.Fprint(stderr, usage)
		return 2
	}
}

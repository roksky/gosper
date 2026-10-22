package main

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/roksky/gosper/report"
)

func TestRunHelp(t *testing.T) {
	cases := [][]string{
		nil,
		{"-h"},
		{"--help"},
		{"help"},
	}

	for _, args := range cases {
		var stdout, stderr bytes.Buffer
		if code := run(args, &stdout, &stderr); code != 0 {
			t.Errorf("run(%q) = %d, want 0", args, code)
		}
		if !strings.Contains(stdout.String(), "Usage:") {
			t.Errorf("run(%q) stdout missing usage: %q", args, stdout.String())
		}
		if stderr.Len() != 0 {
			t.Errorf("run(%q) stderr = %q, want empty", args, stderr.String())
		}
	}
}

func TestRunVersion(t *testing.T) {
	cases := [][]string{
		{"-v"},
		{"--version"},
		{"version"},
	}

	for _, args := range cases {
		var stdout, stderr bytes.Buffer
		if code := run(args, &stdout, &stderr); code != 0 {
			t.Errorf("run(%q) = %d, want 0", args, code)
		}
		if got, want := stdout.String(), "gosper "+report.Version()+"\n"; got != want {
			t.Errorf("run(%q) stdout = %q, want %q", args, got, want)
		}
		if stderr.Len() != 0 {
			t.Errorf("run(%q) stderr = %q, want empty", args, stderr.String())
		}
	}
}

func TestRunUnknownCommand(t *testing.T) {
	cases := [][]string{
		{"frobnicate"},
		{"compile"},
		{"--bogus"},
	}

	for _, args := range cases {
		var stdout, stderr bytes.Buffer
		if code := run(args, &stdout, &stderr); code == 0 {
			t.Errorf("run(%q) = 0, want nonzero", args)
		}
		if !strings.Contains(stderr.String(), "unknown command") {
			t.Errorf("run(%q) stderr missing error: %q", args, stderr.String())
		}
	}
}

func TestVersionDefault(t *testing.T) {
	if report.Version() == "" {
		t.Fatal("version must not be empty")
	}
}

func TestLinkedVersionOverride(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "gosper")
	build := exec.Command("go", "build", "-ldflags", "-X github.com/roksky/gosper/report.version=9.8.7", "-o", bin, ".")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}
	output, err := exec.Command(bin, "--version").CombinedOutput()
	if err != nil || string(output) != "gosper 9.8.7\n" {
		t.Fatalf("linked CLI version = %q, %v", output, err)
	}
}

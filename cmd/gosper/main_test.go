package main

import (
	"bytes"
	"strings"
	"testing"
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
		if got, want := stdout.String(), "gosper "+version+"\n"; got != want {
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
	if version == "" {
		t.Fatal("version must not be empty")
	}
}

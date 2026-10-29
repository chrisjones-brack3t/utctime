package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"
)

// fixed is a non-UTC "now" so the UTC conversion is exercised too.
var fixed = time.Date(2026, 10, 29, 16, 40, 25, 0, time.FixedZone("EDT", -4*60*60))

func TestParseAgo(t *testing.T) {
	tests := []struct {
		in      string
		want    time.Duration
		wantErr bool
	}{
		{"1h", time.Hour, false},
		{"90m", 90 * time.Minute, false},
		{"1d", 24 * time.Hour, false},
		{"1.5d", 36 * time.Hour, false},
		{"xd", 0, true},
		{"foo", 0, true},
	}
	for _, tc := range tests {
		got, err := parseAgo(tc.in)
		if (err != nil) != tc.wantErr {
			t.Errorf("parseAgo(%q) error = %v, wantErr %v", tc.in, err, tc.wantErr)
		}
		if got != tc.want {
			t.Errorf("parseAgo(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantCode   int
		wantStdout string
		wantStderr string // substring
	}{
		{"now in UTC", nil, 0, "2026-10-29T20:40:25Z\n", ""},
		{"one hour ago", []string{"--ago", "1h"}, 0, "2026-10-29T19:40:25Z\n", ""},
		{"single dash works", []string{"-ago", "1h"}, 0, "2026-10-29T19:40:25Z\n", ""},
		{"one day ago", []string{"--ago", "1d"}, 0, "2026-10-28T20:40:25Z\n", ""},
		{"ninety minutes ago", []string{"--ago", "90m"}, 0, "2026-10-29T19:10:25Z\n", ""},
		{"help long", []string{"--help"}, 0, "", "Usage:"},
		{"help short", []string{"-h"}, 0, "", "Usage:"},
		{"invalid duration", []string{"--ago", "foo"}, 2, "", "error: --ago needs"},
		{"negative duration", []string{"--ago", "-1h"}, 2, "", "error: --ago needs"},
		{"unknown flag", []string{"--bogus"}, 2, "", "flag provided but not defined"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run(tc.args, fixed, &stdout, &stderr)
			if code != tc.wantCode {
				t.Errorf("exit code = %d, want %d", code, tc.wantCode)
			}
			if stdout.String() != tc.wantStdout {
				t.Errorf("stdout = %q, want %q", stdout.String(), tc.wantStdout)
			}
			if !strings.Contains(stderr.String(), tc.wantStderr) {
				t.Errorf("stderr = %q, want it to contain %q", stderr.String(), tc.wantStderr)
			}
		})
	}
}

func TestMainExitsZero(t *testing.T) {
	oldArgs, oldExit := os.Args, exit
	defer func() { os.Args, exit = oldArgs, oldExit }()
	os.Args = []string{"utctime", "--ago", "1h"}
	code := -1
	exit = func(c int) { code = c }
	main()
	if code != 0 {
		t.Errorf("main exit code = %d, want 0", code)
	}
}

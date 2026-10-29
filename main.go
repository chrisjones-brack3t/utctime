// utctime prints a UTC timestamp in RFC 3339 format, optionally offset into the past.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

// parseAgo parses a duration like "90m", "1h", or "2d".
// WHY: time.ParseDuration has no day unit, so "d" is handled here; other units fall through to stdlib.
func parseAgo(s string) (time.Duration, error) {
	if days, ok := strings.CutSuffix(s, "d"); ok {
		n, err := strconv.ParseFloat(days, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid duration %q", s)
		}
		return time.Duration(n * float64(24*time.Hour)), nil
	}
	return time.ParseDuration(s)
}

const usage = `utctime - print a UTC timestamp (RFC 3339), now or some time ago.

Usage:
  utctime [--ago DURATION]

With no flags, prints the current time converted to UTC.

Options:
  --ago DURATION   Print the time DURATION in the past, e.g. 30m, 1h, 1d, 1.5d.
                   Units: ns, us, ms, s, m, h, d. Must be positive.
  --version        Print the version.
  -h, --help       Show this help.

Examples:
  utctime               2026-10-29T20:40:25Z
  utctime --ago 1h      one hour ago
  utctime --ago 1d      one day ago
`

// exit and readBuildInfo are variables so tests can replace them.
var (
	exit          = os.Exit
	readBuildInfo = debug.ReadBuildInfo
)

// version returns the module version Go embedded at build time.
// WHY: the git tag is the single source of truth; no version constant to keep in sync.
func version() string {
	if bi, ok := readBuildInfo(); ok {
		return bi.Main.Version
	}
	return "unknown"
}

func main() {
	exit(run(os.Args[1:], time.Now(), os.Stdout, os.Stderr))
}

// run does the work of main and returns the process exit code.
// WHY: taking args, clock and writers as parameters makes it testable without a subprocess.
func run(args []string, now time.Time, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("utctime", flag.ContinueOnError)
	fs.SetOutput(stderr)
	// WHY: the flag package already treats -h, --help and -help as ErrHelp; no extra flag needed.
	fs.Usage = func() { fmt.Fprint(stderr, usage) }
	ago := fs.String("ago", "", "go back in time by `DURATION` instead of printing the current time")
	showVersion := fs.Bool("version", false, "print the version")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	if *showVersion {
		fmt.Fprintln(stdout, "utctime", version())
		return 0
	}

	t := now
	if *ago != "" {
		d, err := parseAgo(*ago)
		if err != nil || d < 0 {
			fmt.Fprintf(stderr, "error: --ago needs a positive duration like 1h or 1d, got %q\n", *ago)
			return 2
		}
		t = t.Add(-d)
	}
	fmt.Fprintln(stdout, t.UTC().Format(time.RFC3339))
	return 0
}

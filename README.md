# utctime

Print the current time, or a time in the past, as a UTC timestamp.

```
$ utctime
2026-10-29T20:52:36Z
$ utctime --ago 1h
2026-10-29T19:52:36Z
```

Output is [RFC 3339](https://www.rfc-editor.org/rfc/rfc3339) in UTC (`YYYY-MM-DDTHH:MM:SSZ`), one line, so it drops straight into API queries, log searches, and SQL filters.

## Install

Requires Go (see `go.mod` for the version).

```sh
go install github.com/chrisjones-brack3t/utctime@latest
```

Or from a clone of this repo:

```sh
go install .
```

Either way, `utctime` lands in `$(go env GOPATH)/bin`; make sure that directory is on your `PATH`.

## Usage

```
utctime [--ago DURATION]
```

| Flag | Description |
| --- | --- |
| `--ago DURATION` | Print the time `DURATION` in the past. Must be positive. |
| `--version` | Print the version, e.g. `utctime v0.1.0`. |
| `-h`, `--help` | Show help. |

With no flags, `utctime` converts your system's current time to UTC.

### Durations

A number followed by a unit: `ns`, `us`, `ms`, `s`, `m`, `h`, `d`. Decimals are allowed.

```sh
utctime --ago 30m    # 30 minutes ago
utctime --ago 1h     # 1 hour ago
utctime --ago 1d     # 1 day ago
utctime --ago 1.5d   # 36 hours ago
```

Combined values such as `1d12h` are not supported; use `36h` or `1.5d`.

### Exit codes

| Code | Meaning |
| --- | --- |
| `0` | Timestamp printed, or help shown. |
| `2` | Invalid or negative `--ago` value, or unknown flag. The error goes to stderr. |

## Development

```sh
go test -cover ./...
```

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the program
go run main.go

# Build
go build -o currency-tracker .

# Run tests
go test ./...

# Run a single test
go test -run TestFunctionName ./...

# Lint (if golangci-lint is available)
golangci-lint run
```

## Architecture

This is a minimal single-file Go program (`main.go`) that demonstrates concurrent HTTP requests using goroutines, channels, and `sync.WaitGroup`.

**Core flow:**
1. `main()` defines a list of named sources (simulating multiple banks), creates a buffered channel and a `WaitGroup`, then spawns one goroutine per source via `go FetchRate(...)`.
2. `FetchRate` hits the public `open.er-api.com` API (no key required), parses the JSON response, and sends a `RateResponse` into the channel before decrementing the `WaitGroup`.
3. A separate goroutine calls `wg.Wait()` then `close(ch)` so `main` can drain the channel with `for range ch` without deadlocking.

**Key types:**
- `RateResponse` — carries the result (or error) from a single fetch; the `Error` field uses `json:"-"` so it is never serialized.

**Concurrency pattern:** fan-out (one goroutine per source) → buffered channel → sequential drain in `main`. The channel buffer equals `len(sources)` to prevent goroutines from blocking before the closer goroutine runs.

The module name is `currency-tracker` (see `go.mod`); there are no external dependencies.

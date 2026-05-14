# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Context

Educational Go program demonstrating concurrent HTTP requests. Source code comments are in Traditional Chinese throughout, making this a learning resource for Go concurrency patterns. The program fetches USD/TWD exchange rates from a public REST API across multiple simulated bank sources simultaneously.

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

Single-file Go program (`main.go`) with zero external dependencies — stdlib only. Demonstrates the fan-out concurrency pattern using goroutines, a buffered channel, and `sync.WaitGroup`.

**Core flow:**
1. `main()` defines a hardcoded list of named sources (`Bank_A`, `Bank_B`, `Bank_C`), creates a buffered channel sized to `len(sources)`, then spawns one goroutine per source via `go FetchRate(...)`.
2. `FetchRate(source, base, target, ch chan<- RateResponse, wg *sync.WaitGroup)` hits the public `open.er-api.com` API (no key required), parses the JSON response, and sends a `RateResponse` into the channel before calling `wg.Done()`.
3. A separate goroutine calls `wg.Wait()` then `close(ch)` so `main` can drain the channel with `for range ch` without deadlocking.

**Key types:**
- `RateResponse` — carries the result (or error) from a single fetch. The `Error` field uses `json:"-"` so it is never serialized.

**Concurrency pattern:** fan-out (one goroutine per source) → buffered channel → sequential drain in `main`. The buffer size equals `len(sources)` to prevent goroutines from blocking before the closer goroutine runs.

**Error handling:** Errors are not returned directly; they are sent as the `Error` field on `RateResponse` through the channel. The `main` drain loop checks `res.Error != nil` and prints the error before continuing.

## Key Conventions

- **Chinese comments**: All source code comments are in Traditional Chinese. Maintain this convention when adding new comments.
- **No external dependencies**: The module has zero third-party imports. Do not introduce any.
- **Single file**: All code lives in `main.go`. This is an intentional constraint for a focused learning example.
- **Hardcoded configuration**: The currency pair (`USD`/`TWD`) and source names are defined in `main()`. This is deliberate — do not over-engineer into config files.
- **Buffered channel sizing**: The channel buffer must equal `len(sources)`. Changing one without the other can cause deadlocks.

## Module Info

- Module name: `currency-tracker` (see `go.mod`)
- Go version: `1.26.2`
- External dependencies: none

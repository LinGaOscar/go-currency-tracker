# Dev Notes

## Setup

- Go 1.26.2+ (matches `go.mod`; verified working with `go1.26.4`).
- No dependencies to install — standard library only, no `go.sum`.

## Run

```bash
go run main.go
```

## Test

No test files exist yet (`go test ./...` currently reports "no test files"). If adding tests, `FetchRate` is the natural unit to cover — consider making the API base URL injectable so tests don't hit the live network.

## Structure

Single-file project:

- `main.go` — everything: `RateResponse` type, `FetchRate` (goroutine worker), `main` (fan-out/fan-in driver).
- `go.mod` — module `currency-tracker`, no external deps.

## Notes

- Network dependent: running the program requires internet access to `open.er-api.com`. No API key required.
- All three "sources" in `main.go` hit the same real endpoint — there's no actual multi-bank comparison logic, it's a concurrency demo.

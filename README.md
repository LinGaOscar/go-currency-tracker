# go-currency-tracker

A minimal Go program that demonstrates concurrent HTTP fetching using goroutines, channels, and `sync.WaitGroup`. It queries the public [open.er-api.com](https://www.exchangerate-api.com/docs/free) exchange-rate API from three simulated "sources" (`Bank_A`, `Bank_B`, `Bank_C`) in parallel, all fetching the same USD → TWD rate, and prints each result (or error) as it arrives.

## What it does

- Spawns one goroutine per source, each calling `https://open.er-api.com/v6/latest/USD`.
- Collects results through a buffered channel.
- Prints the fetched USD/TWD rate per source, or an error if the request/parse fails.

This is a learning/demo project for Go concurrency patterns (fan-out + channel + `WaitGroup`), not a production rate-comparison tool — all three sources hit the same API endpoint.

## Requirements

- Go 1.26.2 or later (see `go.mod`)
- No API key needed — `open.er-api.com`'s free tier requires no authentication.
- No external dependencies (standard library only).

## Build & Run

```bash
# Run directly
go run main.go

# Or build a binary
go build -o currency-tracker .
./currency-tracker
```

## Configuration

None. Base currency (`USD`) and target currency (`TWD`) are hardcoded in `main.go`.
